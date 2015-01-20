package main

import (
	"fmt"
	"log"
	"os"
)

const (
	MaxPeerCount = 60
)

var torrent *Torrent

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
	var err error
	if len(os.Args) <= 2 {
		log.Fatal("Please pass in the location of the .torrent file")
	}
	filename := os.Args[2]
	torrent, err = NewTorrent(filename)
	if err != nil {
		log.Fatal(err)
	}
	err = torrent.Start()
	if err != nil {
		log.Fatal(err)
	}
	var input string
	fmt.Scanln(&input)
}
