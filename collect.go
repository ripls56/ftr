package ftr

import (
	"io"
	"os"
)

type Collect interface {
	Collect(rd io.Reader, opts CollectOpts) error
}

type CollectOpts struct {
	BatchSize int
	FileName  string
}

type FileCollector struct {
}

func (co *FileCollector) Collect(rd io.Reader, opts CollectOpts) error {
	if opts.BatchSize <= 0 {
		// TODO: imrpv errors
		return io.ErrShortBuffer
	}

	file, err := os.Create(opts.FileName)
	if err != nil {
		return err
	}
	defer file.Close()
	defer file.Sync()

	chunk := make([]byte, opts.BatchSize)

	if _, err := io.CopyBuffer(file, rd, chunk); err != nil {
		return err
	}

	return nil
}
