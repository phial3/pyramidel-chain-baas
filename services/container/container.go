package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var _ ContainerService = (*containerService)(nil)

type ContainerService interface {
	Run(context.Context, *container.Config, *container.HostConfig, *network.NetworkingConfig, string) error
	Pull(context.Context, string) error
	Login(context.Context) error
	List(context.Context) ([]types.Container, error)
	Close() error
	Conn() error
}

type containerService struct {
	Host string
	cli  *client.Client
}

func (s *containerService) Pull(ctx context.Context, s2 string) error {
	//TODO implement me
	panic("implement me")
}

func (s *containerService) List(ctx context.Context) ([]types.Container, error) {
	return s.cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
}

func (s *containerService) Close() error {
	return s.cli.Close()
}

func NewContainerService(host string) ContainerService {
	return &containerService{
		Host: host,
	}
}

func (s *containerService) Conn() error {
	address := fmt.Sprintf("tcp://%s:%d", s.Host, 2375)
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(),
		client.WithHost(address))
	if err != nil {
		return err
	}
	s.cli = cli
	return nil
}

func (s *containerService) Run(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) error {
	s.cli.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, "")
	return nil
}

func (s *containerService) Login(ctx context.Context) error {
	return nil
}
