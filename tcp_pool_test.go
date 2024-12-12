package ftr

import (
	"net"
	"testing"
)

func TestWithPoolSize(t *testing.T) {
	tests := []struct {
		name        string
		size        int
		expectPanic bool
	}{
		{"Valid size", 5, false},
		{"Zero size", 0, true},
		{"Negative size", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r != nil) != tt.expectPanic {
					t.Errorf("Expected panic: %v, got: %v", tt.expectPanic, r != nil)
				}
			}()

			_ = NewTcpPool(WithPoolSize(tt.size))
		})
	}
}

func TestNewTcpPool(t *testing.T) {
	tests := []struct {
		name         string
		opts         []TcpPoolOpts
		expectedSize int
		expectedAddr string
	}{
		{
			"Default options",
			nil,
			5,
			":8080",
		},
		{
			"Custom pool size",
			[]TcpPoolOpts{
				WithPoolSize(10),
			},
			10,
			":8080",
		},
		{
			"Custom address",
			[]TcpPoolOpts{
				WithAddress("127.0.0.1:9090"),
			},
			5,
			"127.0.0.1:9090",
		},
		{
			"Custom size and address",
			[]TcpPoolOpts{
				WithPoolSize(8),
				WithAddress("localhost:6060"),
			},
			8,
			"localhost:6060",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewTcpPool(tt.opts...)

			if cap(pool.pool) != tt.expectedSize {
				t.Errorf("Expected pool size: %d, got: %d", tt.expectedSize, cap(pool.pool))
			}

			if pool.address != tt.expectedAddr {
				t.Errorf("Expected address: %s, got: %s", tt.expectedAddr, pool.address)
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		expectError bool
	}{
		{"Valid address", ":8080", false},
		{"Invalid address", "invalid-address:-1", true},
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer listener.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewTcpPool(WithAddress(tt.address), WithPoolSize(1))
			err := pool.Init()

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	addr := ":9090"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer listener.Close()
	// TODO: fix this test
	// go func() {
	// 	for {
	// 		conn, _ := listener.Accept()
	// 		go func(c net.Conn) {
	// 			defer c.Close()
	// 			buf := make([]byte, 1024)
	// 			c.Read(buf)
	// 		}(conn)
	// 	}
	// }()

	pool := NewTcpPool(WithAddress(addr), WithPoolSize(1))
	if err := pool.Init(); err != nil {
		t.Fatalf("Failed to initialize pool: %v", err)
	}
	defer pool.Close()

	tests := []struct {
		name        string
		data        []byte
		expectError bool
	}{
		{"Valid write", []byte("test data"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pool.Write(tt.data)
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestClose(t *testing.T) {
	pool := NewTcpPool(WithPoolSize(2), WithAddress(":9090"))
	addr := ":9090"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer listener.Close()

	err = pool.Init()
	if err != nil {
		t.Fatalf("Failed to initialize pool: %v", err)
	}

	pool.Close()

	if len(pool.pool) != 0 {
		t.Errorf("Expected pool to be empty after close, got size: %d", len(pool.pool))
	}
}
