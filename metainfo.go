package main

import (
	"crypto/sha1"
	"./bencoding"
  "fmt"
  "os"
  "errors"
  "io/ioutil"
  "io"
)

type Metainfo struct {
  Announce       string
  AnnounceList   []string
  CreationDate   int
  InfoDictionary *MetainfoInfoDictionary
}

type MetainfoInfoDictionary struct {
  PieceLength    int
  Pieces         []byte
  Hash           []byte
  MultiFile      bool
  SingleFileInfo *SingleFileInfo
  Files  []*SingleFileInfo
}

type SingleFileInfo struct {
  Name   string
  Length int
  Md5Sum string
  Path   []string // For multifile items
}

func NewMetainfoFromFilename(filename string) (metainfo *Metainfo, err error) {
	if _, e := os.Stat(filename); os.IsNotExist(e) {
		errorMsg := fmt.Sprintf("No such file or directory: %s", filename)
		err = errors.New(errorMsg)
		return
	}
	torrentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
  elements, err := bencoding.ParseString(string(torrentBytes))
  if err != nil {
    return
  }
  metainfo = NewMetainfoFromBencoding(elements)
  return
}

func NewMetainfoFromBencoding(torrentElements []bencoding.Element) *Metainfo {
  announce := torrentElements[0].DictValue["announce"].StringValue
  var announceList []string = []string{}
  if _, present := torrentElements[0].DictValue["announce-list"]; present {
    for _, announceListElement := range torrentElements[0].DictValue["announce-list"].ListValue {
      for _, announceListListElement := range announceListElement.ListValue {
        announceList = append(announceList, announceListListElement.StringValue)
      }
    }
  }
  creationDate := torrentElements[0].DictValue["creation date"].IntValue
  torrentInfo := NewMetainfoInfoDictionary(torrentElements[0].DictValue["info"])

  return &Metainfo {
    Announce:       announce,
    AnnounceList:   announceList,
    CreationDate:   creationDate,
    InfoDictionary: torrentInfo,
  }
}

func NewMetainfoInfoDictionary(infoElement bencoding.Element) *MetainfoInfoDictionary {
  var multifile bool
  var singleFileInfo *SingleFileInfo
  var files []*SingleFileInfo

  pieceLength := infoElement.DictValue["piece length"].IntValue
  pieces := infoElement.DictValue["pieces"].ByteValue
	h := sha1.New()
	io.WriteString(h, infoElement.UnparsedString)
	infoHashBytes := h.Sum(nil)

  if _, present := infoElement.DictValue["files"]; present {
    multifile = true
    fileList := infoElement.DictValue["files"].ListValue
    for _, fileElement := range fileList {
      var md5Sum string
      length := fileElement.DictValue["length"].IntValue
      if _, present := fileElement.DictValue["md5sum"]; present {
        md5Sum = fileElement.DictValue["md5sum"].StringValue
      }
      path := []string{}
      if _, present := fileElement.DictValue["path"]; present {
        for _, pathPartElement := range fileElement.DictValue["path"].ListValue {

          path = append(path, pathPartElement.StringValue)
        }
      }
      mySingleFileInfo := &SingleFileInfo{
        Length: length,
        Md5Sum: md5Sum,
        Path:   path,
      }
      files = append(files, mySingleFileInfo)
    }
  } else {
    var md5Sum string
    multifile = false
    name := infoElement.DictValue["name"].StringValue
    length := infoElement.DictValue["length"].IntValue
    if _, present := infoElement.DictValue["md5sum"]; present {
      md5Sum = infoElement.DictValue["md5sum"].StringValue
    }
    singleFileInfo = &SingleFileInfo{
      Name:   name,
      Length: length,
      Md5Sum: md5Sum,
    }
  }

  return &MetainfoInfoDictionary {
    PieceLength:    pieceLength,
    Pieces:         pieces,
    Hash:           infoHashBytes,
    MultiFile:      multifile,
    SingleFileInfo: singleFileInfo,
    Files:          files,
  }
}

func (m *MetainfoInfoDictionary) Length() int {
  if m.MultiFile {
    length := 0
    for _, file := range m.Files {
      length += file.Length
    }
    return length
  } else {
    return m.SingleFileInfo.Length
  }
}

func (m *Metainfo) AllTrackers() (trackers []*Tracker) {
  trackers = append(trackers, NewTracker(m.Announce, m))
  for _, tracker := range m.AnnounceList {
    trackers = append(trackers, NewTracker(tracker, m))
  }
  return trackers
}
