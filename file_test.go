package main

import(
  "testing"
  "os"
)

func TestNewFileForSingleFileTorrent(t *testing.T) {
  sfi := &SingleFileInfo {
    Name: "Hell_in_Normandy.avi",
    Length: 734160896,
  }
  file := NewFile(sfi)

  if file.Filename != "Hell_in_Normandy.avi" {
    t.Errorf("Expected: %v, Got: %v", sfi.Name, file.Filename)
  }

  if file.Path != "/Users/alex/Downloads" {
    t.Errorf("Expected: %v, Got: %v", "/Users/alex/Downloads", file.Path)
  }

  if file.Length != 734160896 {
    t.Errorf("Expected: %v, Got: %v", 734160896, file.Length)
  }

  if file.FullPath() != "/Users/alex/Downloads/Hell_in_Normandy.avi" {
    t.Errorf("Expected: %v, Got: %v", "/Users/alex/Downloads/Hell_in_Normandy.avi", file.FullPath())
  }

  if file.FullPartPath() != "/Users/alex/Downloads/Hell_in_Normandy.avi.part" {
    t.Errorf("Expected: %v, Got: %v", "/Users/alex/Downloads/Hell_in_Normandy.avi.part", file.FullPath())
  }
}

func TestNewFileForMultiFileTorrent(t *testing.T) {
  sfi := &SingleFileInfo {
    Name: "Google Chrome 16.0 [Win 32 & 64 Bit] Full Vesrion - {RedDragon}",
    Length: 24086584,
    Path: []string{"chrome installer 16.0 .exe"},
  }
  file := NewFile(sfi)

  if file.Filename != "chrome installer 16.0 .exe" {
    t.Errorf("Expected: %v, Got: %v", sfi.Name, file.Filename)
  }

  if file.Path != "/Users/alex/Downloads/Google Chrome 16.0 [Win 32 & 64 Bit] Full Vesrion - {RedDragon}" {
    t.Errorf("Expected: %v, Got: %v", "/Users/alex/Downloads/Google Chrome 16.0 [Win 32 & 64 Bit] Full Vesrion - {RedDragon}", file.Path)
  }

  if file.FullPath() != "/Users/alex/Downloads/Google Chrome 16.0 [Win 32 & 64 Bit] Full Vesrion - {RedDragon}/chrome installer 16.0 .exe" {
    t.Errorf("Expected: %v, Got: %v", "/Users/alex/Downloads/Google Chrome 16.0 [Win 32 & 64 Bit] Full Vesrion - {RedDragon}/chrome installer 16.0 .exe", file.FullPath())
  }
}

func TestNewFileForMultiFileNestedTorrent(t *testing.T) {
  sfi := &SingleFileInfo {
    Name: "Google Chrome 16.0 [Win 32 & 64 Bit] Full Vesrion - {RedDragon}",
    Length: 24086584,
    Path: []string{"installer","chrome installer 16.0 .exe"},
  }
  file := NewFile(sfi)

  if file.Filename != "chrome installer 16.0 .exe" {
    t.Errorf("Expected: %v, Got: %v", sfi.Name, file.Filename)
  }

  if file.Path != "/Users/alex/Downloads/Google Chrome 16.0 [Win 32 & 64 Bit] Full Vesrion - {RedDragon}/installer" {
    t.Errorf("Expected: %v, Got: %v", "/Users/alex/Downloads/Google Chrome 16.0 [Win 32 & 64 Bit] Full Vesrion - {RedDragon}/installer", file.Path)
  }

  if file.FullPath() != "/Users/alex/Downloads/Google Chrome 16.0 [Win 32 & 64 Bit] Full Vesrion - {RedDragon}/installer/chrome installer 16.0 .exe" {
    t.Errorf("Expected: %v, Got: %v", "/Users/alex/Downloads/Google Chrome 16.0 [Win 32 & 64 Bit] Full Vesrion - {RedDragon}/installer/chrome installer 16.0 .exe", file.FullPath())
  }
}

func TestPartFileExists(t *testing.T) {
  partFilePath := "/Users/alex/Downloads/Hell_in_Normandy.avi.part"
  os.Create(partFilePath)
  sfi := &SingleFileInfo {
    Name: "Hell_in_Normandy.avi",
    Length: 734160896,
  }
  file := NewFile(sfi)

  if !file.PartFileExists() {
    t.Errorf("Expected: %v, Got: %v", true, false)
  }
  os.Remove(partFilePath)
}

func TestCompletedFileExists(t *testing.T) {
  filePath := "/Users/alex/Downloads/Hell_in_Normandy.avi"
  os.Create(filePath)
  sfi := &SingleFileInfo {
    Name: "Hell_in_Normandy.avi",
    Length: 734160896,
  }
  file := NewFile(sfi)

  if !file.CompletedFileExists() {
    t.Errorf("Expected: %v, Got: %v", true, false)
  }
  os.Remove(filePath)
}

func TestCreatePartialFile(t *testing.T) {
  sfi := &SingleFileInfo {
    Name: "Hell_in_Normandy.avi",
    Length: 734160896,
  }
  file := NewFile(sfi)
  e := file.CreatePartFile()

  if e != nil {
    t.Error(e)
  }

  stat,e := os.Stat(file.FullPartPath())

  if e != nil {
    t.Error(e)
  }

  if stat.Size() != 734160896 {
    t.Errorf("Expected: %v, Got: %v", 734160896, stat.Size())
  }
  os.Remove(file.FullPartPath())
}

func TestCreatePartialFileNested(t *testing.T) {
  sfi := &SingleFileInfo {
    Name: "Google Chrome 16.0 [Win 32 & 64 Bit] Full Vesrion - {RedDragon}",
    Length: 24086584,
    Path: []string{"installer","chrome installer 16.0 .exe"},
  }
  file := NewFile(sfi)
  e := file.CreatePartFile()

  if e != nil {
    t.Error(e)
  }

  stat,_ := os.Stat(file.FullPartPath())

  if stat.Size() != 24086584 {
    t.Errorf("Expected: %v, Got: %v", 24086584, stat.Size())
  }
  os.Remove(file.FullPartPath())
}
