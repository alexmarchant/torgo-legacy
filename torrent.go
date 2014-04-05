package main

import (
  "sync"
  "fmt"
  "errors"
  "time"
  "os"
  "strconv"
)

var downloadPath = "/Users/alex/Downloads"
var maxPeers = 60
var downloadMessageChan = make(chan *Message)
var lastCheckedDownloadLength int
var lastCheckedDownloadTime time.Time

type Torrent struct {
  Metainfo *Metainfo
  Files []*File
  Trackers []*Tracker
  Peers []*Peer
  Progress *Progress
}

func NewTorrent(metainfoFilePath string) (torrent *Torrent, err error) {
  var metainfo *Metainfo
  var progress *Progress
  metainfo, err = NewMetainfoFromFilename(metainfoFilePath)
  if err != nil {
    return
  }
  progress = NewProgress(
    metainfo.InfoDictionary.PiecesCount(),
    metainfo.InfoDictionary.PieceLength)
  torrent = &Torrent {
    Metainfo: metainfo,
    Progress: progress,
  }
  err = torrent.getFiles()
  return
}

func (t *Torrent) Start() (err error) {
	err = t.getPeers()
	if err != nil {
    return
	}
  t.tryConnectingToMorePeers(maxPeers)
  go t.downloadLoop()
  t.monitorPeerConnections()
  return
}

func (t *Torrent) downloadLoop() {
  for {
    select {
    case <- downloadMessageChan:
      message := <- downloadMessageChan
      err := t.receivedBlock(message)
      if err != nil {
        fmt.Println(err)
      }
    default:
    }
  }
}

func (t *Torrent) monitorPeerConnections() {
  for {
    lowPeerCount := maxPeers - t.activePeerCount()
    if lowPeerCount > 0 {
      t.tryConnectingToMorePeers(lowPeerCount)
    }
    t.assignPiecesToPeers()
    t.printStatus()

    amt := time.Duration(100)
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
  for _,file := range t.Files {
    file.CreatePartFile()
  }
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

func (t *Torrent) assignPiecesToPeers() {
  for _,peer := range t.Peers {
    if peer.State == PeerStateConnected {
      block, err := t.Progress.randomNeededBlock()
      if err == nil {
        m := NewRequestMessage(block)
        go peer.sendMessage(m)
      }
    }
  }
}

func (t *Torrent) printStatus() {
  fmt.Printf("Downloading from %v of %v peers - DL: %vKB/s, UL: %vKB/s - Percent done: %v%\r", 0, t.connectedPeerCount(), t.downloadSpeed(), t.uploadSpeed(), t.percentComplete())
}

func (t *Torrent) percentComplete() string {
  floatPercent := t.Progress.percentComplete()
  return strconv.FormatFloat(floatPercent, 'f', 3, 64)
}

func (t *Torrent) downloadSpeed() string {
  newLength := t.Progress.downloadedLength()
  newTime := time.Now()
  if lastCheckedDownloadLength != 0 {
    deltaDownloadedLength := (newLength - lastCheckedDownloadLength) / 1000
    deltaTime := newTime.Sub(lastCheckedDownloadTime).Seconds()
    speed := float64(deltaDownloadedLength) / deltaTime
    return strconv.FormatFloat(speed, 'f', 1, 64)
  }
  lastCheckedDownloadLength = newLength
  lastCheckedDownloadTime = newTime
  return "0.0"
}

func (t *Torrent) uploadSpeed() string {
  return "0.0"
}

func (t *Torrent) receivedBlock(m *Message) (err error) {
  var index, begin int
  var block *Block
  index, err = uint32BytesToInt(m.Payload[0:4])
  begin, err = uint32BytesToInt(m.Payload[4:8])
  if err != nil {
    block.state = BlockStateNeed
    return
  }
  block, err = t.Progress.findBlockFor(index, begin)
  if err != nil {
    block.state = BlockStateNeed
    return
  }
  blockBytes := m.Payload[8:]
  if len(blockBytes) != block.length {
    block.state = BlockStateNeed
    err = errors.New("Block length mismatch")
    return
  }
  err = t.writeBlockToFile(block, blockBytes)
  if err != nil {
    block.state = BlockStateNeed
    return
  }
  block.state = BlockStateHave
  return
}

func (t *Torrent) writeBlockToFile(block *Block, blockBytes []byte) (err error) {
  var osFile *os.File
  if len(t.Files) > 1 {
    err = errors.New("TODO MultiFile")
  }
  offset := block.index * t.Metainfo.InfoDictionary.PieceLength + block.begin
  file := t.Files[0]
  osFile, err = file.OpenPartFileWrite()
  if err != nil {
    return
  }
  _, err = osFile.WriteAt(blockBytes, int64(offset))
  return
}
