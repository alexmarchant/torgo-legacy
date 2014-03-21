package main

import (
	"./bencoding"
  "fmt"
)

type ConnectionPool struct {
  Connections []*Connection
}

func NewConnectionPool(peers []Peer, torrentInfo []bencoding.Element) *ConnectionPool {
  var connections []*Connection
  for _, peer := range peers {
    c := NewConnection(peer, torrentInfo)
    connections = append(connections, c)
  }
  return &ConnectionPool {
    Connections: connections,
  }
}

func (cp *ConnectionPool) Start() {
  for _, connection := range cp.Connections {
    connection.Open()
  }
  for _, connection := range cp.OpenConnections() {
    connection.SendHandshakeMessage()
  }
}

func (cp *ConnectionPool) OpenConnections() (connections []*Connection) {
  for _, connection := range cp.Connections {
    if connection.State == ConnectionStateConnected {
      connections = append(connections, connection)
    }
  }
  return
}

func (cp *ConnectionPool) PrintReport() {
  fmt.Printf("Peers: %v\n", len(cp.Connections))
  for i := -1; i <= 1; i++ {
    var count int
    for _, connection := range cp.Connections {
      if connection.State == i {
        count++
      }
    }
    fmt.Printf("%v: %v\n", ConnectionStateString(i), count)
  }
}
