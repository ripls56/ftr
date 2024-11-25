package ftr

import (
	"bufio"
	"io"
)

type Splitter interface {
	Split(rd io.Reader, opts SplitOpts) <-chan Batch
}

type SplitOpts struct {
	BatchSize int
	FileName  string
}

type FileSplitter struct {
}

func (sp *FileSplitter) Split(rd io.Reader, opts SplitOpts) <-chan Batch {
	ch := make(chan Batch)

	if opts.BatchSize <= 0 {
		panic("split: batch size should be greater then zero")
	}
	go func() {
		defer close(ch)

		buf := bufio.NewReader(rd)

		id := 1

		for {
			chunk := make([]byte, opts.BatchSize)
			n, err := buf.Read(chunk)

			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}

			ch <- Batch{
				Id: id,
				Meta: FileMeta{
					FileName: opts.FileName,
				},
				Content: chunk[:n],
			}
			id++
		}
	}()
	return ch
}
