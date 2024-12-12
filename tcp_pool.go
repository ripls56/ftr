package ftr

import (
	"net"
)

type TcpPool struct {
	n    int // Max number of connections
	pool chan net.Conn

	network string
	address string
}

type TcpPoolOpts func(*TcpPool)

func WithPoolSize(size int) TcpPoolOpts {
	return func(tp *TcpPool) {
		if size <= 0 {
			panic("tcp pool: pool size should be greater then zero")
		}

		tp.n = size
	}
}

func WithAddress(address string) TcpPoolOpts {
	return func(tp *TcpPool) {
		tp.address = address
	}
}

func NewTcpPool(opts ...TcpPoolOpts) *TcpPool {
	tcpPool := TcpPool{
		n:       5,
		network: "tcp",
		address: ":8080",
	}

	for _, opt := range opts {
		opt(&tcpPool)
	}

	pool := make(chan net.Conn, tcpPool.n)
	tcpPool.pool = pool

	return &tcpPool
}

func (tp *TcpPool) Init() error {
	for range tp.n {
		conn, err := net.Dial(tp.network, tp.address)
		if err != nil {
			return err
		}
		tp.pool <- conn
	}

	return nil
}

func (tp *TcpPool) Write(p []byte) (n int, err error) {
	conn := <-tp.pool
	defer func() {
		tp.pool <- conn
	}()

	n, err = conn.Write(p)
	if err != nil {
		return n, err
	}

	return n, err
}

func (tp *TcpPool) Close() {
	for len(tp.pool) > 0 {
		conn := <-tp.pool
		conn.Close()
	}
}
