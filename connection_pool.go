package main

import ()

type ConnectionPool struct {
	Connections []*Connection
}

const (
	MaxPeers = 10
)

func NewConnectionPool(peers []*Peer, metainfo *Metainfo) *ConnectionPool {
	var connections []*Connection
	for _, peer := range peers {
		c := NewConnection(peer, metainfo)
		connections = append(connections, c)
	}
	return &ConnectionPool{
		Connections: connections,
	}
}

func (cp *ConnectionPool) Start() {
	for i, connection := range cp.Connections {
		if i > MaxPeers {
			break
		}
		connection.Open()
	}
	for i, connection := range cp.OpenConnections() {
		if i > MaxPeers {
			break
		}
		go connection.Listen()
		connection.SendHandshakeMessage()
	}
	//for i, connection := range cp.OpenConnections() {
	//if i > MaxPeers { break }
	//interestedMessage := InterestedMessage()
	//requestMessage := RequestMessage(0,0)
	//connection.SendMessage(interestedMessage)
	//connection.SendMessage(requestMessage)
	//}
}

func (cp *ConnectionPool) OpenConnections() (connections []*Connection) {
	for _, connection := range cp.Connections {
		if connection.State == ConnectionStateConnected {
			connections = append(connections, connection)
		}
	}
	return
}
