package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	volume2 "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
)

var _ ContainerService = (*peerService)(nil)

type peerService struct {
	Host             string
	PublicIP         string
	Port             string
	orgUscc          string
	serverName       string
	serverDomain     string
	dbport           int
	ccport           int
	bootstrapAddress string
	cli              *client.Client
}

func (p *peerService) Run(ctx context.Context, config *container.Config, config2 *container.HostConfig, s string) error {
	couchServe := NewCouchService(p.Host, p.PublicIP, "couch_"+p.serverName, "couch."+p.serverName+".com", p.serverName, p.dbport)
	if err := couchServe.Conn(); err != nil {
		return err
	}
	couchConfig, couchHostConfig := couchServe.GenConfig(ctx)
	if err := couchServe.Run(ctx, couchConfig, couchHostConfig, "couch."+p.serverName+".com"); err != nil {
		return err
	}
	createRes, err := p.cli.ContainerCreate(ctx, config, config2, nil, nil, s)
	if err != nil {
		return err
	}

	//if err := p.SetNetwork(ctx, createRes.ID); err != nil {
	//	return err
	//}
	log.Println(config, config2)
	if err := p.cli.ContainerStart(ctx, createRes.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func (p *peerService) Pull(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (p *peerService) Login(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (p *peerService) List(ctx context.Context) ([]types.Container, error) {
	//TODO implement me
	panic("implement me")
}

func (p *peerService) Close() error {
	return p.cli.Close()
}

func (p *peerService) Conn() error {
	address := fmt.Sprintf("tcp://%s:%d", p.Host, 2376)
	cacertPath := fmt.Sprintf("/root/txhyjuicefs/%s/certs/ca.pem", p.PublicIP)
	certPath := fmt.Sprintf("/root/txhyjuicefs/%s/certs/client.pem", p.PublicIP)
	keyPath := fmt.Sprintf("/root/txhyjuicefs/%s/certs/client-key.pem", p.PublicIP)
	log.Println(address, cacertPath)
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(),
		client.WithHost(address), client.WithTLSClientConfig(cacertPath, certPath, keyPath))
	if err != nil {
		return err
	}
	p.cli = cli
	return nil
}

func (p *peerService) GenConfig(ctx context.Context) (*container.Config, *container.HostConfig) {
	mspId := fmt.Sprintf("CORE_PEER_LOCALMSPID=%sMsp", p.orgUscc)
	peerId := fmt.Sprintf("CORE_PEER_ID=%s", p.serverDomain)
	peerAddress := fmt.Sprintf("CORE_PEER_ADDRESS=%s:%s", p.serverDomain, p.Port)
	listenAddress := fmt.Sprintf("CORE_PEER_LISTENADDRESS=0.0.0.0:%s", p.Port)
	chaincodeAddress := fmt.Sprintf("CORE_PEER_CHAINCODEADDRESS=%s:%d", p.serverDomain, p.ccport)
	chaincodeListenAddress := fmt.Sprintf("CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:%d", p.ccport)
	gossipBootstrapAddress := fmt.Sprintf("CORE_PEER_GOSSIP_BOOTSTRAP=%s", p.bootstrapAddress)
	gossipEndpoint := fmt.Sprintf("CORE_PEER_GOSSIP_EXTERNALENDPOINT=%s:%s", p.serverDomain, p.Port)
	couchAddress := fmt.Sprintf("CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=%s:%d", "couch."+p.serverName+".com", 5984)
	couchUser := fmt.Sprintf("CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin")
	couchPw := fmt.Sprintf("CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=adminpw")

	portset := p.buildContainerPorts()
	portbinding := p.buildContainerPortBindingOptions()

	containerConfig := &container.Config{
		Image: "harbor.sxtxhy.com/gcbaas-gm/fabric-peer:latest",
		Cmd:   []string{"peer", "node", "start"},
		Env: []string{
			"CORE_PEER_PROFILE_ENABLED=true",
			"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
			"CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=fabric_network",
			"FABRIC_LOGGING_SPEC=INFO", "CORE_PEER_TLS_ENABLED=true",
			"CORE_PEER_PROFILE_ENABLED=true",
			"CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt",
			"CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key",
			"CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt",
			"CORE_CHAINCODE_GOLANG_RUNTIME=harbor.sxtxhy.com/gcbaas-gm/fabric-baseos:latest",
			"CORE_CHAINCODE_BUILDER=harbor.sxtxhy.com/gcbaas-gm/fabric-ccenv:latest", "CORE_LEDGER_STATE_STATEDATABASE=CouchDB",
			couchAddress,
			couchUser,
			couchPw,
			mspId,
			peerId,
			peerAddress,
			listenAddress,
			chaincodeAddress,
			chaincodeListenAddress,
			gossipEndpoint,
			gossipBootstrapAddress,
		},
		Hostname:     p.serverName,
		Domainname:   p.serverDomain,
		ExposedPorts: portset,
		Labels: map[string]string{
			"service": "hyperledger-fabric",
		},
		NetworkDisabled: false,
	}

	// 创建卷
	volume, err := p.cli.VolumeCreate(context.Background(), volume2.CreateOptions{
		Name: p.serverDomain,
	})
	if err != nil {
		panic(err)
	}
	mspBind := fmt.Sprintf("/root/txhyjuicefs/organizations/%s/peers/%s/msp:/etc/hyperledger/fabric/msp", p.orgUscc, p.serverDomain)
	tlsBind := fmt.Sprintf("/root/txhyjuicefs/organizations/%s/peers/%s/tls:/etc/hyperledger/fabric/tls", p.orgUscc, p.serverDomain)
	volumeBind := fmt.Sprintf("%s:/var/hyperledger/production", volume.Name)
	hostConfig := &container.HostConfig{
		PortBindings: portbinding,
		Binds: []string{
			"/var/run/docker.sock:/host/var/run/docker.sock",
			mspBind,
			tlsBind,
			volumeBind,
		},
		NetworkMode: "fabric_network",
		Resources: container.Resources{
			CPUShares: 1024,
			Memory:    536870912,
		},
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}

	return containerConfig, hostConfig
}

func (p *peerService) SetNetwork(ctx context.Context, containerId string) error {
	// 获取网络
	networkName := "fabric_network"
	fabric_network, err := p.cli.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{})
	if err != nil {
		return err
	}

	endpointConfig := &network.EndpointSettings{
		NetworkID: fabric_network.ID,
	}
	err = p.cli.NetworkConnect(ctx, fabric_network.ID, containerId, endpointConfig)
	if err != nil {
		return err
	}
	return nil
}

func NewPeerService(host, ip, port, orgUscc, serverName, serverDomain string, dbport, ccport int) ContainerService {
	return &peerService{
		Host:         host,
		Port:         port,
		PublicIP:     ip,
		orgUscc:      orgUscc,
		serverName:   serverName,
		serverDomain: serverDomain,
		dbport:       dbport,
		ccport:       ccport,
	}
}

func (p *peerService) buildContainerPorts() nat.PortSet {
	ports := nat.PortSet{}
	ccp := nat.Port(fmt.Sprintf("%d/%s", p.ccport, "tcp"))
	ports[ccp] = struct{}{}

	sp := nat.Port(fmt.Sprintf("%s/%s", p.Port, "tcp")) //service port
	ports[sp] = struct{}{}

	return ports
}

func (p *peerService) buildContainerPortBindingOptions() nat.PortMap {
	bindings := nat.PortMap{}
	ccp := nat.Port(fmt.Sprintf("%d/%s", p.ccport, "tcp")) // chaincode port
	hostccp := fmt.Sprintf("%d", p.ccport)
	binding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: hostccp,
	}
	bindings[ccp] = append(bindings[ccp], binding)

	sp := nat.Port(fmt.Sprintf("%s/%s", p.Port, "tcp")) //service port
	hostsp := fmt.Sprintf("%s", p.Port)
	binding = nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: hostsp,
	}
	bindings[sp] = append(bindings[sp], binding)
	return bindings
}
