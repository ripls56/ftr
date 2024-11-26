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
	"os"
	"path/filepath"
	"testing"
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
			name:      "Small file",
			splitter:  &FileSplitter{},
			filePath:  "testdata/small.txt",
			batchSize: 1000,
		},
		{
			name:      "Large jpg image",
			splitter:  &FileSplitter{},
			filePath:  "testdata/large-img.jpg",
			batchSize: 128,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				file, err := os.Open(bm.filePath)
				if err != nil {
					b.Fatalf("Failed to open file: %v", err)
				}
				defer file.Close()

				fileName := filepath.Base(bm.filePath)
				outFile, err := os.CreateTemp(os.TempDir(), fileName)
				if err != nil {
					b.Fatalf("Failed to create temp file: %v", err)
				}
				defer outFile.Close()
				defer os.Remove(outFile.Name())

				ch := bm.splitter.Split(file, SplitOpts{
					BatchSize: bm.batchSize,
					FileName:  fileName,
				})

				for batch := range ch {
					_, err := outFile.Write(batch.Content)
					if err != nil {
						b.Fatalf("Failed to write batch: %v", err)
					}
				}
			}
		})
	}
}

func BenchmarkSplitV2(b *testing.B) {
	benchmarks := []struct {
		name      string
		splitter  SplitV2
		filePath  string
		batchSize int
	}{
		{
			name:      "Small file",
			splitter:  &FileSplitterV2{},
			filePath:  "testdata/small.txt",
			batchSize: 1,
		},
		{
			name:      "Small file",
			splitter:  &FileSplitterV2{},
			filePath:  "testdata/small.txt",
			batchSize: 5,
		},
		{
			name:      "Small file",
			splitter:  &FileSplitterV2{},
			filePath:  "testdata/small.txt",
			batchSize: 1000,
		},
		{
			name:      "Large jpg image",
			splitter:  &FileSplitterV2{},
			filePath:  "testdata/large-img.jpg",
			batchSize: 128,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				file, err := os.Open(bm.filePath)
				if err != nil {
					b.Fatalf("Failed to open file: %v", err)
				}
				defer file.Close()

				fileName := filepath.Base(bm.filePath)
				outFile, err := os.CreateTemp(os.TempDir(), fileName)
				if err != nil {
					b.Fatalf("Failed to create temp file: %v", err)
				}
				defer outFile.Close()
				defer os.Remove(outFile.Name())

				err = bm.splitter.Split(file, outFile, SplitOpts{
					BatchSize: bm.batchSize,
					FileName:  fileName,
				})
				if err != nil {
					b.Fatalf("Failed to split file: %v", err)
				}
			}
		})
	}
}
