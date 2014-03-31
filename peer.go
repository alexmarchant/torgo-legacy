package main

import (
  "strconv"
)

const (
  PeerStateNotConnected = iota
  PeerStateConnected
  PeerStateConnectionFailed
)

type Peer struct {
  Ip         string
  PeerId     string
  Port       int
  Connection *Connection
  State      int
}

func NewPeer() *Peer {
  return &Peer {
    State: PeerStateNotConnected,
  }
}

func (p *Peer) Address() string {
  return p.Ip + ":" + strconv.Itoa(p.Port)
}

func (p *Peer) Connect() (err error) {
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

func (p *Peer) connectionFailed() {
  p.State = PeerStateConnectionFailed
  p.disconnect()
}

func (p *Peer) connected() {
  p.State = PeerStateConnected
}

func (p *Peer) disconnect() {
  p.Connection.Close()
}

func (p *Peer) dial() (err error) {

}

func (p *Peer) sendHandshake() (err error) {

}

func (p *Peer) listenForHandshake() bool {

}
