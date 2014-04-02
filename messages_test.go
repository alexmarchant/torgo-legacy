package main

import (
  "testing"
  "bytes"
)

func TestNewPieceRequestMessage(t *testing.T) {
  m := NewPieceRequestMessage(0,0)

  if m.Id != 6 {
    t.Errorf("Wrong message id")
  }

  if m.Length != 13 {
    t.Errorf("Wrong message length")
  }

  if !bytes.Equal(m.Payload, []byte{0,0,0,0,0,0,0,0,0,0,64,0}) {
    t.Errorf("Wrong message payload")
  }
}

func TestDeliverableBytes(t *testing.T) {
  m := NewPieceRequestMessage(0,0)
  deliverable_bytes := m.DeliverableBytes()
  expected_bytes := []byte{0,0,0,13,6,0,0,0,0,0,0,0,0,0,0,64,0}

  if !bytes.Equal(deliverable_bytes, expected_bytes) {
    t.Errorf("Expected %v, got %v", expected_bytes, deliverable_bytes)
  }
}

func ReadMessageBytes(t *testing.T) {
  message_bytes := []byte{0,0,0,13,6,0,0,0,0,0,0,0,0,0,0,64,0}
  m, e := ReadMessage(message_bytes)
  expected_payload := []byte{0,0,0,0,0,0,0,0,0,0,64,0}

  if e != nil {
    t.Errorf("Got error: %v", e)
  }

  if m.Length != 13 {
    t.Errorf("Message length parsed incorrectly")
  }

  if m.Id != 6 {
    t.Errorf("Message id parsed incorrectly")
  }

  if bytes.Equal(m.Payload, expected_payload) {
    t.Errorf("Expected &v, got &v", expected_payload, m.Payload)
  }
}
