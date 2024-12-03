package ftr

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestFileMeta_Serialize(t *testing.T) {
	type fields struct {
		path    string
		batchID uint32
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Valid data - simple path",
			fields: fields{
				path:    "testdata/small.txt",
				batchID: 1,
			},
			want: func() []byte {
				var buf bytes.Buffer
				_ = binary.Write(&buf, binary.LittleEndian, uint16(len("testdata/small.txt")))
				_ = binary.Write(&buf, binary.LittleEndian, []rune("testdata/small.txt"))
				_ = binary.Write(&buf, binary.LittleEndian, uint32(1))
				return buf.Bytes()
			}(),
			wantErr: false,
		},
		{
			name: "Empty path",
			fields: fields{
				path:    "",
				batchID: 123,
			},
			want: func() []byte {
				var buf bytes.Buffer
				_ = binary.Write(&buf, binary.LittleEndian, uint16(0))
				_ = binary.Write(&buf, binary.LittleEndian, uint32(123))
				return buf.Bytes()
			}(),
			wantErr: false,
		},
		{
			name: "Long path",
			fields: fields{
				path:    "a/very/long/path/to/some/file/that/keeps/going.and.have.a.very.long.ext",
				batchID: 42,
			},
			want: func() []byte {
				var buf bytes.Buffer
				path := "a/very/long/path/to/some/file/that/keeps/going.and.have.a.very.long.ext"
				_ = binary.Write(&buf, binary.LittleEndian, uint16(len(path)))
				_ = binary.Write(&buf, binary.LittleEndian, []rune(path))
				_ = binary.Write(&buf, binary.LittleEndian, uint32(42))
				return buf.Bytes()
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm := &Metadata{
				path:    tt.fields.path,
				batchID: tt.fields.batchID,
			}
			meta, err := fm.Serialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("FileMeta.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(meta, tt.want) {
				t.Errorf("FileMeta.Serialize() = %b, want %b", meta, tt.want)
			}
		})
	}
}

func TestFileMeta_Deserialize(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *Metadata
		wantErr bool
	}{
		{
			name: "Valid data",
			data: func() []byte {
				var buf bytes.Buffer
				_ = binary.Write(&buf, binary.LittleEndian, uint16(len("testdata/small.txt")))
				_ = binary.Write(&buf, binary.LittleEndian, []rune("testdata/small.txt"))
				_ = binary.Write(&buf, binary.LittleEndian, uint32(1))
				return buf.Bytes()
			}(),
			want: &Metadata{
				path:    "testdata/small.txt",
				batchID: 1,
			},
			wantErr: false,
		},
		{
			name:    "Empty data",
			data:    []byte{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Corrupted data",
			data: func() []byte {
				var buf bytes.Buffer
				_ = binary.Write(&buf, binary.LittleEndian, uint16(5))
				_ = binary.Write(&buf, binary.LittleEndian, []rune("ab"))
				return buf.Bytes()
			}(),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := &Metadata{}
			var buf bytes.Buffer
			buf.Write(tt.data)
			err := meta.Deserialize(&buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileMeta.Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && (meta.path != tt.want.path || meta.batchID != tt.want.batchID) {
				t.Errorf("FileMeta.Deserialize() = %v, want %v", meta, tt.want)
			}
		})
	}
}

func TestFileMeta_SerializeAndDeserialize(t *testing.T) {
	tests := []struct {
		name    string
		meta    *Metadata
		wantErr bool
	}{
		{
			name: "Valid data",
			meta: &Metadata{
				path:    "testdata/small.txt",
				batchID: 123,
			},
			wantErr: false,
		},
		{
			name: "Empty path",
			meta: &Metadata{
				path:    "",
				batchID: 456,
			},
			wantErr: false,
		},
		{
			name: "Long path",
			meta: &Metadata{
				path:    "a/very/long/path/to/some/file/that/keeps/going.and.have.a.very.long.ext",
				batchID: 789,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serialized, err := tt.meta.Serialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("FileMeta.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var buf bytes.Buffer
			buf.Write(serialized)

			meta := &Metadata{}
			err = meta.Deserialize(&buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileMeta.Deserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("data: \n%+v\n%+v", tt.meta, meta)
			if tt.meta.path != meta.path || tt.meta.batchID != meta.batchID {
				t.Errorf("Serialize and Deserialize mismatch: meta = %v, want = %v", meta, tt.meta)
			}
		})
	}
}
