package main

import(
  "testing"
  "bytes"
)

func TestIntToUint32Bytes(t *testing.T) {
  n := 29284291
  nBytes := intToUint32Bytes(n)
  expectedBytes := []byte{1,190,215,195}

  if !bytes.Equal(nBytes, expectedBytes) {
    t.Errorf("Expected: %v, Got: %v", expectedBytes, nBytes)
  }
}

func TestUint32BytesToInt(t *testing.T) {
  intBytes := []byte{1,190,215,195}
  n, e := uint32BytesToInt(intBytes)
  expectedN := 29284291

  if e != nil {
    t.Errorf("Error: %v", e)
  }

  if n != expectedN {
    t.Errorf("Expected: %v, Got: %v", expectedN, n)
  }
}
