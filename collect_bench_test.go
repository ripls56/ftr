package ftr

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkCollect(b *testing.B) {
	benchmarks := []struct {
		name      string
		collector Collect
		filePath  string
		batchSize int
	}{
		{
			name:      "Small file",
			collector: &FileCollector{},
			filePath:  "testdata/small.txt",
			batchSize: 1,
		},
		{
			name:      "Small file",
			collector: &FileCollector{},
			filePath:  "testdata/small.txt",
			batchSize: 5,
		},
		{
			name:      "Small file",
			collector: &FileCollector{},
			filePath:  "testdata/small.txt",
			batchSize: 1000,
		},
		{
			name:      "Large jpg image",
			collector: &FileCollector{},
			filePath:  "testdata/large-img.jpg",
			batchSize: 128,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			file, err := os.Open(bm.filePath)
			assert.NoError(b, err)
			defer file.Close()

			var buf bytes.Buffer
			splitter := FileSplitterV2{}
			err = splitter.Split(file, &buf, SplitOpts{
				BatchSize: bm.batchSize,
				FileName:  filepath.Base(bm.filePath),
			})
			assert.NoError(b, err)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err = bm.collector.Collect(&buf, CollectOpts{
					BatchSize: bm.batchSize,
					FileName:  file.Name(),
				})
				assert.NoError(b, err)
			}
		})
	}
}
