package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
)

var _ ContainerService = (*caService)(nil)

type ContainerService interface {
	Run(context.Context, *container.Config, *container.HostConfig, string) error
	Pull(context.Context, string) error
	Login(context.Context) error
	List(context.Context) ([]types.Container, error)
	Close() error
	Conn() error
	GenConfig(context.Context) (*container.Config, *container.HostConfig)
	SetNetwork(context.Context, string) error
}

type caService struct {
	Host         string
	PublicIP     string
	Port         string
	caUser       string
	caPw         string
	orgUscc      string
	serverName   string
	serverDomain string
	cli          *client.Client
}

func (s *caService) Login(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *caService) Pull(ctx context.Context, s2 string) error {
	//TODO implement me
	panic("implement me")
}

func (s *caService) List(ctx context.Context) ([]types.Container, error) {
	return s.cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
}

func (s *caService) Close() error {
	return s.cli.Close()
}

func NewCaContainerService(host, ip, port, caUser, caPw, orgUscc, serverName, serverDomain string) ContainerService {
	return &caService{
		Host:         host,
		Port:         port,
		PublicIP:     ip,
		caPw:         caPw,
		caUser:       caUser,
		orgUscc:      orgUscc,
		serverName:   serverName,
		serverDomain: serverDomain,
	}
}

func (s *caService) Conn() error {
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

func (s *caService) Run(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, containerName string) error {
	createRes, err := s.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}

	if err := s.SetNetwork(ctx, createRes.ID); err != nil {
		return err
	}
	log.Println(config, hostConfig)
	if err := s.cli.ContainerStart(ctx, createRes.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func (s *caService) GenConfig(ctx context.Context) (*container.Config, *container.HostConfig) {
	port := fmt.Sprintf("%s/tcp", s.Port)
	portenv := fmt.Sprintf("FABRIC_CA_SERVER_PORT=%s", s.Port)
	binddir := fmt.Sprintf("/root/txhyjuicefs/organizations/fabric-ca/%s:/etc/hyperledger/fabric-ca-server", s.orgUscc)
	user_pass := fmt.Sprintf("BOOTSTRAP_USER_PASS=%s:%s", s.caUser, s.caPw)
	ca_server_name := fmt.Sprintf("FABRIC_CA_SERVER_CA_NAME=%s", s.serverName)
	containerConfig := &container.Config{Image: "harbor.sxtxhy.com/gcbaas-gm/fabric-ca:latest", Cmd: []string{"sh", "-c", "/usr/local/bin/start_ca.sh"}, Env: []string{"FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server",
		ca_server_name, "FABRIC_CA_SERVER_TLS_ENABLED=true", portenv, user_pass}, ExposedPorts: map[nat.Port]struct{}{nat.Port(port): struct {
	}{}}, Hostname: s.serverName, Domainname: s.serverDomain}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(port): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: s.Port}},
		},
		Binds: []string{binddir},
		Resources: container.Resources{
			CPUShares: 1024,
			Memory:    536870912,
		},
	}

	return containerConfig, hostConfig
}

func (s *caService) SetNetwork(ctx context.Context, containerId string) error {
	// 获取网络
	networkName := "fabric_network"
	fabric_network, err := s.cli.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{})
	if err != nil {
		return err
	}

	endpointConfig := &network.EndpointSettings{
		NetworkID: fabric_network.ID,
	}
	err = s.cli.NetworkConnect(ctx, fabric_network.ID, containerId, endpointConfig)
	if err != nil {
		return err
	}
	return nil
}
