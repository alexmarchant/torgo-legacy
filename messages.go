package main

import (
  "encoding/binary"
)

const (
  MessageByteLength = 16384
)

type Message struct {
  Length  int
  Id      int
  Payload []byte
}

func ReadMessage(bytes []byte) (message *Message, err error) {
  lengthBytes := bytes[0:4]
  length := int(binary.BigEndian.Uint32(lengthBytes))
  idBytes := bytes[4:5]
  id := int(binary.BigEndian.Uint32(idBytes))
  payloadBytes := bytes[5:]

  message = &Message {
    Length:  length,
    Id:      id,
    Payload: payloadBytes,
  }
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

func RequestMessage(index int, begin int) (m *Message) {
  indexBytes := intToUint32Bytes(index)
  beginBytes := intToUint32Bytes(begin)
  lengthBytes := intToUint32Bytes(MessageByteLength)
  payload := []byte{}
  payload = append(payload, indexBytes...)
  payload = append(payload, beginBytes...)
  payload = append(payload, lengthBytes...)
  m = &Message {
    Id:      6,
    Payload: payload,
  }
  m.Length = m.CalcLength()
  return
}


func (m *Message) DeliverableBytes() (delivery []byte) {
  lengthBytes := intToUint32Bytes(m.Length)
  delivery = append(delivery, lengthBytes...)
  if m.Id != 0 {
    delivery = append(delivery, byte(m.Id))
  }
  if len(m.Payload) == 0 {
    delivery = append(delivery, m.Payload...)
  }
  return
}

func (m *Message) CalcLength() (length int) {
  if (m.Id >= 0) {
    length += 1
  }
  length += len(m.Payload)
  return
}
