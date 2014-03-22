package main

import (
  "fmt"
)

type ConnectionPool struct {
  Connections []*Connection
}

func NewConnectionPool(peers []*Peer, metainfo *Metainfo) *ConnectionPool {
  var connections []*Connection
  for _, peer := range peers {
    c := NewConnection(peer, metainfo)
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
    go connection.Listen()
    connection.SendHandshakeMessage()
  }
  for _, connection := range cp.OpenConnections() {
    connection.SendMessage(InterestedMessage())
    connection.SendMessage(RequestMessage(0,0))
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
