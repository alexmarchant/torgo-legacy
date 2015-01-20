package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

const (
	PeerStateNotConnected = iota
	PeerStateConnecting
	PeerStateConnectionFailed
	PeerStateConnected
	PeerStateDownloading
)

var peerId = []byte("15620985492012023883")

type Peer struct {
	Ip           string
	PeerId       string
	Port         int
	State        int
	AmChoking    bool
	AmInterested bool
	Choked       bool
	Interested   bool
	Connection   net.Conn
	Bitfield     Bitfield
}

func NewPeer() *Peer {
	return &Peer{
		State:        PeerStateNotConnected,
		AmChoking:    true,
		AmInterested: false,
		Choked:       true,
		Interested:   false,
	}
}

func (p *Peer) Address() string {
	return p.Ip + ":" + strconv.Itoa(p.Port)
}

func (p *Peer) Connect() (err error) {
	p.connecting()
	err = p.dial()
	if err != nil {
		p.connectionFailed()
		return
	}
	err = p.sendHandshake()
	if err != nil {
		p.connectionFailed()
		return
	}
	receivedHandshake := p.listenForHandshake()
	if receivedHandshake {
		p.connected()
		p.listen()
	} else {
		p.connectionFailed()
	}
	return
}

func (p *Peer) connecting() {
	p.State = PeerStateConnecting
}

func (p *Peer) connectionFailed() {
	p.State = PeerStateConnectionFailed
	p.disconnect()
}

func (p *Peer) connected() {
	p.State = PeerStateConnected
}

func (p *Peer) downloading() {
	p.State = PeerStateDownloading
}

func (p *Peer) disconnect() {
	if p.Connection != nil {
		p.Connection.Close()
	}
}

func (p *Peer) dial() (err error) {
	timeout, _ := time.ParseDuration("3s")
	p.Connection, err = net.DialTimeout("tcp", p.Address(), timeout)
	return
}

func (p *Peer) sendHandshake() (err error) {
	handshakeMessage := NewHandshakeMessage(torrent.Metainfo)
	handshakeBytes := handshakeMessage.DeliverableBytes()
	num, err := p.Connection.Write(handshakeBytes)
	if err != nil {
		return
	}
	if num != len(handshakeBytes) {
		err = errors.New(fmt.Sprintf("Problem sending handshake to: %v\n", p.Address()))
	}
	return
}

func (p *Peer) listenForHandshake() bool {
	var timeout int
	for {
		timeout += 1
		if timeout > 30 {
			break
		}
		var buf []byte
		buf = make([]byte, 512)
		_, err := p.Connection.Read(buf)
		if err != nil {
			return false
		}
		if buf[0] != 0 {
			hm, err := ParseHandshakeMessageFromBytes(buf)
			if err != nil {
				return false
			}
			if bytes.Equal(hm.InfoHash, torrent.Metainfo.InfoDictionary.Hash) {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func (p *Peer) listen() (piece []byte, err error) {
	for {
		var message *Message
		message, err = p.readMessage()
		if err != nil {
			break
		}
		p.actOnMessage(message)
	}
	return
}

func (p *Peer) readMessage() (m *Message, err error) {
	var length int
	var lengthBytes []byte

	lengthBytes, length, err = p.readMessageLength()
	if err != nil {
		return
	}
	if length > 130*1024 {
		err = fmt.Errorf("Message size too large: ", length)
		return
	}

	var buf []byte
	if length == 0 {
		// keep-alive - we want an empty message
		buf = make([]byte, 1)
	} else {
		buf = make([]byte, length)
	}

	_, err = io.ReadFull(p.Connection, buf)
	if err != nil {
		return
	}

	// put length back on
	fullMessageBytes := append(lengthBytes, buf...)
	m, err = ReadMessage(fullMessageBytes)
	return
}

func (p *Peer) readMessageLength() (bytes []byte, length int, err error) {
	bytes = make([]byte, 4)
	_, err = p.Connection.Read(bytes)
	if err != nil {
		return
	}
	length, err = uint32BytesToInt(bytes)
	return
}

func (p *Peer) sendMessage(m *Message) (err error) {
	var num int
	messageBytes := m.DeliverableBytes()
	num, err = p.Connection.Write(messageBytes)
	if err != nil {
		return
	}
	if num != len(messageBytes) {
		err = errors.New("Message sending unsuccesful")
	}
	return
}

func (p *Peer) actOnMessage(m *Message) {
	switch m.Id {
	case -1: // keep-alive
	case 0:
		p.Choked = true
	case 1:
		p.Choked = false
	case 2:
		p.Interested = true
	case 3:
		p.Interested = false
	case 4: // have
		pieceIndex, err := uint32BytesToInt(m.Payload)
		if err != nil {
			break
		}
		p.Bitfield[pieceIndex] = true
	case 5: // bitfield
		piecesCount := torrent.Metainfo.InfoDictionary.PiecesCount()
		bitfield, err := BytesToBitfield(m.Payload, piecesCount)
		if err != nil {
			break
		}
		p.Bitfield = bitfield
	case 6: // request
	case 7: // piece
		downloadMessageChan <- m
	case 8: // cancel
	case 9: // port
	default:
		fmt.Println("No message id, that shouldn't happen")
	}
}
