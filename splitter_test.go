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

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	tests := []struct {
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

	for _, tc := range tests {
		t.Logf("batch size: %d", tc.batchSize)

		t.Run(tc.name, func(t *testing.T) {
			file, err := os.Open(tc.filePath)
			assert.NoError(t, err)
			defer file.Close()

			origInfo, err := file.Stat()
			assert.NoError(t, err)

			fileName := filepath.Base(tc.filePath)
			outFile, err := os.CreateTemp(os.TempDir(), fileName)
			assert.NoError(t, err)
			defer outFile.Close()
			defer os.Remove(outFile.Name())

			ch := tc.splitter.Split(file, SplitOpts{
				BatchSize: tc.batchSize,
				FileName:  fileName,
			})

			for batch := range ch {
				_, err := outFile.Write(batch.Content)
				assert.NoError(t, err)
			}

			err = outFile.Sync()
			assert.NoError(t, err)

			reasmInfo, err := outFile.Stat()
			assert.NoError(t, err)
			assert.Equal(
				t,
				origInfo.Size(),
				reasmInfo.Size(),
				"The size of the assembled file does not match the original file",
			)
		})
	}
}

func TestSplitV2(t *testing.T) {
	tests := []struct {
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

	for _, tc := range tests {
		t.Logf("batch size: %d", tc.batchSize)

		t.Run(tc.name, func(t *testing.T) {
			file, err := os.Open(tc.filePath)
			assert.NoError(t, err)
			defer file.Close()

			origInfo, err := file.Stat()
			assert.NoError(t, err)

			fileName := filepath.Base(tc.filePath)
			outFile, err := os.CreateTemp(os.TempDir(), fileName)
			assert.NoError(t, err)
			defer outFile.Close()
			defer os.Remove(outFile.Name())

			err = tc.splitter.Split(file, outFile, SplitOpts{
				BatchSize: tc.batchSize,
				FileName:  fileName,
			})

			assert.NoError(t, err)

			err = outFile.Sync()
			assert.NoError(t, err)

			reasmInfo, err := outFile.Stat()
			assert.NoError(t, err)
			assert.Equal(
				t,
				origInfo.Size(),
				reasmInfo.Size(),
				"The size of the assembled file does not match the original file",
			)
		})
	}
}
