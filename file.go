package main

import (
	"os"
	"strings"
)

type File struct {
	Filename string
	Path     string
	Length   int
	Md5Sum   string
}

func NewFile(sfi *SingleFileInfo) *File {
	var filename string
	var path string
	if len(sfi.Path) > 0 {
		filename = sfi.Path[len(sfi.Path)-1]
		var pathParts []string
		pathParts = append(pathParts, downloadPath)
		pathParts = append(pathParts, sfi.Name)
		if len(sfi.Path) > 1 {
			pathParts = append(pathParts, sfi.Path[0:len(sfi.Path)-1]...)
		}
		path = strings.Join(pathParts, "/")
	} else {
		filename = sfi.Name
		path = downloadPath
	}
	length := sfi.Length
	return &File{
		Filename: filename,
		Path:     path,
		Length:   length,
	}
}

func (f *File) FullPath() string {
	fullpath := []string{f.Path, f.Filename}
	return strings.Join(fullpath, "/")
}

func (f *File) FullPartPath() string {
	return f.FullPath() + ".part"
}

func (f *File) PartFileExists() bool {
	_, e := os.Stat(f.FullPartPath())
	if os.IsNotExist(e) {
		return false
	} else {
		return true
	}
}

func (f *File) CompletedFileExists() bool {
	_, e := os.Stat(f.FullPath())
	if os.IsNotExist(e) {
		return false
	} else {
		return true
	}
}

func (f *File) CreatePartFile() (err error) {
	os.MkdirAll(f.Path, 0755)
	of, err := os.Create(f.FullPartPath())
	if err != nil {
		return
	}
	err = of.Truncate(int64(f.Length))
	return
}

func (f *File) OpenPartFileWrite() (file *os.File, err error) {
	file, err = os.OpenFile(f.FullPartPath(), os.O_WRONLY, 0666)
	return
}
