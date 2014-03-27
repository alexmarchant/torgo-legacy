package main

import (
  "errors"
)

type HandshakeMessage struct {
  PStrLen  int
  PStr     string
  InfoHash []byte
  PeerId   []byte
}

func ParseHandshakeMessageFromBytes(b []byte) (hm *HandshakeMessage, err error) {
  if len(b) == 0 {
    err = errors.New("Can't parse empty buffer")
  }
  pstrlen := int(b[0])

  if pstrlen != 19 {
    err = errors.New("Invalid handshake message")
    return
  }

  pstrStart := 1
  pstrEnd := pstrStart + pstrlen
  if len(b) < pstrEnd {
    err = errors.New("Unexpected EOF")
    return
  }
  pstr := string(b[pstrStart:pstrEnd])

  infoHashStart := pstrEnd + 8
  infoHashEnd := infoHashStart + 20
  if len(b) < infoHashEnd {
    err = errors.New("unexpected EOF")
    return
  }
  infoHash := b[infoHashStart:infoHashEnd]

  peerIdStart := infoHashEnd
  peerIdEnd := infoHashEnd + 20
  if len(b) < peerIdEnd {
    err = errors.New("unexpected EOF")
    return
  }
  peerId := b[peerIdStart:peerIdEnd]
  hm = &HandshakeMessage{
    PStrLen:  pstrlen,
    PStr:     pstr,
    InfoHash: infoHash,
    PeerId:   peerId,
  }
  return
}

func NewHandshakeMessage(m *Metainfo) *HandshakeMessage {
	pstrlen := 19
	pstr := "BitTorrent protocol"

  return &HandshakeMessage {
    PStrLen:  pstrlen,
    PStr:     pstr,
    InfoHash: m.InfoDictionary.Hash,
    PeerId:   peerId,
  }
}

func (hm *HandshakeMessage) DeliverableBytes() (delivery []byte) {
  pstrlen := []byte{byte(hm.PStrLen)}
  pstr := []byte(hm.PStr)
	reserved := []byte{0, 0, 0, 0, 0, 0, 0, 0}
  infoHash := hm.InfoHash
  peerId := []byte(hm.PeerId)

	delivery = append(delivery, pstrlen...)
	delivery = append(delivery, pstr...)
	delivery = append(delivery, reserved...)
	delivery = append(delivery, infoHash...)
	delivery = append(delivery, peerId...)
  return
}
