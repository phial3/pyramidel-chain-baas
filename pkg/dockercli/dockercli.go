package dockercli

import (
	"github.com/docker/docker/client"
)

func Connect(host string) (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), client.WithHost(host))
	if err != nil {
		return nil, err
	}
	return cli, err
}
