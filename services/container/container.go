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
	Host     string
	PublicIP string
	cli      *client.Client
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

func NewContainerService(host, ip string) ContainerService {
	return &containerService{
		Host:     host,
		PublicIP: ip,
	}
}

func (s *containerService) Conn() error {
	address := fmt.Sprintf("tcp://%s:%d", s.Host, 2376)
	cacertPath := fmt.Sprintf("/root/txhyjuicefs/%s/certs/ca.pem", s.PublicIP)
	certPath := fmt.Sprintf("/root/txhyjuicefs/%s/certs/client.pem", s.PublicIP)
	keyPath := fmt.Sprintf("/root/txhyjuicefs/%s/certs/client-key.pem", s.PublicIP)
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(),
		client.WithHost(address), client.WithTLSClientConfig(cacertPath, certPath, keyPath))
	if err != nil {
		return err
	}
	s.cli = cli
	return nil
}

func (s *containerService) Run(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) error {
	createRes, err := s.cli.ContainerCreate(ctx, config, hostConfig, networkingConfig, nil, "")
	if err != nil {
		return err
	}
	if err := s.cli.ContainerStart(ctx, createRes.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func (s *containerService) Login(ctx context.Context) error {
	return nil
}
