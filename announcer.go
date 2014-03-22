package main

import (
  "fmt"
  "errors"
  "net/url"
	"io/ioutil"
	"net/http"
	"./bencoding"
)

type Announcer struct {
  Url      string
  Metainfo *Metainfo
  Peers    []*Peer
}

func NewAnnouncer(url string, metainfo *Metainfo) *Announcer {
  return &Announcer {
    Url:      url,
    Metainfo: metainfo,
  }
}

func (a *Announcer) GetPeers() (err error) {
  fmt.Printf("Attempting to connect to tracker: %v\n", a.Url)
  body, err := a.httpGetPeerRequest()
  if err != nil {
    fmt.Printf("Connection to tracker failed.\nMessage: %v\n", err)
  }
  if body == "" {
    err = errors.New("Could not connect to a tracker")
    return
  }
  fmt.Printf("Successfully connected to tracker: %v\n", a.Url)
  a.Peers, err = a.parsePeerResponse(body)

  for _, peer := range a.Peers {
    fmt.Printf("Peer added %v\n", peer)
  }
  if err != nil {
    return
  }
  fmt.Printf("Peer list downloaded\n")
  return
}

func (a *Announcer) urlWithParams() string {
	params := a.generateParams()
	return a.Url + params
}

func (a *Announcer) httpGetPeerRequest() (body string, err error) {
  fmt.Printf("GET %v\n", a.urlWithParams())
	resp, err := http.Get(a.urlWithParams())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	body = string(bodyBytes)
  fmt.Printf("Peer list string %v\n", body)
  bencodeError := bencoding.ValidBencodingString(body)
  if bencodeError != nil {
    err = fmt.Errorf("Peer list response invalid:\nPeer list: %v\nError message: %v\n", body, bencodeError)
    body = ""
    return
  }
	return
}

func (a *Announcer) generateParams() (urlParams string) {
	paramDict := map[string]string{}
	paramDict["info_hash"] = a.infoHashParam()
	paramDict["peer_id"] = a.peerIdParam()
	paramDict["port"] = a.portParam()
	paramDict["uploaded"] = a.uploadedParam()
	paramDict["downloaded"] = a.downloadedParam()
	paramDict["left"] = a.leftParam()
	paramDict["compact"] = a.compactParam()
	paramDict["no_peer_id"] = a.noPeerIdParam()
	paramDict["event"] = a.eventParam()

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

func (a *Announcer) infoHashParam() string {
  hash := string(a.Metainfo.InfoDictionary.Hash)
	return url.QueryEscape(hash)
}

func (a *Announcer) peerIdParam() string {
	return "15620985492012023883"
}

func (a *Announcer) leftParam() string {
  // TODO make this dynamic
	return fmt.Sprintf("%v", a.Metainfo.InfoDictionary.Length())
}

func (a *Announcer) portParam() string {
	return "6881"
}

func (a *Announcer) uploadedParam() string {
	return "0"
}

func (a *Announcer) downloadedParam() string {
	return "0"
}

func (a *Announcer) compactParam() string {
	return "1"
}

func (a *Announcer) noPeerIdParam() string {
	return "0"
}

func (a *Announcer) eventParam() string {
	return "started"
}

func (a *Announcer) parsePeerResponse(body string) (peers []*Peer, err error) {
  // TODO this is dictory model, take care of binary model
	benResponse, err := bencoding.ParseString(body)
  fmt.Printf("bencode parsed body %v\n", benResponse)
	if err != nil {
		return
	}
	peerList := benResponse[0].DictValue["peers"].ListValue
	for _, peerElement := range peerList {
		newPeer := &Peer{
			Ip:     peerElement.DictValue["ip"].StringValue,
			PeerId: peerElement.DictValue["peer id"].StringValue,
			Port:   peerElement.DictValue["port"].IntValue,
		}
		peers = append(peers, newPeer)
	}
	return
}
