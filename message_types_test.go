package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewKeepAliveMessage(t *testing.T) {
	m := NewKeepAliveMessage()
	deliverableBytes := m.DeliverableBytes()
	expectedPayloadBytes := []byte{}
	expectedDeliverableBytes := []byte{0, 0, 0, 0}

	expectedMessage := &Message{
		Length:  0,
		Id:      -1,
		Payload: expectedPayloadBytes,
	}

	if !reflect.DeepEqual(m, expectedMessage) {
		t.Errorf("Expected %+v, got %+v", expectedMessage, m)
	}

	if !bytes.Equal(deliverableBytes, expectedDeliverableBytes) {
		t.Errorf("Expected %v, got %v", expectedDeliverableBytes, deliverableBytes)
	}
}

func TestNewChokeMessage(t *testing.T) {
	m := NewChokeMessage()
	deliverableBytes := m.DeliverableBytes()
	expectedPayloadBytes := []byte{}
	expectedDeliverableBytes := []byte{0, 0, 0, 1, 0}

	expectedMessage := &Message{
		Length:  1,
		Id:      0,
		Payload: expectedPayloadBytes,
	}

	if !reflect.DeepEqual(m, expectedMessage) {
		t.Errorf("Expected %+v, got %+v", expectedMessage, m)
	}

	if !bytes.Equal(deliverableBytes, expectedDeliverableBytes) {
		t.Errorf("Expected %v, got %v", expectedDeliverableBytes, deliverableBytes)
	}
}

func TestNewUnchokeMessage(t *testing.T) {
	m := NewUnchokeMessage()
	deliverableBytes := m.DeliverableBytes()
	expectedPayloadBytes := []byte{}
	expectedDeliverableBytes := []byte{0, 0, 0, 1, 1}

	expectedMessage := &Message{
		Length:  1,
		Id:      1,
		Payload: expectedPayloadBytes,
	}

	if !reflect.DeepEqual(m, expectedMessage) {
		t.Errorf("Expected %+v, got %+v", expectedMessage, m)
	}

	if !bytes.Equal(deliverableBytes, expectedDeliverableBytes) {
		t.Errorf("Expected %v, got %v", expectedDeliverableBytes, deliverableBytes)
	}
}

func TestNewInterestedMessage(t *testing.T) {
	m := NewInterestedMessage()
	deliverableBytes := m.DeliverableBytes()
	expectedPayloadBytes := []byte{}
	expectedDeliverableBytes := []byte{0, 0, 0, 1, 2}

	expectedMessage := &Message{
		Length:  1,
		Id:      2,
		Payload: expectedPayloadBytes,
	}

	if !reflect.DeepEqual(m, expectedMessage) {
		t.Errorf("Expected %+v, got %+v", expectedMessage, m)
	}

	if !bytes.Equal(deliverableBytes, expectedDeliverableBytes) {
		t.Errorf("Expected %v, got %v", expectedDeliverableBytes, deliverableBytes)
	}
}

func TestNewNotInterestedMessage(t *testing.T) {
	m := NewNotInterestedMessage()
	deliverableBytes := m.DeliverableBytes()
	expectedPayloadBytes := []byte{}
	expectedDeliverableBytes := []byte{0, 0, 0, 1, 3}

	expectedMessage := &Message{
		Length:  1,
		Id:      3,
		Payload: expectedPayloadBytes,
	}

	if !reflect.DeepEqual(m, expectedMessage) {
		t.Errorf("Expected %+v, got %+v", expectedMessage, m)
	}

	if !bytes.Equal(deliverableBytes, expectedDeliverableBytes) {
		t.Errorf("Expected %v, got %v", expectedDeliverableBytes, deliverableBytes)
	}
}

