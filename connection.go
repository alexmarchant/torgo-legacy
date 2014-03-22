package main

import (
	"io"
  "net"
  "fmt"
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
		var n uint32
		n, err := readNBOUint32(c.TcpConn)
		if err != nil {
			break
		}
		if n > MessageByteLength {
      fmt.Println("Message size too large: \n", n)
			break
		}

		var buf []byte
		if n == 0 {
			// keep-alive - we want an empty message
			buf = make([]byte, 1)
		} else {
			buf = make([]byte, n)
		}

		_, err = io.ReadFull(c.TcpConn, buf)
		if err != nil {
			break
		}

    fmt.Printf("%v", buf)
	}
}

func readNBOUint32(conn net.Conn) (n uint32, err error) {
	var buf [4]byte
	_, err = conn.Read(buf[0:])
	if err != nil {
		return
	}
	n = bytesToUint32(buf[0:])
	return
}

func bytesToUint32(buf []byte) uint32 {
	return (uint32(buf[0]) << 24) |
		(uint32(buf[1]) << 16) |
		(uint32(buf[2]) << 8) | uint32(buf[3])
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

func (c *Connection) handshakeMessage() (message []byte) {
	pstrlen := []byte{19}
	pstr := []byte("BitTorrent protocol")
	reserved := []byte{0, 0, 0, 0, 0, 0, 0, 0}

	message = []byte{}
	message = append(message, pstrlen...)
	message = append(message, pstr...)
	message = append(message, reserved...)
	message = append(message, c.Metainfo.InfoDictionary.Hash...)
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
