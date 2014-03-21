package main

import (
  "strconv"
)

type Peer struct {
  Ip     string
  PeerId string
  Port   int
}

func (p Peer) Address() string {
  return p.Ip + ":" + strconv.Itoa(p.Port)
}
