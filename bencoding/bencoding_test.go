package bencoding

import (
  "testing"
  "io/ioutil"
  "bytes"
)

func TestParseString(t *testing.T) {
  // Can parse string elements
  in := "4:test"
  if x, e := ParseString(in);
    x[0].StringValue != "test" ||
    e != nil {
		t.Errorf("ParseString(%v) = (%v, %v), want (%v, %v)", in, x, e, "x[0].StringValue = 'test'", nil)
	}

  // Can parse int elements
  in = "i16e"
  if x, e := ParseString(in);
    x[0].IntValue != 16 ||
    e != nil {
    t.Errorf("ParseString(%v) = (%v, %v), want (%v, %v)", in, x, e, "x[0].IntValue = 16", nil)
  }

  // Can parse list elements
  in = "l4:testi16ee"
  if x, e := ParseString(in);
    x[0].ListValue[0].StringValue != "test" ||
    x[0].ListValue[1].IntValue != 16 ||
    e != nil {
    t.Errorf("ParseString(%v) = (%v, %v), want (%v, %v)", in, x, e, "A valid list", nil)
  }

  // Can parse dict elements
  in = "d15:meaning-of-lifei42ee"
  if x, e := ParseString(in);
    x[0].DictValue["meaning-of-life"].IntValue != 42 ||
    e != nil {
    t.Errorf("ParseString(%v) = (%v, %v), want (%v, %v)", in, x, e, "A valid dict", nil)
  }

  // Saves a copy of the unparsed string
  in = "d15:meaning-of-lifei42ee"
  if x, e := ParseString(in);
    x[0].UnparsedString != in ||
    e != nil {
    t.Errorf("ParseString(%v) = (%v, %v), want (%v, %v)", in, x, e, "UnparsedString", nil)
  }

  // Can handle all 4 types mixed
  in = "4:testi16el4:testi16eed15:meaning-of-lifei42ee"
  if x, e := ParseString(in);
    x[0].StringValue != "test" ||
    x[1].IntValue != 16 ||
    x[2].ListValue[0].StringValue != "test" ||
    x[2].ListValue[1].IntValue != 16 ||
    x[3].DictValue["meaning-of-life"].IntValue != 42 ||
    e != nil {
    t.Errorf("ParseString(%v) = (%v, %v), want (%v, %v)", in, x, e, "Mix of types", nil)
  }

  // Can handle nested types
  in = "d20:bencoding-data-typesl6:string3:int4:list4:dictee"
  if x, e := ParseString(in);
    x[0].DictValue["bencoding-data-types"].ListValue[0].StringValue != "string" ||
    x[0].DictValue["bencoding-data-types"].ListValue[1].StringValue != "int" ||
    x[0].DictValue["bencoding-data-types"].ListValue[2].StringValue != "list" ||
    x[0].DictValue["bencoding-data-types"].ListValue[3].StringValue != "dict" ||
    e != nil {
    t.Errorf("ParseString(%v) = (%v, %v), want (%v, %v)", in, x, e, "Nested types", nil)
  }

  // Real test file
  fileBytes,_ := ioutil.ReadFile("./test_resources/normandy.torrent")
  pieceBytes,_ := ioutil.ReadFile("./test_resources/piece_bytes.txt")
  pieceBytes = pieceBytes[:len(pieceBytes) - 1] // Extra byte in the txt file, not sure why
  in = string(fileBytes)
  if x, e := ParseString(in);
    x[0].DictValue["announce"].StringValue != "http://files.publicdomaintorrents.com/bt/announce.php" ||
    x[0].DictValue["creation date"].IntValue != 1119095315 ||
    x[0].DictValue["info"].DictValue["length"].IntValue != 734160896 ||
    x[0].DictValue["info"].DictValue["name"].StringValue != "Hell_in_Normandy.avi" ||
    x[0].DictValue["info"].DictValue["piece length"].IntValue != 262144 ||
    bytes.Compare(x[0].DictValue["info"].DictValue["pieces"].ByteValue, pieceBytes) != 0 ||
    e != nil {
    t.Errorf("ParseString(%v) = (%v, %v), want (%v, %v)", "<.torrent file>", x, e, "Real file comparison", nil)
  }
}
