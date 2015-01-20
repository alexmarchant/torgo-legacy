package main

import (
	"bytes"
	"testing"
)

func TestDeliverableBytes(t *testing.T) {
	block := &Block{
		index:  0,
		begin:  0,
		length: 16384,
		state:  BlockStateDownloading,
	}
	m := NewRequestMessage(block)
	deliverable_bytes := m.DeliverableBytes()
	expected_bytes := []byte{0, 0, 0, 13, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0}

	if !bytes.Equal(deliverable_bytes, expected_bytes) {
		t.Errorf("Expected %v, got %v", expected_bytes, deliverable_bytes)
	}
}

func TestReadMessageBytes(t *testing.T) {
	message_bytes := []byte{0, 0, 0, 13, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0}
	m, e := ReadMessage(message_bytes)
	expected_payload := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0}

	if e != nil {
		t.Errorf("Got error: %v", e)
	}

	if m.Length != 13 {
		t.Errorf("Message length parsed incorrectly")
	}

	if m.Id != 6 {
		t.Errorf("Message id parsed incorrectly")
	}

	if !bytes.Equal(m.Payload, expected_payload) {
		t.Errorf("Expected %v, got %v", expected_payload, m.Payload)
	}
}
