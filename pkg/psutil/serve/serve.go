package psutilserve

import (
	"errors"
	"flag"
	"fmt"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/jsonrpcServer"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/check"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/localcache"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/freeport"
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

type Response struct {
	Port int `json:"port"`
}

func (h *Host) GetPort(req interface{}, res *Response) error {
	result := Response{}
	port, err := freeport.GetFreePort()
	if err != nil {
		*res = result
		return err
	}
	log.Println(port)
	*res = result
	res.Port = port
	return nil
}

func Serve() {
	var port = flag.Int("port", 8082, "jsonrpc use port default 8082")
	var interval = flag.Int("interval", 10, "cache cleanup interval seconds")
	flag.Parse()
	inter := (time.Duration)(*interval)
	localcache.Cache = localcache.NewCache(100*time.Minute, 100*time.Minute)
	var address = fmt.Sprintf("0.0.0.0:%d", *port)
	go func() {
		for {
			if host, err := check.CheckHost(); err != nil {
				localcache.Cache.Set("info", host, (inter+5)*time.Second)
				log.Println(err)
			}
			time.Sleep(time.Second * inter)
		}
	}()
	jsonrpcServer.Runserve(new(Host), address)
}
