package main

import (
	"github.com/alexmarchant/bencoding"
	"io/ioutil"
	"testing"
)

func TestNewTrackerResponse(t *testing.T) {
	peerResponse := getSamplePeerResponse()
	tr, e := NewTrackerResponse(peerResponse)

	if e != nil {
		t.Errorf("Error parsing binary peer list: %v", e)
	}

	if len(tr.Peers) != 8 {
		t.Errorf("Peer list length doesn't match declared length")
	}

	if tr.Interval != 1800 {
		t.Errorf("Error generating peer list interval")
	}

	if tr.State != TrackerResponseStateSuccess {
		t.Errorf("Error generating tracker response state")
	}
}

func TestParseBinaryPeerList(t *testing.T) {
	peerList := getSampleBinaryPeerList()
	peers, e := parseBinaryPeerList(peerList)

	if e != nil {
		t.Errorf("Error parsing binary peer list: %v", e)
	}

	if len(peers) != 8 {
		t.Errorf("Peer list length doesn't match declared length")
	}
}

func getSamplePeerResponse() (fileBytes []byte) {
	sampleTrackerResponseFilepath := "./test_resources/sample_binary_peer_response.txt"
	fileBytes, _ = ioutil.ReadFile(sampleTrackerResponseFilepath)
	return
}

func getSampleBinaryPeerList() (peerList bencoding.Element) {
	fileBytes := getSamplePeerResponse()
	benResponse, _ := bencoding.ParseString(string(fileBytes))
	peerList = benResponse[0].DictValue["peers"]
	return
}
