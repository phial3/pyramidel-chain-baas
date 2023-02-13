package psutilclient

import (
	"github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/check"
	"net/rpc"
)

func CallPsutil(client *rpc.Client) (check.HostInfo, error) {
	hostinfo := check.HostInfo{}
	if err := client.Call("Host.Get", nil, &hostinfo); err != nil {
		return hostinfo, err
	}
	return hostinfo, nil
}
