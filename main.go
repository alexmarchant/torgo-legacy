package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"./bencoding"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
	torrentInfo, err := getTorrentData()
	if err != nil {
		log.Fatal(err)
	}
	peers, err := getPeers(torrentInfo)
	if err != nil {
		log.Fatal(err)
	}
  cp := NewConnectionPool(peers, torrentInfo)
  cp.Start()
  cp.PrintReport()
}

func getTorrentData() (elements []bencoding.Element, err error) {
	if len(os.Args) <= 2 {
		err = errors.New("Please pass in the location of the .torrent file")
		return
	}
	filename := os.Args[2]
	if _, e := os.Stat(filename); os.IsNotExist(e) {
		errorMsg := fmt.Sprintf("No such file or directory: %s", filename)
		err = errors.New(errorMsg)
		return
	}
	torrentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	elements, err = bencoding.ParseString(string(torrentBytes))
	return
}


func getPeers(torrentInfo []bencoding.Element) (peers []Peer, err error) {
	requestUrl := generatePeerRequestUrl(torrentInfo)
	body, err := httpGetPeerRequest(requestUrl)
	if err != nil {
		return
	}
	peers, err = parsePeerResponse(body)
	return
}

func generatePeerRequestUrl(torrentInfo []bencoding.Element) string {
	trackerUrl := torrentInfo[0].DictValue["announce"].StringValue
	params := generateParams(torrentInfo)
	return trackerUrl + params
}

func httpGetPeerRequest(requestUrl string) (body string, err error) {
	resp, err := http.Get(requestUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	body = string(bodyBytes)
	return
}

func parsePeerResponse(body string) (peers []Peer, err error) {
	benResponse, err := bencoding.ParseString(body)
	if err != nil {
		return
	}
	peerList := benResponse[0].DictValue["peers"].ListValue
	for _, peerElement := range peerList {
		newPeer := Peer{
			Ip:     peerElement.DictValue["ip"].StringValue,
			PeerId: peerElement.DictValue["peer id"].StringValue,
			Port:   peerElement.DictValue["port"].IntValue,
		}
		peers = append(peers, newPeer)
	}
	return
}

func generateParams(torrentInfo []bencoding.Element) (urlParams string) {
	paramDict := map[string]string{}
	paramDict["info_hash"] = infoHashParam(torrentInfo)
	paramDict["peer_id"] = peerIdParam(torrentInfo)
	paramDict["port"] = portParam(torrentInfo)
	paramDict["uploaded"] = uploadedParam(torrentInfo)
	paramDict["downloaded"] = downloadedParam(torrentInfo)
	paramDict["left"] = leftParam(torrentInfo)
	paramDict["compact"] = compactParam(torrentInfo)
	paramDict["no_peer_id"] = noPeerIdParam(torrentInfo)
	paramDict["event"] = eventParam(torrentInfo)

	for key, value := range paramDict {
		if urlParams == "" {
			urlParams += "?"
		} else {
			urlParams += "&"
		}
		urlParams += key + "=" + value
	}
	return
}

func infoHashParam(torrentInfo []bencoding.Element) string {
	h := sha1.New()
	io.WriteString(h, torrentInfo[0].DictValue["info"].UnparsedString)
	infoHashBytes := h.Sum(nil)
	return url.QueryEscape(string(infoHashBytes))
}

func peerIdParam(torrentInfo []bencoding.Element) string {
	return "15620985492012023883"
}

func leftParam(torrentInfo []bencoding.Element) string {
	// TODO multi file lengths
	return fmt.Sprintf("%v", torrentInfo[0].DictValue["info"].DictValue["length"].IntValue)
}

func portParam(torrentInfo []bencoding.Element) string {
	return "6881"
}

func uploadedParam(torrentInfo []bencoding.Element) string {
	return "0"
}

func downloadedParam(torrentInfo []bencoding.Element) string {
	return "0"
}

func compactParam(torrentInfo []bencoding.Element) string {
	return "0"
}

func noPeerIdParam(torrentInfo []bencoding.Element) string {
	return "0"
}

func eventParam(torrentInfo []bencoding.Element) string {
	return "started"
}
