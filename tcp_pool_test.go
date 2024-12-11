package ftr

import (
	"net"
	"testing"
)

func TestTcpPool_Write(t *testing.T) {
	type fields struct {
		n       int
		pool    chan net.Conn
		network string
		address string
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		{
			name: "Normal pool",
			fields: fields{
				n:       5,
				pool:    make(chan net.Conn, 5),
				network: "tcp",
				address: ":8080",
			},
			args: args{
				p: []byte("asd"),
			},
			wantN:   3,
			wantErr: false,
		},
	}

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &TcpPool{
				n:       tt.fields.n,
				pool:    tt.fields.pool,
				network: tt.fields.network,
				address: tt.fields.address,
			}
			defer tp.Close()

			if err := tp.Init(); (err != nil) != tt.wantErr {
				t.Errorf("TcpPool.Init() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotN, err := tp.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("TcpPool.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("TcpPool.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
