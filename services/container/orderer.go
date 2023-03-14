package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	volume2 "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var _ ContainerService = (*ordererService)(nil)

type ordererService struct {
	host         string
	publicIP     string
	Port         string
	orgUscc      string
	serverName   string
	serverDomain string
	cli          *client.Client
}

func (o *ordererService) buildContainerPorts() nat.PortSet {
	//TODO implement me
	panic("implement me")
}

func (o *ordererService) buildContainerPortBindingOptions() nat.PortMap {
	//TODO implement me
	panic("implement me")
}

func (o *ordererService) Run(ctx context.Context, config *container.Config, config2 *container.HostConfig, s string) error {
	//TODO implement me
	panic("implement me")
}

func (o *ordererService) Pull(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (o *ordererService) Login(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (o *ordererService) List(ctx context.Context) ([]types.Container, error) {
	//TODO implement me
	panic("implement me")
}

func (o *ordererService) Close() error {
	//TODO implement me
	panic("implement me")
}

func (o *ordererService) Conn() error {
	address := fmt.Sprintf("tcp://%s:%d", o.host, 2376)
	cacertPath := fmt.Sprintf("/txhyjuicefs/%s/certs/ca.pem", o.publicIP)
	certPath := fmt.Sprintf("/txhyjuicefs/%s/certs/client.pem", o.publicIP)
	keyPath := fmt.Sprintf("/txhyjuicefs/%s/certs/client-key.pem", o.publicIP)
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(),
		client.WithHost(address), client.WithTLSClientConfig(cacertPath, certPath, keyPath))
	if err != nil {
		return err
	}
	o.cli = cli
	return nil
}

func (o *ordererService) GenConfig(ctx context.Context) (*container.Config, *container.HostConfig) {
	mspId := fmt.Sprintf("ORDERER_GENERAL_LOCALMSPID=%sMsp", o.orgUscc)
	orderer_host := fmt.Sprintf("ORDERER_HOST=%s", o.serverDomain)
	listenPort := fmt.Sprintf("ORDERER_GENERAL_LISTENPORT=%s", o.Port)
	portset := o.buildContainerPorts()
	portbinding := o.buildContainerPortBindingOptions()
	containerConfig := &container.Config{
		Image: "harbor.sxtxhy.com/gcbaas-gm/fabric-orderer:latest",
		Cmd:   []string{"orderer"},
		Env: []string{
			"ORDERER_GENERAL_GENESISMETHOD=file",
			"ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block",
			"ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp",
			"FABRIC_LOGGING_SPEC=INFO",
			"ORDERER_GENERAL_LISTENADDRESS=0.0.0.0",
			"ORDERER_GENERAL_TLS_ENABLED=true",
			"ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key",
			"ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt",
			"ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]",
			"ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/server.crt",
			"ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/server.key",
			"ORDERER_GENERAL_CLUSTER_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]",
			mspId,
			orderer_host,
			listenPort,
		},
		Hostname:     o.serverDomain,
		Domainname:   o.serverDomain,
		ExposedPorts: portset,
		Labels: map[string]string{
			"service": "hyperledger-fabric",
		},
		NetworkDisabled: false,
	}

	// 创建卷
	volume, _ := o.cli.VolumeCreate(context.Background(), volume2.CreateOptions{
		Name: o.serverDomain,
	})
	mspBind := fmt.Sprintf("/txhyjuicefs/organizations/%s/orderers/%s/msp:/etc/hyperledger/fabric/msp", o.orgUscc, o.serverDomain)
	tlsBind := fmt.Sprintf("/txhyjuicefs/organizations/%s/orderers/%s/tls:/etc/hyperledger/fabric/tls", o.orgUscc, o.serverDomain)
	volumeBind := fmt.Sprintf("%s:/var/hyperledger/production/orderer", volume.Name)
	hostConfig := &container.HostConfig{
		PortBindings: portbinding,
		Binds: []string{
			"/txhyjuicefs/system-genesis-block/genesis.block:/var/hyperledger/orderer/orderer.genesis.block",
			mspBind,
			tlsBind,
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

func (o *ordererService) SetNetwork(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}