func TestNewHaveMessage(t *testing.T) {
	m := NewHaveMessage(257)
	deliverableBytes := m.DeliverableBytes()
	expectedPayloadBytes := []byte{0, 0, 1, 1}
	expectedDeliverableBytes := []byte{0, 0, 0, 5, 4, 0, 0, 1, 1}

	expectedMessage := &Message{
		Length:  5,
		Id:      4,
		Payload: expectedPayloadBytes,
	}

	if !reflect.DeepEqual(m, expectedMessage) {
		t.Errorf("Expected %+v, got %+v", expectedMessage, m)
	}

	if !bytes.Equal(deliverableBytes, expectedDeliverableBytes) {
		t.Errorf("Expected %v, got %v", expectedDeliverableBytes, deliverableBytes)
	}
}

func TestNewBitfieldMessage(t *testing.T) {
	bitfield := Bitfield{true, true, true, false, false, true, true, true, false, true}
	m := NewBitfieldMessage(bitfield)
	deliverableBytes := m.DeliverableBytes()
	expectedPayloadBytes := []byte{231, 64}
	expectedDeliverableBytes := []byte{0, 0, 0, 3, 5, 231, 64}

	expectedMessage := &Message{
		Length:  3,
		Id:      5,
		Payload: expectedPayloadBytes,
	}

	if !reflect.DeepEqual(m, expectedMessage) {
		t.Errorf("Expected %+v, got %+v", expectedMessage, m)
	}

	if !bytes.Equal(deliverableBytes, expectedDeliverableBytes) {
		t.Errorf("Expected %v, got %v", expectedDeliverableBytes, deliverableBytes)
	}
}

func TestNewRequestMessage(t *testing.T) {
	block := &Block{
		index:  1,
		begin:  1,
		length: 16384,
		state:  BlockStateDownloading,
	}
	m := NewRequestMessage(block)
	deliverableBytes := m.DeliverableBytes()
	expectedPayloadBytes := []byte{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 64, 0}
	expectedDeliverableBytes := []byte{0, 0, 0, 13, 6, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 64, 0}

	expectedMessage := &Message{
		Length:  13,
		Id:      6,
		Payload: expectedPayloadBytes,
	}

	if !reflect.DeepEqual(m, expectedMessage) {
		t.Errorf("Expected %+v, got %+v", expectedMessage, m)
	}

	if !bytes.Equal(deliverableBytes, expectedDeliverableBytes) {
		t.Errorf("Expected %v, got %v", expectedDeliverableBytes, deliverableBytes)
	}
}

func TestNewPieceMessage(t *testing.T) {
	block := []byte{241, 13, 1, 9, 245, 133}
	m := NewPieceMessage(1, 1, block)
	deliverableBytes := m.DeliverableBytes()
	expectedPayloadBytes := []byte{0, 0, 0, 1, 0, 0, 0, 1, 241, 13, 1, 9, 245, 133}
	expectedDeliverableBytes := []byte{0, 0, 0, 15, 7, 0, 0, 0, 1, 0, 0, 0, 1, 241, 13, 1, 9, 245, 133}

	expectedMessage := &Message{
		Length:  15,
		Id:      7,
		Payload: expectedPayloadBytes,
	}

	if !reflect.DeepEqual(m, expectedMessage) {
		t.Errorf("Expected %+v, got %+v", expectedMessage, m)
	}

	if !bytes.Equal(deliverableBytes, expectedDeliverableBytes) {
		t.Errorf("Expected %v, got %v", expectedDeliverableBytes, deliverableBytes)
	}
}

func TestNewCancelMessage(t *testing.T) {
	m := NewCancelMessage(1, 1, MessageByteLength)
	deliverableBytes := m.DeliverableBytes()
	expectedPayloadBytes := []byte{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 64, 0}
	expectedDeliverableBytes := []byte{0, 0, 0, 13, 8, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 64, 0}

	expectedMessage := &Message{
		Length:  13,
		Id:      8,
		Payload: expectedPayloadBytes,
	}

	if !reflect.DeepEqual(m, expectedMessage) {
		t.Errorf("Expected %+v, got %+v", expectedMessage, m)
	}

	if !bytes.Equal(deliverableBytes, expectedDeliverableBytes) {
		t.Errorf("Expected %v, got %v", expectedDeliverableBytes, deliverableBytes)
	}
}
