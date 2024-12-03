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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
			splitter:  &FileSplitter{},
			filePath:  "testdata/small.txt",
			batchSize: 1,
		},
		{
			name:      "Small file",
			splitter:  &FileSplitter{},
			filePath:  "testdata/small.txt",
			batchSize: 5,
		},
		{
			name: "Small file",
			// 1 kbit
			splitter: &FileSplitter{},
			filePath: "testdata/small.txt",
			// 1 kbit
			batchSize: 1 << 10,
		},
		{
			name:     "Small file",
			splitter: &FileSplitter{},
			filePath: "testdata/small.txt",
			// 0,5 kb
			batchSize: 1 << 12,
		},
		{
			name:      "Large jpg image",
			splitter:  &FileSplitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: 128,
		},
		{
			name:     "Large jpg image",
			splitter: &FileSplitter{},
			filePath: "testdata/large-img.jpg",
			// 0,5 kb
			batchSize: 1 << 12,
		},
		{
			name:     "Large jpg image",
			splitter: &FileSplitter{},
			filePath: "testdata/large-img.jpg",
			// 32 kb
			batchSize: 1 << 18,
		},
		{
			name:     "Large jpg image",
			splitter: &FileSplitter{},
			filePath: "testdata/large-img.jpg",
			// 128 kb
			batchSize: 1 << 20,
		},
		{
			name:     "Large jpg image",
			splitter: &FileSplitter{},
			filePath: "testdata/large-img.jpg",
			// 1 mbit
			batchSize: 1 << 23,
		},
		{
			name:     "Large jpg image",
			splitter: &FileSplitter{},
			filePath: "testdata/large-img.jpg",
			// 2 mb
			batchSize: 1 << 24,
		},
	}

	for _, bm := range benchmarks {
		name := fmt.Sprintf("%s batch size %d", bm.name, bm.batchSize)

		b.Run(name, func(b *testing.B) {
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
				err = bm.splitter.Split(file, outFile, bm.batchSize)
				assert.NoError(b, err)

				outInfo, err := outFile.Stat()
				assert.NoError(b, err)

				assert.Equal(b, origInfo.Size(), outInfo.Size())
			}
		})
	}
}
