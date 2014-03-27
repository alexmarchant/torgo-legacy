package main

import (
  "sync"
  "log"
  "fmt"
)

var downloadPath = "/Users/alex/Downloads"

type Torrent struct {
  Metainfo *Metainfo
  Files []*File
  Trackers []*Tracker
}

func NewTorrent(metainfoFilePath string) (torrent *Torrent, err error) {
  metainfo, err := NewMetainfoFromFilename(metainfoFilePath)
  if err != nil {
    return
  }
  files, err := getFiles(metainfo)
  if err != nil {
    return
  }
  torrent = &Torrent {
    Metainfo: metainfo,
    Files:    files,
  }
  return
}

func (t *Torrent) StartDownloading() (err error) {
	peers, err := getPeers(t.Metainfo)
	if err != nil {
		log.Fatal(err)
	}
	if len(peers) == 0 {
		log.Fatal("No peers found")
	}
  cp := NewConnectionPool(peers, t.Metainfo)
  cp.Start()
  return
}

func getFiles(metainfo *Metainfo) (files []*File, err error) {
  if metainfo.InfoDictionary.MultiFile {
    for _,fileInfo := range metainfo.InfoDictionary.Files {
      files = append(files, NewFile(fileInfo))
    }
  } else {
    fileInfo := metainfo.InfoDictionary.SingleFileInfo
    files = append(files, NewFile(fileInfo))
  }
}

func getPeers(metainfo *Metainfo) (peers []*Peer, err error) {
  var wg sync.WaitGroup
  trackers := metainfo.AllTrackers()
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
  fmt.Printf("%v peers added\n", len(peers))
	return
}
