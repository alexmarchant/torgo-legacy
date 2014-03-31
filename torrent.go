package main

import (
  "sync"
  "log"
  "fmt"
  "errors"
)

var downloadPath = "/Users/alex/Downloads"

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
		log.Fatal(err)
	}
  return
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
  fmt.Printf("%v peers added\n", len(peers))
  t.Peers = peers
	return
}
