package bencoding

import (
  "strings"
  "strconv"
  "errors"
  "fmt"
  "unicode/utf8"
)

type Element struct {
  // 0 = string, 1 = int, 2 = list, 3 = dict
  ElementType int
  StringValue string
  ByteValue []byte
  IntValue int
  ListValue []Element
  DictValue map[string]Element
  UnparsedString string
}

func (e Element) String() string {
  switch e.ElementType {
  case 0:
    if e.ReadableString() {
      return fmt.Sprintf("%v", e.StringValue)
    } else {
      return fmt.Sprintf("<[]byte length:%v>", len(e.ByteValue))
    }
  case 1:
    return fmt.Sprintf("%v", e.IntValue)
  case 2:
    return fmt.Sprintf("%v", e.ListValue)
  case 3:
    return fmt.Sprintf("%v", e.DictValue)
  default:
    return ""
  }
}

func (e Element) ReadableString() bool {
  if e.ElementType == 0 {
    if utf8.Valid(e.ByteValue) {
      return true
    }
  }
  return false
}

func ParseString(s string) (elements []Element, err error) {
  for len(s) > 0 {
    elementType, err := nextElementType(s)
    if err != nil {
      return nil, err
    }
    var nextElement Element
    var nextElementLength int
    switch elementType {
    case 0:
      nextElement, nextElementLength, err = parseNextStringElement(s)
    case 1:
      nextElement, nextElementLength, err = parseNextIntElement(s)
    case 2:
      nextElement, nextElementLength, err = parseNextListElement(s)
    case 3:
      nextElement, nextElementLength, err = parseNextDictElement(s)
    }
    if err != nil {
      return nil, err
    }
    elements = append(elements, nextElement)
    s = s[nextElementLength:]
  }
  return
}

func nextElementType(s string) (elementType int, err error) {
  firstChar := string(s[0])
  switch firstChar {
  case "i":
    return 1, nil
  case "l":
    return 2, nil
  case "d":
    return 3, nil
  case "1","2","3","4","5","6","7","8","9":
    return 0, nil
  default:
    return -1, errors.New("The next element type isn't recogized")
  }
}

func parseNextStringElement(s string) (element Element, elementLength int, err error) {
  split := strings.Split(s, ":")
  stringLength, err := strconv.Atoi(split[0])
  if err != nil {
    return
  }
  stringLengthIndicatorStringLength := len(split[0]) + 1
  stringValue := string(s[stringLengthIndicatorStringLength:stringLength + stringLengthIndicatorStringLength])
  byteValue := []byte(stringValue)
  if len(byteValue) != stringLength {
    return Element{}, 0, errors.New("Invalid bendcoding: String value doesn't match string length.")
  }
  elementLength = len(split[0]) + 1 + stringLength
  element = Element{
    ElementType: 0,
    StringValue: stringValue,
    ByteValue: byteValue,
    UnparsedString: s[:elementLength],
  }
  return
}

func parseNextIntElement(s string) (element Element, elementLength int, err error) {
  endIndex := strings.Index(s, "e")
  elementString := s[:endIndex + 1]
  elementLength = len(elementString)
  intString := strings.TrimPrefix(elementString, "i")
  intString = strings.TrimSuffix(intString, "e")
  intValue, err := strconv.Atoi(intString)
  if err != nil {
    return
  }
  element = Element{
    ElementType: 1,
    IntValue: intValue,
    UnparsedString: s[:elementLength],
  }
  return
}

func parseNextListElement(s string) (element Element, elementLength int, err error) {
  originalS := s
  element = Element{ElementType: 2, ListValue: []Element{}}
  elementLength = 2
  s = strings.TrimPrefix(s, "l")
  for !strings.HasPrefix(s, "e") {
    elementType, err := nextElementType(s)
    if err != nil {
      return Element{}, -1, err
    }
    var nextElement Element
    var nextElementLength int
    switch elementType {
    case 0:
      nextElement, nextElementLength, err = parseNextStringElement(s)
    case 1:
      nextElement, nextElementLength, err = parseNextIntElement(s)
    case 2:
      nextElement, nextElementLength, err = parseNextListElement(s)
    case 3:
      nextElement, nextElementLength, err = parseNextDictElement(s)
    }
    if err != nil {
      return Element{}, -1, err
    }
    element.ListValue = append(element.ListValue, nextElement)
    s = s[nextElementLength:]
    elementLength += nextElementLength
  }
  element.UnparsedString = originalS[:elementLength - 1]
  return
}

func parseNextDictElement(s string) (element Element, elementLength int, err error) {
  originalS := s
  element = Element{ElementType: 3, DictValue: map[string]Element{}}
  elementLength = 2
  s = strings.TrimPrefix(s, "d")
  for !strings.HasPrefix(s, "e") {
    keyElementType, err := nextElementType(s)
    if err != nil {
      return Element{}, -1, err
    }
    if keyElementType != 0 {
      return Element{}, -1, errors.New("No key value found for dictionary") 
    }
    keyElement, keyElementLength, err := parseNextStringElement(s)
    s = s[keyElementLength:]

    valueElementType, err := nextElementType(s)
    if err != nil {
      return Element{}, -1, err
    }
    var valueElement Element
    var valueElementLength int
    switch valueElementType {
    case 0:
      valueElement, valueElementLength, err = parseNextStringElement(s)
    case 1:
      valueElement, valueElementLength, err = parseNextIntElement(s)
    case 2:
      valueElement, valueElementLength, err = parseNextListElement(s)
    case 3:
      valueElement, valueElementLength, err = parseNextDictElement(s)
    }
    if err != nil {
      return Element{}, -1, err
    }
    element.DictValue[keyElement.StringValue] = valueElement
    s = s[valueElementLength:]
    elementLength += keyElementLength + valueElementLength
  }
  element.UnparsedString = originalS[:elementLength]
  return
}
