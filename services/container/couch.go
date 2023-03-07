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

var _ ContainerService = (*couchService)(nil)

type couchService struct {
	host     string
	publicIP string
	name     string
	domain   string
	port     int
	pw       string
	cli      *client.Client
}

func (c *couchService) buildContainerPorts() nat.PortSet {
	ports := nat.PortSet{}
	p := nat.Port(fmt.Sprintf("%d/%s", c.port, "tcp"))
	ports[p] = struct{}{}
	return ports
}

func (c *couchService) buildContainerPortBindingOptions() nat.PortMap {
	bindings := nat.PortMap{}
	p := nat.Port(fmt.Sprintf("%d/%s", c.port, "tcp")) // chaincode port
	hostp := fmt.Sprintf("%d", c.port)
	binding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: hostp,
	}
	bindings[p] = append(bindings[p], binding)
	return bindings
}

func NewCouchService(host, publicIP, name, domain, pw string, port int) ContainerService {
	return &couchService{
		host:     host,
		publicIP: publicIP,
		name:     name,
		domain:   domain,
		pw:       pw,
		port:     port,
	}
}

func (c *couchService) Run(ctx context.Context, config *container.Config, config2 *container.HostConfig, s string) error {
	createRes, err := c.cli.ContainerCreate(ctx, config, config2, nil, nil, s)
	if err != nil {
		return err
	}

	//if err := c.SetNetwork(ctx, createRes.ID); err != nil {
	//	return err
	//}
	log.Println(config, config2)
	if err := c.cli.ContainerStart(ctx, createRes.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func (c *couchService) Pull(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (c *couchService) Login(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (c *couchService) List(ctx context.Context) ([]types.Container, error) {
	//TODO implement me
	panic("implement me")
}

func (c *couchService) Close() error {
	return c.cli.Close()
}

func (c *couchService) Conn() error {
	address := fmt.Sprintf("tcp://%s:%d", c.host, 2376)
	cacertPath := fmt.Sprintf("/root/txhyjuicefs/%s/certs/ca.pem", c.publicIP)
	certPath := fmt.Sprintf("/root/txhyjuicefs/%s/certs/client.pem", c.publicIP)
	keyPath := fmt.Sprintf("/root/txhyjuicefs/%s/certs/client-key.pem", c.publicIP)
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(),
		client.WithHost(address), client.WithTLSClientConfig(cacertPath, certPath, keyPath))
	if err != nil {
		return err
	}
	c.cli = cli
	return nil
}

func (c *couchService) GenConfig(ctx context.Context) (*container.Config, *container.HostConfig) {
	portset := c.buildContainerPorts()
	portbidng := c.buildContainerPortBindingOptions()

	containerConfig := &container.Config{
		Image: "harbor.sxtxhy.com/gcbaas-gm/couchdb:3.1.1",
		Env: []string{
			"COUCHDB_USER=admin",
			"COUCHDB_PASSWORD=adminpw",
		},
		Hostname:     c.domain,
		Domainname:   c.domain,
		ExposedPorts: portset,
		Labels: map[string]string{
			"service": "hyperledger-fabric",
		},
		NetworkDisabled: false,
	}

	volumeBind := fmt.Sprintf("/root/txhyjuicefs/%s:/opt/couchdb/data", c.name)
	hostConfig := &container.HostConfig{
		PortBindings: portbidng,
		Binds: []string{
			volumeBind,
		},
		Resources: container.Resources{
			CPUShares: 1024,
			Memory:    536870912,
		},
		NetworkMode: "fabric_network",
	}

	return containerConfig, hostConfig
}

func (c *couchService) SetNetwork(ctx context.Context, containerId string) error {
	// 获取网络
	networkName := "fabric_network"
	fabric_network, err := c.cli.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{})
	if err != nil {
		return err
	}

	endpointConfig := &network.EndpointSettings{
		NetworkID: fabric_network.ID,
	}
	err = c.cli.NetworkConnect(ctx, fabric_network.ID, containerId, endpointConfig)
	if err != nil {
		return err
	}
	return nil
}
