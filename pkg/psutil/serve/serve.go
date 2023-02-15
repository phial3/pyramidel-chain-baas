package psutilserve

import (
	"errors"
	"flag"
	"fmt"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/jsonrpcServer"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/check"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/localcache"
	"log"
	"time"
)

type Host struct {
}

func (h *Host) Get(req interface{}, info *check.HostInfo) error {
	var hostInfo check.HostInfo
	var err error = nil
	val, ok := localcache.Cache.Get("info")
	if !ok {
		hostInfo, err = check.CheckHost()
		if err != nil {
			return err
		}
		*info = hostInfo
		return err
	} else {
		hostInfo, ok = val.(check.HostInfo)
		if !ok {
			err = errors.New("check type error")
		}
		*info = hostInfo
		return err
	}

}

func Serve() {

	var port = flag.Int("port", 8082, "jsonrpc use port default 8081")
	var interval = flag.Int("interval", 10, "cache cleanup interval seconds")
	flag.Parse()
	inter := (time.Duration)(*interval)
	localcache.Cache = localcache.NewCache(100*time.Minute, 100*time.Minute)
	var address = fmt.Sprintf("0.0.0.0:%d", *port)
	go func() {
		for {
			if _, err := check.CheckHost(); err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second * inter)
		}
	}()
	jsonrpcServer.Runserve(new(Host), address)
}
