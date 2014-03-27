package main

import (
  "encoding/binary"
  "errors"
  "bytes"
)

func intToUint32Bytes(n int) (intBytes []byte) {
  intBytes = make([]byte, 4)
  binary.BigEndian.PutUint32(intBytes, uint32(n))
  return
}

func uint32BytesToInt(intBytes []byte) (n int, err error) {
  if len(intBytes) != 4 {
    err = errors.New("Must pass 4 byte slice")
    return
  }
  var nInt32 int32
  buf := bytes.NewReader(intBytes)
  err = binary.Read(buf, binary.BigEndian, &nInt32)
  if err != nil {
    return
  }
  n = int(nInt32)
  return
}
