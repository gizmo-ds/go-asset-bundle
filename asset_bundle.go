package bundle

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// FI 文件头部标识
var FI = []byte("GoAB")

type (
	// AssetBundle 资源包
	AssetBundle struct {
		Size    int64
		Head    []byte
		Files   []ABFileInfo
		Version uint16
		file    *os.File
		fileMap map[string]ABFileInfo
		offset  int
	}
	// ABFileInfo 资源包文件信息
	ABFileInfo struct {
		Path    string `json:"p"`
		ModTime int64  `json:"t"`
		Size    int64  `json:"s"`
		At      int64  `json:"a"`
	}
)

// NewAssetBundle 创建一个新的资源包
func NewAssetBundle(name string) (*AssetBundle, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	ab := AssetBundle{
		file: f,
	}
	return &ab, nil
}

// Close 关闭资源包
func (ab *AssetBundle) Close() error {
	return ab.file.Close()
}

// Bundle 捆绑指定目录的资源
func (ab *AssetBundle) Bundle(root string, version uint16) error {
	var err error
	if ab.fileMap == nil {
		ab.fileMap = make(map[string]ABFileInfo)
	}
	err = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			_path, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			fileInfo := ABFileInfo{
				Path:    _path,
				ModTime: info.ModTime().Unix(),
				Size:    info.Size(),
				At:      ab.Size,
			}
			ab.Files = append(ab.Files, fileInfo)
			ab.fileMap[_path] = fileInfo
			ab.Size += info.Size()
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(ab.file, bytes.NewReader(FI))
	if err != nil {
		return err
	}

	_version := make([]byte, 2)
	binary.LittleEndian.PutUint16(_version, version)
	_, err = io.Copy(ab.file, bytes.NewReader(_version))
	if err != nil {
		return err
	}
	ab.Version = version

	ab.Head, _ = json.Marshal(ab.Files)

	headSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(headSize, uint32(len(ab.Head)))
	_, err = io.Copy(ab.file, bytes.NewReader(headSize))
	if err != nil {
		return err
	}

	_, err = io.Copy(ab.file, bytes.NewReader(ab.Head))
	if err != nil {
		return err
	}

	ab.offset = len(FI) + 2 + 4 + len(ab.Head)

	for i := 0; i < len(ab.Files); i++ {
		f, err := os.Open(filepath.Join(root, ab.Files[i].Path))
		if err != nil {
			return err
		}
		_, err = io.Copy(ab.file, f)
		if err != nil {
			return err
		}
	}
	return nil
}

// OpenAssetBundle 从文件打开资源包
func OpenAssetBundle(name string) (*AssetBundle, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	ab := AssetBundle{
		file: f,
	}

	size := len(FI) + 2 // 标识加上版本号
	buf := make([]byte, size)
	_, err = f.Read(buf)
	if err != nil {
		return nil, err
	}

	if string(buf[:len(FI)]) != string(FI) {
		return nil, errors.New("Unsupported file type")
	}

	versionData := buf[len(FI):]
	ab.Version = binary.LittleEndian.Uint16(versionData)

	buf = make([]byte, 4)
	_, err = f.Read(buf)
	if err != nil {
		return nil, err
	}
	headSize := binary.LittleEndian.Uint32(buf)

	buf = make([]byte, headSize)
	_, err = f.Read(buf)
	if err != nil {
		return nil, err
	}
	ab.Head = buf

	ab.offset = size + 4 + len(buf)

	ab.fileMap = make(map[string]ABFileInfo)
	err = json.Unmarshal(ab.Head, &ab.Files)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(ab.Files); i++ {
		ab.fileMap[ab.Files[i].Path] = ab.Files[i]
	}
	return &ab, nil
}
