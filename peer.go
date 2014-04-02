package main

import (
  "strconv"
  "net"
  "time"
  "fmt"
  "errors"
  "bytes"
)

const (
  PeerStateNotConnected = iota
  PeerStateConnecting
  PeerStateConnectionFailed
  PeerStateConnected
  PeerStateDownloading
)

var peerId = []byte("15620985492012023883")
var pieceSize = 16384

type Peer struct {
  Ip           string
  PeerId       string
  Port         int
  State        int
  AmChoking    bool
  AmInterested bool
  Choking      bool
  Interested   bool
  Connection   net.Conn
}

func NewPeer() *Peer {
  return &Peer {
    State:          PeerStateNotConnected,
    AmChoking:      true,
    AmInterested:   false,
    Choking:        true,
    Interested:     false,
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
  } else {
    p.connectionFailed()
  }
  return
}

func (p *Peer) DownloadPiece() (piece []byte, err error) {
  p.downloading()
  err = p.sendPieceRequest()
  if err != nil {
    p.connectionFailed()
    return
  }
  piece, err = p.listenForPiece()
  if err != nil {
    p.connectionFailed()
    return
  }
  fmt.Printf("%v size byte gotten\n", len(piece))
  p.connected()
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
    if timeout > 30 { break }
		var buf []byte
    buf = make([]byte, 512)
    _,err := p.Connection.Read(buf)
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

func (p *Peer) sendPieceRequest() (err error) {
  pieceRequest := NewPieceRequestMessage(0,0)
  pieceRequestBytes := pieceRequest.DeliverableBytes()
  num, err := p.Connection.Write(pieceRequestBytes)
  if err != nil {
    return
  }
  if num != len(pieceRequestBytes) {
    err = errors.New(fmt.Sprintf("Problem sending piece request message to: %v\n", p.Address()))
  }
  return
}

func (p *Peer) listenForPiece() (piece []byte, err error) {
  var timeout int
  for {
    timeout += 1
    if timeout > 30 { break }
		var buf []byte
    buf = make([]byte, 9 + pieceSize)
    _,err = p.Connection.Read(buf)
		if err != nil {
      return
		}
    if buf[0] != 0 {
      fmt.Printf("%v, %v\n", len(buf), string(buf))
      var message *Message
      message, err = ReadMessage(buf)
      if err != nil {
        return
      }
      if message.Id != 7 {
        fmt.Println("Listening for piece, got something else")
        break
      }
      if message.Length != message.CalcLength() {
        fmt.Println("Message length mismatch")
        break
      }
      piece = message.Payload
      return
    }
  }
  errors.New("Listen timeout")
  return 
}
