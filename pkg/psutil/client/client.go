package psutilclient

import (
	"github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/check"
	psutilserve "github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/serve"
	"net/rpc"
)

func CallPsutil(client *rpc.Client) (check.HostInfo, error) {
	hostinfo := check.HostInfo{}
	if err := client.Call("Host.Get", nil, &hostinfo); err != nil {
		return hostinfo, err
	}
	return hostinfo, nil
}

func CallGetPort(client *rpc.Client) (int, error) {
	res := psutilserve.Response{}
	if err := client.Call("Host.GetPort", nil, &res); err != nil {
		return res.Port, err
	}
	return res.Port, nil
}
