package main

import (
  "math"
  "strconv"
  "errors"
)

type Bitfield []bool

func (b Bitfield) ToBytes() []byte {
  byteCount := numberOfBytesNeededForNumberOfBits(len(b))
  bitfieldBytes := make([]byte, byteCount)
  for i := 0; i < byteCount; i++ {
    var byteString string
    for y := 0; y < 8; y++ {
      index := (i * 8) + y
      if index < len(b) && b[index] {
        byteString += "1"
      } else {
        byteString += "0"
      }
    }
    bitfieldInt,_ := strconv.ParseInt(byteString, 2, 0)
    bitfieldByte := byte(bitfieldInt)
    bitfieldBytes[i] = bitfieldByte
  }
  return bitfieldBytes
}

func BytesToBitfield(bytes []byte, length int) (bitfield Bitfield, err error) {
  var bitstring string
  byteCount := numberOfBytesNeededForNumberOfBits(length)
  if byteCount != len(bytes) {
    err = errors.New("Length mismatch")
    return
  }
  bitfield = make(Bitfield, length)
  for _,bitfieldByte := range bytes {
    byteString := strconv.FormatUint(uint64(bitfieldByte), 2)
    byteStringMissingBits := 8 - len(byteString)
    if byteStringMissingBits > 0 {
      for i := 0; i < byteStringMissingBits; i++ {
        byteString = "0" + byteString
      }
    }
    bitstring += byteString
  }
  for i,bit := range bitstring {
    if i >= length { break }
    if string(bit) == "1" {
      bitfield[i] = true
    } else {
      bitfield[i] = false
    }
  }
  return
}

func BitfieldsEqual(a Bitfield, b Bitfield) bool {
  if len(a) != len(b) {
    return false
  }
  for i,_ := range a {
    if a[i] != b[i] {
      return false
    }
  }
  return true
}

func numberOfBytesNeededForNumberOfBits(bitCount int) int {
  count64 := float64(bitCount)
  return int(math.Ceil(count64 / float64(8)))
}
