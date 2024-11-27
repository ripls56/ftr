package ftr

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// FileMeta
type FileMeta struct {
	path    string
	batchID uint32
}

type Serialize interface {
	Serialize() ([]byte, error)
}

type Deserialize interface {
	Deserialize(data []byte) (*FileMeta, error)
}

type FileMetaOpts struct {
	Path    string
	BatchID uint32
}

func NewMeta(opts FileMetaOpts) *FileMeta {
	return &FileMeta{
		path:    opts.Path,
		batchID: opts.BatchID,
	}
}

func (fm *FileMeta) Deserialize(data []byte) (*FileMeta, error) {
	if len(data) <= 0 {
		return nil, errors.New("deserialize: buffer is empty")
	}
	buf := bytes.NewReader(data)
	meta := &FileMeta{}

	path, err := readPath(buf)
	if err != nil {
		return nil, err
	}
	meta.path = path

	if err := binary.Read(buf, binary.BigEndian, &meta.batchID); err != nil {
		return nil, err
	}

	return meta, nil
}

func readPath(buf io.Reader) (string, error) {
	var pathLen uint16
	if err := binary.Read(buf, binary.BigEndian, &pathLen); err != nil {
		if !errors.Is(err, io.EOF) {
			return "", err
		}
	}

	path := make([]rune, pathLen)
	if err := binary.Read(buf, binary.BigEndian, path); err != nil {
		if !errors.Is(err, io.EOF) {
			return "", err
		}
	}

	return string(path), nil
}

func (fm *FileMeta) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	runes := []rune(fm.path)
	if err := binary.Write(&buf, binary.BigEndian, uint16(len(runes))); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, runes); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, fm.batchID); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
