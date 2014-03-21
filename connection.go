package main

import (
	"crypto/sha1"
	"./bencoding"
	"io"
  "net"
  "fmt"
  //"bufio"
  "time"
  "errors"
)

const (
  ConnectionStateFailed     = -1
  ConnectionStateConnecting = 0
  ConnectionStateConnected  = 1
)

var peerId = []byte("15620985492012023883")

type Connection struct {
	Peer           Peer
	TorrentInfo    []bencoding.Element
	AmChoking      bool
	AmInterested   bool
	PeerChoking    bool
	PeerInterested bool
  TcpConn        net.Conn
  State          int
}

func NewConnection(peer Peer, torrentInfo []bencoding.Element) *Connection {
  return &Connection{
    Peer:           peer,
    TorrentInfo:    torrentInfo,
    AmChoking:      true,
    AmInterested:   false,
    PeerChoking:    true,
    PeerInterested: false,
    State:          0,
  }
}

func (c *Connection) Open() (err error) {
  timeout, _ := time.ParseDuration("3s")
  c.TcpConn, err = net.DialTimeout("tcp", c.Peer.Address(), timeout)
  if err != nil {
    c.State = ConnectionStateFailed
    return
  }
  c.State = ConnectionStateConnected
	return
}

func (c *Connection) Close() {
  c.TcpConn.Close()
}

func (c *Connection) SendHandshakeMessage() (err error) {
  handshakeBytes := c.handshakeMessage()
  num, err := c.TcpConn.Write(handshakeBytes)
  if err != nil {
    c.State = ConnectionStateFailed
    return
  }
  if num != len(handshakeBytes) {
    c.State = ConnectionStateFailed
    err = errors.New(fmt.Sprintf("Problem sending handshake to: %v\n", c.Peer.Ip))
  }
  return
}

func (c *Connection) handshakeMessage() (message []byte) {

	pstrlen := []byte{19}
	pstr := []byte("BitTorrent protocol")
	reserved := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	h := sha1.New()
	io.WriteString(h, c.TorrentInfo[0].DictValue["info"].UnparsedString)
	infoHash := h.Sum(nil)

	message = []byte{}
	message = append(message, pstrlen...)
	message = append(message, pstr...)
	message = append(message, reserved...)
	message = append(message, infoHash...)
	message = append(message, peerId...)

	return
}

func (c *Connection) StateString() string {
  return ConnectionStateString(c.State)
}

func ConnectionStateString(state int) string {
  switch state {
  case ConnectionStateFailed: return "Failed to connect"
  case ConnectionStateConnecting: return "Connecting"
  case ConnectionStateConnected: return "Connected"
  default: return "Unknown"
  }
}
