package jsonrpcServer

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func Runserve(rcvr interface{}, address string) {
	err := rpc.Register(rcvr)
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := l.Close(); err != nil {
			panic(err)
		}
	}()

	for {
		conn, e := l.Accept()
		if e != nil {
			continue
		}
		defer func() {
			if err := conn.Close(); err != nil {
				panic(err)
			}
		}()
		go jsonrpc.ServeConn(conn)
	}
}
