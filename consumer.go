package ftr

import (
	"errors"
	"io"
)

var (
	ErrMetaDeserialize = errors.New("consume: failed to deserilize metadata")
	ErrBatchCollect    = errors.New("consume: failed to collect batch with given collector")
)

type Consumer struct {
	col Collect
}

func NewTcpConsumer(col Collect) *Consumer {
	return &Consumer{
		col: col,
	}
}

func (c *Consumer) Consume(rd io.Reader, wr io.Writer, batchSize uint32) error {
	meta, err := Deserialize(rd)
	if err != nil {
		return err
	}

	err = c.col.Collect(rd, CollectOpts{
		BatchSize: int(batchSize),
		FilePath:  meta.path,
	})
	if err != nil {
		return err
	}

	return nil
}
