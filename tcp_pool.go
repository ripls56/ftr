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

type TcpPoolOpts struct {
	PoolSize int
	Network  string
	Address  string
}

func NewTcpPool(opts TcpPoolOpts) *TcpPool {
	if opts.PoolSize <= 0 {
		panic("tcp pool: pool size should be greater then zero")
	}

	pool := make(chan net.Conn, opts.PoolSize)
	return &TcpPool{
		n:       opts.PoolSize,
		pool:    pool,
		network: opts.Network,
		address: opts.Address,
	}
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
	type result struct {
		N   int
		Err error
	}
	resCh := make(chan result)

	go func() {
		conn := <-tp.pool
		defer func() {
			tp.pool <- conn
		}()

		n, err = conn.Write(p)
		if err != nil {
			resCh <- result{
				N:   n,
				Err: err,
			}
		}

		resCh <- result{
			N:   n,
			Err: nil,
		}
	}()

	res := <-resCh
	return res.N, res.Err
}

func (tp *TcpPool) Close() {
	for len(tp.pool) > 0 {
		conn := <-tp.pool
		conn.Close()
	}
}
