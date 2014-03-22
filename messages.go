package main

import (
  "encoding/binary"
)

const (
  MessageByteLength = 16384
)

type Message struct {
  Id      int
  Payload []byte
}

func (m *Message) DeliverableBytes() (delivery []byte) {
  delivery = append(delivery, m.Length()...)
  if m.Id != 0 {
    delivery = append(delivery, byte(m.Id))
  }
  if len(m.Payload) == 0 {
    delivery = append(delivery, m.Payload...)
  }
  return
}

func (m *Message) Length() (bytes []byte) {
  length := 0
  if (m.Id >= 0) {
    length += 1
  }
  length += len(m.Payload)
  bytes = make([]byte, 4)
  binary.BigEndian.PutUint32(bytes, uint32(length))
  return
}

func KeepAliveMessage() *Message {
  return &Message {
    Id: -1,
  }
}

func InterestedMessage() *Message {
  return &Message {
    Id: 2,
  }
}

func RequestMessage(index int, begin int) *Message {
  indexBytes := make([]byte, 4)
  binary.BigEndian.PutUint32(indexBytes, uint32(index))
  beginBytes := make([]byte, 4)
  binary.BigEndian.PutUint32(beginBytes, uint32(begin))
  lengthBytes := make([]byte, 4)
  binary.BigEndian.PutUint32(lengthBytes, uint32(MessageByteLength))
  payload := []byte{}
  payload = append(payload, indexBytes...)
  payload = append(payload, beginBytes...)
  payload = append(payload, lengthBytes...)
  return &Message {
    Id:      6,
    Payload: payload,
  }
}

