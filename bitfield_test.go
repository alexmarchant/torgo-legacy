package main

import (
	"bytes"
	"testing"
)

func TestToBytes(t *testing.T) {
	bitfield := Bitfield{true, true, true, false, false, true, true, true, false, true}
	bitfieldBytes := bitfield.ToBytes()
	expectedBytes := []byte{231, 64}

	if !bytes.Equal(bitfieldBytes, expectedBytes) {
		t.Errorf("Expected: %v, got: %v", expectedBytes, bitfieldBytes)
	}
}

func TestBytesToBitfield(t *testing.T) {
	bytes := []byte{231, 64}
	bitfield, e := BytesToBitfield(bytes, 10)
	expectedBitfield := Bitfield{true, true, true, false, false, true, true, true, false, true}

	if e != nil {
		t.Errorf("%v", e)
	}

	if !BitfieldsEqual(bitfield, expectedBitfield) {
		t.Errorf("Expected: %v, got: %v", expectedBitfield, bitfield)
	}
}
