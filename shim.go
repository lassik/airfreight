package airfreight

import (
	"net/http"
	"os"
	"time"
)

type EntFile struct {
	name     string
	modTime  int64
	contents []byte
	offset   int
}

type FileSystem struct {
	files map[string]Ent
}

func MapFileSystem(files map[string]Ent) FileSystem {
	return FileSystem{files: files}
}

func (ifs FileSystem) Open(name string) (http.File, error) {
	ent, exists := ifs.files[name]
	if !exists {
		return nil, os.ErrNotExist
	}
	return &EntFile{name: name, contents: []byte(ent.Contents)}, nil
}

func (EntFile) Close() error {
	return nil
}

func (entFile *EntFile) Read(p []byte) (n int, err error) {
	i := entFile.offset
	n = len(entFile.contents) - i
	if n > cap(p) {
		n = cap(p)
	}
	copy(p, entFile.contents[i:i+n])
	entFile.offset = i + n
	return n, nil
}

func (entFile EntFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		entFile.offset = int(offset)
	case 1:
		entFile.offset += int(offset)
	case 2:
		entFile.offset = len(entFile.contents) + int(offset)
	}
	return int64(entFile.offset), nil
}

func (entFile EntFile) Readdir(count int) ([]os.FileInfo, error) {
	return []os.FileInfo{}, nil
}

func (entFile EntFile) Stat() (os.FileInfo, error) {
	return entFile, nil
}

func (entFile EntFile) Name() string {
	return entFile.name
}

func (entFile EntFile) Size() int64 {
	return int64(len(entFile.contents))
}

func (entFile EntFile) Mode() os.FileMode {
	return 0644
}

func (entFile EntFile) ModTime() time.Time {
	return time.Unix(entFile.modTime, 0)
}

func (entFile EntFile) IsDir() bool {
	return false
}

func (entFile EntFile) Sys() interface{} {
	return nil
}
