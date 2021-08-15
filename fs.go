package bundle

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type (
	File struct {
		*bytes.Reader
		offset int64
		fi     FileInfo
		file   *os.File
	}

	FileInfo struct {
		name    string
		size    int64
		modTime int64
		at      int64
		offset  int
	}
)

func (ab *AssetBundle) Open(name string) (http.File, error) {
	_name := name
	if _name == "/" {
		_name += "index.html"
	}
	_name = filepath.Clean(_name)
	if _name[0] == '\\' {
		_name = _name[1:]
	}
	info, ok := ab.fileMap[_name]
	if !ok {
		return nil, os.ErrNotExist
	}
	fileinfo := FileInfo{
		name:    filepath.Base(name),
		modTime: info.ModTime,
		size:    info.Size,
		at:      info.At,
		offset:  ab.offset,
	}
	return &File{
		fi:   fileinfo,
		file: ab.file,
	}, nil
}

func (f *File) Stat() (os.FileInfo, error) {
	return &f.fi, nil
}

func (f *File) Readdir(_ int) ([]os.FileInfo, error) {
	return nil, os.ErrNotExist
}

func (f *File) Close() error {
	return nil
}

func (f *File) Read(b []byte) (int, error) {
	if f.offset >= int64(f.fi.size) {
		return 0, io.EOF
	}
	if f.offset < 0 {
		return 0, &fs.PathError{Op: "read", Path: f.fi.name, Err: fs.ErrInvalid}
	}
	buf := make([]byte, f.fi.size-f.offset)
	n, err := f.file.ReadAt(buf, f.fi.at+int64(f.fi.offset)+f.offset)
	if err != nil {
		return 0, err
	}
	n = copy(b, buf)
	f.offset += int64(n)
	return n, nil
}

func (f *FileInfo) Name() string {
	return f.name
}

func (f *FileInfo) Size() int64 {
	return f.size
}

func (f *FileInfo) ModTime() time.Time {
	return time.Unix(f.modTime, 0)
}

func (f *FileInfo) Mode() os.FileMode {
	return 0444
}

func (f *FileInfo) IsDir() bool {
	if f.name == "\\" {
		return true
	}
	return false
}

func (f *FileInfo) Sys() interface{} {
	return nil
}
