package main

import (
	"io"
  "net"
  "fmt"
  "time"
  "errors"
  "bytes"
)

const (
  ConnectionStateFailed     = -1
  ConnectionStateConnecting = 0
  ConnectionStateConnected  = 1
  ConnectionStateHandshakeSent
  ConnectionStateHandshakeReceived
)

var peerId = []byte("15620985492012023883")

type Connection struct {
	Peer           *Peer
	Metainfo       *Metainfo
	AmChoking      bool
	AmInterested   bool
	PeerChoking    bool
	PeerInterested bool
  TcpConn        net.Conn
  State          int
}

func NewConnection(peer *Peer, metainfo *Metainfo) *Connection {
  return &Connection{
    Peer:           peer,
    Metainfo:       metainfo,
    AmChoking:      true,
    AmInterested:   false,
    PeerChoking:    true,
    PeerInterested: false,
    State:          ConnectionStateConnecting,
  }
}

func (c *Connection) Open() (err error) {
  timeout, _ := time.ParseDuration("3s")
  fmt.Printf("Attempting to connect to peer: %v\n", c.Peer.Address())
  c.TcpConn, err = net.DialTimeout("tcp", c.Peer.Address(), timeout)
  if err != nil {
    fmt.Printf("Failed to connect to peer: %v\n", c.Peer.Address())
    c.State = ConnectionStateFailed
    return
  }
  fmt.Printf("Successfully connected to peer: %v\n", c.Peer.Address())
  c.State = ConnectionStateConnected
	return
}

func (c *Connection) Close() {
  defer c.TcpConn.Close()
}

func (c *Connection) Listen() {
  fmt.Printf("Listening for responses from: %v\n", c.Peer.Address())
  for {
		var buf []byte
    buf = make([]byte, 512)
    _,err := c.TcpConn.Read(buf)
		if err != nil {
      fmt.Println(err)
			break
		}
    if buf[0] != 0 {
      hm, err := ParseHandshakeMessageFromBytes(buf)
      if err != nil {
        fmt.Println(err)
        fmt.Println(buf[0:100])
        break
      }
      if !bytes.Equal(hm.InfoHash, c.Metainfo.InfoDictionary.Hash) {
        c.Close()
        fmt.Println("Disconnected, info hash didn't match")
        break
      }
      fmt.Println("Handshake received, matches info")
    } else {
      fmt.Println("Ping")
    }
	}
}

func readNBOUint32(conn net.Conn) (n int, err error) {
	var buf [4]byte
	_, err = conn.Read(buf[0:])
	if err != nil {
		return
	}
	n, err = uint32BytesToInt(buf[0:])
	return
}

func (c *Connection) ReadPiece() {
  fmt.Printf("Downloading block from %v\n", c.Peer.Address())
  buf := make([]byte, MessageByteLength)
  n, err := io.ReadFull(c.TcpConn, buf)
  fmt.Printf("Received %d bytes in response: %#v\n", n, buf[:n])
  if err != nil {
      fmt.Println("read error:", err)
  }
}

func (c *Connection) SendHandshakeMessage() (err error) {
  fmt.Printf("Sending handshake to peer: %v\n", c.Peer.Address())
  handshakeMessage := NewHandshakeMessage(c.Metainfo)
  handshakeBytes := handshakeMessage.DeliverableBytes()
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

func (c *Connection) SendMessage(m *Message) (err error) {
  fmt.Printf("Sending message to peer: %v\n", c.Peer.Address())
  messageBytes := m.DeliverableBytes()
  num, err := c.TcpConn.Write(messageBytes)
  if err != nil {
    c.State = ConnectionStateFailed
    return
  }
  if num != len(messageBytes) {
    c.State = ConnectionStateFailed
    err = errors.New(fmt.Sprintf("Problem sending message to: %v\n", c.Peer.Ip))
  }
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
