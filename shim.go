package airfreight

import (
	"net/http"
	"os"
	"time"
)

type entFile struct {
	name     string
	modTime  int64
	contents []byte
	offset   int
}

// FileSystem implements the http.FileSystem interface so Airfreight
// can be plugged directly into a Go net/http web server.
type FileSystem struct {
	files map[string]Ent
}

// HTTPFileSystem makes a http.FileSystem compatible adapter from an
// Airfreight map.
func HTTPFileSystem(files map[string]Ent) FileSystem {
	return FileSystem{files: files}
}

// Open a file for reading from the Airfreight map.
func (ifs FileSystem) Open(name string) (http.File, error) {
	ent, exists := ifs.files[name]
	if !exists {
		return nil, os.ErrNotExist
	}
	return &entFile{
		name:     name,
		modTime:  ent.ModTime,
		contents: []byte(ent.Contents),
	}, nil
}

func (entFile) Close() error {
	return nil
}

func (efile *entFile) Read(p []byte) (n int, err error) {
	i := efile.offset
	n = len(efile.contents) - i
	if n > cap(p) {
		n = cap(p)
	}
	copy(p, efile.contents[i:i+n])
	efile.offset = i + n
	return n, nil
}

func (efile *entFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		efile.offset = int(offset)
	case 1:
		efile.offset += int(offset)
	case 2:
		efile.offset = len(efile.contents) + int(offset)
	}
	return int64(efile.offset), nil
}

func (efile entFile) Readdir(count int) ([]os.FileInfo, error) {
	return []os.FileInfo{}, nil
}

func (efile entFile) Stat() (os.FileInfo, error) {
	return efile, nil
}

func (efile entFile) Name() string {
	return efile.name
}

func (efile entFile) Size() int64 {
	return int64(len(efile.contents))
}

func (efile entFile) Mode() os.FileMode {
	return 0444
}

func (efile entFile) ModTime() time.Time {
	return time.Unix(efile.modTime, 0)
}

func (efile entFile) IsDir() bool {
	return false
}

func (efile entFile) Sys() interface{} {
	return nil
}
