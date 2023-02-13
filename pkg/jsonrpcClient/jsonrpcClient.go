package jsonrpcClient

import (
	"net/rpc"
	"net/rpc/jsonrpc"
)

func ConnetJsonrpc(addr string) (*rpc.Client, error) {
	conn, err := jsonrpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
