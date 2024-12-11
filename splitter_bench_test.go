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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_        = iota
	KB int64 = 1 << (10 * iota)
	MB
)

func BenchmarkSplit(b *testing.B) {
	benchmarks := []struct {
		name      string
		splitter  Split
		filePath  string
		batchSize int
	}{
		{
			name:      "Small file",
			splitter:  &Splitter{},
			filePath:  "testdata/small.txt",
			batchSize: 1,
		},
		{
			name:      "Small file",
			splitter:  &Splitter{},
			filePath:  "testdata/small.txt",
			batchSize: 5,
		},
		{
			name:      "Small file",
			splitter:  &Splitter{},
			filePath:  "testdata/small.txt",
			batchSize: int(KB),
		},
		{
			name:      "Small file",
			splitter:  &Splitter{},
			filePath:  "testdata/small.txt",
			batchSize: int(4 * KB),
		},
		{
			name:      "Small file",
			splitter:  &Splitter{},
			filePath:  "testdata/small.txt",
			batchSize: int(32 * KB),
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: 128,
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: int(KB),
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: int(4 * KB),
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: int(16 * KB),
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: int(32 * KB),
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: int(64 * KB),
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: int(256 * KB),
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: int(MB),
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: int(2 * MB),
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: int(4 * MB),
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: int(8 * MB),
		},
		{
			name:      "Large jpg image",
			splitter:  &Splitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: int(16 * MB),
		},
	}

	for _, bm := range benchmarks {
		name := fmt.Sprintf("%s batch size %d", bm.name, bm.batchSize)

		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				file, err := os.Open(bm.filePath)
				assert.NoError(b, err)
				defer file.Close()

				fileName := filepath.Base(bm.filePath)
				outFile, err := os.CreateTemp(os.TempDir(), fileName)
				assert.NoError(b, err)

				defer outFile.Close()
				defer os.Remove(outFile.Name())

				origInfo, err := file.Stat()
				assert.NoError(b, err)
				err = bm.splitter.Split(file, outFile, SplitOpts{
					BatchSize: bm.batchSize,
					FilePath:  bm.filePath,
				})
				if err != io.EOF {
					assert.NoError(b, err)
				}

				outInfo, err := outFile.Stat()
				assert.NoError(b, err)

				assert.Equal(b, origInfo.Size(), outInfo.Size())
			}
		})
	}
}
