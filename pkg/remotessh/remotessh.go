package remotessh

import (
	"github.com/melbahja/goph"
	"go.uber.org/zap"
	"time"
)

var sshLogger = zap.L().Named("remote/ssh")

func ConnectToHost(pass, ip string, port uint) (*goph.Client, error) {
	callback, err := goph.DefaultKnownHosts()
	if err != nil {
		return nil, err
	}
	client, err := goph.NewConn(&goph.Config{
		Auth:     goph.Password(pass),
		User:     "root",
		Addr:     ip,
		Port:     port,
		Timeout:  10 * time.Second,
		Callback: callback,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
