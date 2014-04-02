package main

import (
  "encoding/binary"
  "errors"
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
  if len(bytes) < 6 {
    err = errors.New("Message is too short")
    return
  }
  lengthBytes := bytes[0:4]
  length := int(binary.BigEndian.Uint32(lengthBytes))
  idByte := bytes[4]
  id := int(idByte)
  payloadBytes := bytes[5:]

  message = &Message {
    Length:  length,
    Id:      id,
    Payload: payloadBytes,
  }
  return
}

func NewKeepAliveMessage() *Message {
  return &Message {
    Length: 0,
    Id: -1,
    Payload: []byte{},
  }
}

func NewChokeMessage() *Message {
  return &Message {
    Length: 1,
    Id: 0,
    Payload: []byte{},
  }
}

func NewUnchokeMessage() *Message {
  return &Message {
    Length: 1,
    Id: 1,
    Payload: []byte{},
  }
}

func NewInterestedMessage() *Message {
  return &Message {
    Length: 1,
    Id: 2,
    Payload: []byte{},
  }
}

func NewNotInterestedMessage() *Message {
  return &Message {
    Length: 1,
    Id: 3,
    Payload: []byte{},
  }
}

func NewHaveMessage(pieceIndex int) *Message {
  pieceIndexBytes := intToUint32Bytes(pieceIndex)
  return &Message {
    Length: 5,
    Id: 4,
    Payload: pieceIndexBytes,
  }
}

func NewBitfieldMessage(bitfield Bitfield) *Message {
  bitfieldBytes := bitfield.ToBytes()
  length := 1 + len(bitfieldBytes)
  return &Message {
    Length: length,
    Id: 5,
    Payload: bitfieldBytes,
  }
}

func NewRequestMessage(index int, begin int, length int) *Message {
  indexBytes := intToUint32Bytes(index)
  beginBytes := intToUint32Bytes(begin)
  lengthBytes := intToUint32Bytes(length)
  payload := []byte{}
  payload = append(payload, indexBytes...)
  payload = append(payload, beginBytes...)
  payload = append(payload, lengthBytes...)
  return &Message {
    Length: 13,
    Id:      6,
    Payload: payload,
  }
}

func NewPieceMessage(index int, begin int, block []byte) (m *Message) {
  indexBytes := intToUint32Bytes(index)
  beginBytes := intToUint32Bytes(begin)
  payload := []byte{}
  payload = append(payload, indexBytes...)
  payload = append(payload, beginBytes...)
  payload = append(payload, block...)
  m = &Message {
    Id: 7,
    Payload: payload,
  }
  m.Length = m.CalcLength()
  return
}

func NewCancelMessage(index int, begin int, length int) *Message {
  indexBytes := intToUint32Bytes(index)
  beginBytes := intToUint32Bytes(begin)
  lengthBytes := intToUint32Bytes(length)
  payload := []byte{}
  payload = append(payload, indexBytes...)
  payload = append(payload, beginBytes...)
  payload = append(payload, lengthBytes...)
  return &Message {
    Length: 13,
    Id:      8,
    Payload: payload,
  }
}

func (m *Message) DeliverableBytes() (delivery []byte) {
  lengthBytes := intToUint32Bytes(m.Length)
  delivery = append(delivery, lengthBytes...)
  if m.Id >= 0 {
    delivery = append(delivery, byte(m.Id))
  }
  if len(m.Payload) > 0 {
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
