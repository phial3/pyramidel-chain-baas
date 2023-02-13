package serve

import (
	"flag"
	"fmt"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/jsonrpcServer"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/check"
)

type Host struct {
}

func (h *Host) Get(req interface{}, info *check.HostInfo) error {
	hostInfo, err := check.CheckHost()
	if err != nil {
		return err
	}
	*info = hostInfo
	return nil
}

func Serve() {
	var port = flag.Int("port", 8082, "jsonrpc use port default 8081")
	flag.Parse()
	var address = fmt.Sprintf("0.0.0.0:%d", *port)
	jsonrpcServer.Runserve(new(Host), address)
}
