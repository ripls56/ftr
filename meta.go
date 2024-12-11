package ftr

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Metadata struct {
	path    string
	batchID uint32
}

type MetaOpts struct {
	Path    string
	BatchID uint32
}

func NewMetadata(opts MetaOpts) *Metadata {
	return &Metadata{
		path:    opts.Path,
		batchID: opts.BatchID,
	}
}

func Deserialize(rd io.Reader) (*Metadata, error) {
	path, err := readPath(rd)
	if err != nil {
		return nil, err
	}

	metadata := Metadata{}
	metadata.path = path

	if err := binary.Read(rd, binary.LittleEndian, &metadata.batchID); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func readPath(rd io.Reader) (string, error) {
	var pathLen uint16
	if err := binary.Read(rd, binary.LittleEndian, &pathLen); err != nil {
		if !errors.Is(err, io.EOF) {
			return "", err
		}
	}

	path := make([]rune, pathLen)
	if err := binary.Read(rd, binary.LittleEndian, path); err != nil {
		if !errors.Is(err, io.EOF) {
			return "", err
		}
	}

	return string(path), nil
}

func (fm *Metadata) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	runes := []rune(fm.path)
	if err := binary.Write(&buf, binary.LittleEndian, uint16(len(runes))); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, runes); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, fm.batchID); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
