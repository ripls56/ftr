// Copyright 2024 ripls56
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ftr

import (
	"io"
)

type Split interface {
	Split(rd io.Reader, wr io.Writer, opts SplitOpts) error
}

type SplitOpts struct {
	BatchSize int
	Filepath  string
}

type FileSplitter struct{}

func NewFileSplitter(ser Serialize) Split {
	return &FileSplitter{}
}

func (sp *FileSplitter) Split(rd io.Reader, wr io.Writer, opts SplitOpts) error {
	if opts.BatchSize <= 0 {
		return io.ErrShortBuffer
	}

	batch := make([]byte, opts.BatchSize)

	errCh := make(chan error)
	go func() {
		for {
			n, err := rd.Read(batch)
			if err != nil {
				errCh <- err
			}

			if n > 0 {
				// TODO: write metadata first
				if _, err := wr.Write(batch[:n]); err != nil {
					errCh <- err
				}
				continue
			}

			errCh <- io.EOF
		}
	}()

	return <-errCh
}
