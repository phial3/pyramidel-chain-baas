package jsonrpcClient

import (
	"github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/check"
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

func CallPsutil(client *rpc.Client) (check.HostInfo, error) {
	hostinfo := check.HostInfo{}
	if err := client.Call("Host.Get", nil, &hostinfo); err != nil {
		return hostinfo, err
	}
	return hostinfo, nil
}
