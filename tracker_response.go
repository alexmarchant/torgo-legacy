package main

import (
	"github.com/alexmarchant/bencoding"
  "fmt"
  "math"
  "strings"
  "bytes"
  "encoding/binary"
  "strconv"
)

const (
  TrackerResponseStateNotSent = iota
  TrackerResponseStateError
  TrackerResponseStateSuccess
)

type TrackerResponse struct {
  Interval int
  Peers    []*Peer
  State    int
}

func NewTrackerResponse(responseBytes []byte) (tr *TrackerResponse, err error) {
  tr = &TrackerResponse{
    State: TrackerResponseStateNotSent,
  }
	benResponse, err := bencoding.ParseString(string(responseBytes))
	if err != nil {
    tr.State = TrackerResponseStateError
		return
	}
  var peers []*Peer
  peerList := benResponse[0].DictValue["peers"]
  if peerList.ElementType == bencoding.BencodingListType {
    peers = append(peers, parseBencodingPeerList(peerList)...)
  } else {
    newPeers, err := parseBinaryPeerList(peerList)
    if err != nil {
      tr.State = TrackerResponseStateError
      return nil, err
    }
    peers = append(peers, newPeers...)
  }
  interval := benResponse[0].DictValue["interval"].IntValue

  tr.Interval = interval
  tr.Peers = peers
  tr.State = TrackerResponseStateSuccess
  return
}

func parseBencodingPeerList(peerList bencoding.Element) (peers []*Peer) {
	for _, peerElement := range peerList.ListValue {
    newPeer := NewPeer()
    newPeer.Ip = peerElement.DictValue["ip"].StringValue
    newPeer.PeerId = peerElement.DictValue["peer id"].StringValue
    newPeer.Port = peerElement.DictValue["port"].IntValue
		peers = append(peers, newPeer)
	}
  return
}

func parseBinaryPeerList(peerList bencoding.Element) (peers []*Peer, err error) {
  peerListBytes := peerList.ByteValue
  peerCountF := float64(len(peerListBytes)) / 6
  peerCountI := int(math.Floor(peerCountF))
  for i := 0; i < peerCountI; i++ {
    offset := i * 6
    peerBytes := peerListBytes[offset:offset + 6]
    peerIpBytes := peerBytes[0:3]
    peerPortBytes := peerBytes[4:6]
    peerIpString, err := ipStringFromBytes(peerIpBytes)
    if err != nil {
      return nil, err
    }
    peerPort, err := portStringFromBytes(peerPortBytes)
    if err != nil {
      return nil, err
    }
    newPeer := NewPeer()
    newPeer.Ip = peerIpString
    newPeer.Port = peerPort
    peers = append(peers, newPeer)
  }
  return
}

func ipStringFromBytes(ipBytes []byte) (ipString string, err error) {
  var ipStringParts []string
  for i := 0; i < 4; i++ {
    ipByte := ipBytes[i:i+1]
    var byteInt uint8
    buf := bytes.NewReader(ipByte)
    err := binary.Read(buf, binary.BigEndian, &byteInt)
    if err != nil {
      fmt.Println("ip, binary.Read failed:", err)
    }
    byteString := strconv.FormatUint(uint64(byteInt), 10)
    ipStringParts = append(ipStringParts, byteString)
  }
  ipString = strings.Join(ipStringParts, ".")
  return
}

func portStringFromBytes(portBytes []byte) (portInt int, err error) {
  var byteInt uint16
  buf := bytes.NewReader(portBytes)
  err = binary.Read(buf, binary.BigEndian, &byteInt)
  if err != nil {
    fmt.Println("port, binary.Read failed:", err)
  }
  portInt = int(byteInt)
  return
}
