package main

import (
	"errors"
	"fmt"
	"log"
	"os"
  "sync"
)

const (
  MaxPeerCount = 60
)

func main() {
	parseCliInput()
}

func parseCliInput() {
	if len(os.Args) <= 1 {
		usage := `
USAGE: torgo <task>

The torgo tasks are:
  download <torrent file>    Downloads the specified torrent
`
		log.Fatal(usage)
	}
	switch os.Args[1] {
	case "download":
		download()
	default:
		errorMsg := fmt.Sprintf("We don't recognize the command: %v", os.Args[1])
		log.Fatal(errorMsg)
	}
}

func download() {
	metainfo, err := getTorrentData()
	if err != nil {
		log.Fatal(err)
	}
	peers, err := getPeers(metainfo)
	if err != nil {
		log.Fatal(err)
	}
  cp := NewConnectionPool(peers, metainfo)
  cp.Start()

  var input string
  fmt.Scanln(&input)
}

func getTorrentData() (metainfo *Metainfo, err error) {
	if len(os.Args) <= 2 {
		err = errors.New("Please pass in the location of the .torrent file")
		return
	}
	filename := os.Args[2]
  metainfo, err = NewMetainfoFromFilename(filename)
	return
}

func getPeers(metainfo *Metainfo) (peers []*Peer, err error) {
  var wg sync.WaitGroup
  announcers := metainfo.AllAnnouncers()
  for _, announcer := range announcers {
    wg.Add(1)
    go func(announcer *Announcer) {
      defer wg.Done()
      announcer.GetPeers()
    }(announcer)
  }
  wg.Wait()
  for _, announcer := range announcers {
    peers = append(peers, announcer.Peers...)
  }
  fmt.Printf("%v peers added\n", len(peers))
	return
}
