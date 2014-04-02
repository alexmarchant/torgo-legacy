package main

import (
  "fmt"
  "net/url"
	"io/ioutil"
	"net/http"
	"net"
	"time"
)

type Tracker struct {
  Url             string
  Metainfo        *Metainfo
  TrackerResponse *TrackerResponse
}

func NewTracker(url string, metainfo *Metainfo) *Tracker {
  return &Tracker {
    Url:      url,
    Metainfo: metainfo,
    TrackerResponse: &TrackerResponse{
      State: TrackerResponseStateNotSent,
    },
  }
}

func (t *Tracker) SendRequest() (err error) {
  responseBody, err := t.httpGetPeerRequest()
  if err != nil {
    return
  }
  t.TrackerResponse, err = NewTrackerResponse(responseBody)
  return
}

func (t *Tracker) urlWithParams() string {
	params := t.generateParams()
	return t.Url + params
}

func (t *Tracker) httpGetPeerRequest() (responseBody []byte, err error) {
  //fmt.Printf("GET %v\n", t.urlWithParams())
  transport := http.Transport{
    Dial: dialTimeout,
  }
  client := http.Client{
    Transport: &transport,
  }
  resp, err := client.Get(t.urlWithParams())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	responseBody, err = ioutil.ReadAll(resp.Body)
  if err != nil {
    return
  }
	return
}

func (t *Tracker) generateParams() (urlParams string) {
	paramDict := map[string]string{}
	paramDict["info_hash"] = t.infoHashParam()
	paramDict["peer_id"] = t.peerIdParam()
	paramDict["port"] = t.portParam()
	paramDict["uploaded"] = t.uploadedParam()
	paramDict["downloaded"] = t.downloadedParam()
	paramDict["left"] = t.leftParam()
	paramDict["compact"] = t.compactParam()
	paramDict["no_peer_id"] = t.noPeerIdParam()
	paramDict["event"] = t.eventParam()

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

func (t *Tracker) infoHashParam() string {
  hash := string(t.Metainfo.InfoDictionary.Hash)
	return url.QueryEscape(hash)
}

func (t *Tracker) peerIdParam() string {
	return "15620985492012023883"
}

func (t *Tracker) leftParam() string {
  // TODO make this dynamic
	return fmt.Sprintf("%v", t.Metainfo.InfoDictionary.Length())
}

func (t *Tracker) portParam() string {
	return "6881"
}

func (t *Tracker) uploadedParam() string {
	return "0"
}

func (t *Tracker) downloadedParam() string {
	return "0"
}

func (t *Tracker) compactParam() string {
	return "1"
}

func (t *Tracker) noPeerIdParam() string {
	return "0"
}

func (t *Tracker) eventParam() string {
	return "started"
}

func dialTimeout(network, addr string) (net.Conn, error) {
  var timeout = time.Duration(2 * time.Second)
  return net.DialTimeout(network, addr, timeout)
}
