package main

import (
  "sync"
  "fmt"
  "errors"
  "time"
)

var downloadPath = "/Users/alex/Downloads"
var maxPeers = 10

type Torrent struct {
  Metainfo *Metainfo
  Files []*File
  Trackers []*Tracker
  Peers []*Peer
}

func NewTorrent(metainfoFilePath string) (torrent *Torrent, err error) {
  metainfo, err := NewMetainfoFromFilename(metainfoFilePath)
  if err != nil {
    return
  }
  torrent = &Torrent {
    Metainfo: metainfo,
  }
  err = torrent.getFiles()
  return
}

func (t *Torrent) StartDownloading() (err error) {
	err = t.getPeers()
	if err != nil {
    return
	}
  t.tryConnectingToMorePeers(maxPeers)
  t.downloadLoop()
  return
}

func (t *Torrent) downloadLoop() {
  for {
    lowPeerCount := maxPeers - t.activePeerCount()
    if lowPeerCount > 0 {
      t.tryConnectingToMorePeers(lowPeerCount)
    }
    t.tryDownloadingPieces()

    t.printStatus()

    amt := time.Duration(3000)
    time.Sleep(time.Millisecond * amt)
  }
}

func (t *Torrent) getFiles() (err error) {
  var files []*File
  if t.Metainfo.InfoDictionary.MultiFile {
    for _,fileInfo := range t.Metainfo.InfoDictionary.Files {
      files = append(files, NewFile(fileInfo))
    }
  } else {
    fileInfo := t.Metainfo.InfoDictionary.SingleFileInfo
    files = append(files, NewFile(fileInfo))
  }
  t.Files = files
  return
}

func (t *Torrent) getPeers() (err error) {
  var peers []*Peer
  var wg sync.WaitGroup
  trackers := t.Metainfo.AllTrackers()
  for _, tracker := range trackers {
    wg.Add(1)
    go func(tracker *Tracker) {
      defer wg.Done()
      tracker.SendRequest()
    }(tracker)
  }
  wg.Wait()
  for _, tracker := range trackers {
    if tracker.TrackerResponse.State == TrackerResponseStateSuccess {
      peers = append(peers, tracker.TrackerResponse.Peers...)
    }
  }
  if len(peers) == 0 {
    err = errors.New("No peers found")
    return
  }
  t.Peers = peers
	return
}

func (t *Torrent) totalPeerCount() (count int) {
  count = len(t.Peers)
  return
}

func (t *Torrent) activePeerCount() (count int) {
  count += t.connectingPeerCount()
  count += t.connectedPeerCount()
  return
}

func (t *Torrent) connectingPeerCount() (count int) {
  for _, peer := range t.Peers {
    if peer.State == PeerStateConnecting {
      count++
    }
  }
  return
}

func (t *Torrent) connectedPeerCount() (count int) {
  for _, peer := range t.Peers {
    if peer.State == PeerStateConnected {
      count++
    }
  }
  return
}

func (t *Torrent) tryConnectingToMorePeers(count int) {
  currentCount := 0
  for _, peer := range t.Peers {
    if currentCount >= count { return }
    if peer.State == PeerStateNotConnected {
      go peer.Connect()
      currentCount++
    }
  }
}

func (t *Torrent) tryDownloadingPieces() {
  for _,peer := range t.Peers {
    if peer.State == PeerStateConnected {
      go peer.DownloadPiece()
    }
  }
}

func (t *Torrent) printStatus() {
  fmt.Printf("Downloading from %v of %v peers - DL: %vKB/s, UL: %vKB/s\n", 0, t.connectedPeerCount(), t.downloadSpeed(), t.uploadSpeed())
}

func (t *Torrent) downloadSpeed() float64 {
  return 0
}

func (t *Torrent) uploadSpeed() float64 {
  return 0
}
