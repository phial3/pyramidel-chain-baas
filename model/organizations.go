package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/jsonrpcClient"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/mq"
	psutilclient "github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/client"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/remotessh"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/request/organizations"
	"github.com/hxx258456/pyramidel-chain-baas/services/container"
	"github.com/hxx258456/pyramidel-chain-baas/services/loadbalance"
	"github.com/melbahja/goph"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	"time"
)

type Organization struct {
	Uscc           string `json:"uscc" gorm:"column:uscc;unique"`     // 统一信用代码
	Domain         string `json:"domain" gorm:"column:domain;unique"` // 域名 ${uscc}.example.com
	CaHostId       uint   `json:"caHostId" gorm:"column:caHostId"`    // ca服务运行服务器id
	Host           Host   `json:"host" gorm:"foreignKey:CaHostId"`
	CaUser         string `json:"caUser" gorm:"column:caUser"`                        // ca服务root用户
	CaPassword     string `json:"caPassword" gorm:"column:caPassword"`                // ca服务root密码
	CaServerPort   uint   `json:"caServerPort" gorm:"column:caServerPort"`            // ca服务运行端口
	CaServerDomain string `json:"caServerDomain" gorm:"column:caServerDomain;unique"` // ca服务域名 容器名
	CaServerName   string `json:"caServerName" gorm:"column:caServerName;unique"`     // ca服务名 FABRIC_CA_SERVER_CA_NAME
	Peers          []Peer
	Orderers       []Orderer
	Status         int `json:"status" gorm:"column:status"` // 状态
	Base
}

func (Organization) TableName() string {
	return "baas_organization"
}

func (o *Organization) Exists(uscc string) (bool, error) {
	var result = make([]Organization, 1)
	if err := db.Where("uscc = ?", uscc).Find(&result).Error; err != nil {
		return false, err
	}

	if len(result) > 0 {
		*o = result[0]
		return true, nil
	} else {
		return false, nil
	}
}
func (o *Organization) Create(param organizations.Organizations, balancer loadbalance.LBS) error {
	tx := db.Session(&gorm.Session{
		SkipDefaultTransaction: true,
	})
	exists, err := o.Exists(param.OrgUscc)
	if err != nil {
		return err
	}
	if err := o.Host.QueryById(o.CaHostId, &o.Host); err != nil {
		tx.Rollback()
		return err
	}
	// ssh申请节点证书
	sshcli, err := remotessh.ConnectToHost(o.Host.Username, o.Host.Pw, o.Host.UseIp, 22)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer func(sshcli *goph.Client) {
		err := sshcli.Close()
		if err != nil {
			log.Println(err)
		}
	}(sshcli)
	if !exists {
		if err := tx.Create(o).Error; err != nil {
			tx.Rollback()
			return err
		}

		caService := container.NewCaContainerService(o.Host.UseIp, o.Host.PublicIp, strconv.Itoa(int(o.CaServerPort)), o.CaUser, o.CaPassword, o.Uscc, o.CaServerName, o.CaServerDomain)
		if err := caService.Conn(); err != nil {
			tx.Rollback()
			return err
		}
		ctx := context.Background()
		config, D := caService.GenConfig(ctx)
		if err := caService.Run(ctx, config, D, o.CaServerName); err != nil {
			tx.Rollback()
			return err
		}
		time.Sleep(time.Second * 5)
		if err := remotessh.EnrollBootstrapCa(sshcli, o.Uscc, strconv.Itoa(int(o.CaServerPort))); err != nil {
			log.Println("error bootstrapca脚本执行报错", err)
			// TODO:添加错误处理机制
		}
	}

	cureentPeer := &Peer{}
	if err := cureentPeer.GetMaxSerial(tx, o.ID); err != nil {
		tx.Rollback()
		return err
	}
	cureentOrderer := &Orderer{}
	if err := cureentOrderer.GetMaxSerial(tx, o.ID); err != nil {
		tx.Rollback()
		return err
	}
	peerList, ordererList, coNum, err := GroupList(param.NodeList)
	if err != nil {
		tx.Rollback()
		return err
	}
	peerCh := make(chan Peer)
	ordererCh := make(chan Orderer)
	wg := &sync.WaitGroup{}
	wg.Add(coNum)
	var result []interface{}
	go func(wg *sync.WaitGroup) {
		for {
			select {
			case peer, ok := <-peerCh:
				if ok {
					result = append(result, peer)
					log.Println("消息队里改变:", result)
				}
			case orderer, ok := <-ordererCh:
				if ok {
					result = append(result, orderer)
					log.Println("消息队里改变:", result)
				}
			default:
				if len(result) == coNum {
					// TODO:amqp send result
					mqcli := mq.NewRabbitMQ()
					if err := mqcli.Connect(); err != nil {
						//wg.Wait()
						//return
						panic(err)
					}
					var message []byte
					body, err := json.Marshal(&result)
					if err != nil {
						message = []byte(err.Error())
					} else {
						message = body
					}
					log.Println("message: ", string(message))
					if err := mqcli.Publish(message); err != nil {
						log.Println(err)
						wg.Wait()
						return
					}
					wg.Wait()
					return
				}
			}
		}
	}(wg)
	for i, v := range peerList {
		for j := 0; uint(j) < v.NodeNumber; j++ {
			var (
				hostid uint
				port   int
			)
			lbid := balancer.NextService()
			if lbid == 0 {
				hostid = o.CaHostId
			} else {
				hostid = lbid
			}
			host := &Host{}
			err = host.QueryById(hostid, host)

			if err != nil {
				return err
			}
			cli, err := jsonrpcClient.ConnetJsonrpc(o.Host.UseIp + ":8082")
			if err != nil {
				return err
			}
			defer cli.Close()
			port, err = psutilclient.CallGetPort(cli)
			if err != nil {
				return err
			}
			ccport, err := psutilclient.CallGetPort(cli)
			if err != nil {
				return err
			}
			dbport, err := psutilclient.CallGetPort(cli)
			if err != nil {
				return err
			}
			domain := fmt.Sprintf("peer%d.%s.pcb.com", cureentPeer.SerialNumber+uint(i)+uint(j)+1, o.Uscc)
			name := fmt.Sprintf("%s_peer%d", o.Uscc, cureentPeer.SerialNumber+uint(i)+uint(j)+1)

			//
			peer := Peer{
				Domain:         domain,
				DueTime:        param.DueTime,
				RestartTime:    param.RestartTime,
				NodeBandwidth:  v.NodeBandwidth,
				NodeCore:       v.NodeCore,
				NodeDisk:       v.NodeDisk,
				NodeMemory:     v.NodeMemory,
				Name:           name,
				SerialNumber:   cureentPeer.SerialNumber + uint(i) + uint(j) + 1,
				HostId:         hostid,
				OrganizationId: o.ID,
				Port:           uint(port),
				OrgPackageId:   param.OrgPackageId,
				Status:         0,
				CCPort:         ccport,
				DBPort:         dbport,
			}
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				// ssh申请节点证书
				sshcli, err := remotessh.ConnectToHost(o.Host.Username, o.Host.Pw, o.Host.UseIp, 22)
				if err != nil {
					peer.Status = 0
					peer.Error = err.Error()
					peerCh <- peer
					return
				}
				defer func(sshcli *goph.Client) {
					err := sshcli.Close()
					if err != nil {
						log.Println(err)
					}
				}(sshcli)
				// ssh申请节点证书
				log.Println("开始创建peer节点", domain)
				if err := remotessh.RegisterPeer(sshcli, o.Uscc, name, domain, strconv.Itoa(int(o.CaServerPort))); err != nil {
					peer.Status = 0
					peer.Error = err.Error()
					peerCh <- peer
					return
				}
				log.Println(port, dbport, ccport)
				// 启动peer节点
				peerServe := container.NewPeerService(host.PublicIp, host.PublicIp, strconv.Itoa(port), o.Uscc, name, domain, dbport, ccport)
				if err := peerServe.Conn(); err != nil {
					peer.Status = 0
					peer.Error = err.Error()
					peerCh <- peer
					return
				}
				defer peerServe.Close()
				ctx := context.Background()
				containerConf, hostConf := peerServe.GenConfig(ctx)
				if err := peerServe.Run(ctx, containerConf, hostConf, domain); err != nil {
					peer.Status = 0
					peer.Error = err.Error()
					peerCh <- peer
					return
				}
				peer.Status = 1
				peerCh <- peer
				return
			}(wg)
			if err := peer.Create(tx); err != nil {
				tx.Rollback()
				return err
			}
		}

	}
	for i, v := range ordererList {
		for j := 0; uint(j) < v.NodeNumber; j++ {
			var (
				hostid uint
				port   int
			)
			lbid := balancer.NextService()
			if lbid == 0 {
				hostid = o.CaHostId
			} else {
				hostid = lbid
			}
			host := &Host{}
			err = host.QueryById(hostid, host)

			if err != nil {
				return err
			}
			cli, err := jsonrpcClient.ConnetJsonrpc(o.Host.UseIp + ":8082")
			if err != nil {
				return err
			}
			defer cli.Close()
			port, err = psutilclient.CallGetPort(cli)
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}
			domain := fmt.Sprintf("orderer%d.%s.pcb.com", cureentOrderer.SerialNumber+uint(i)+uint(j)+1, o.Uscc)
			name := fmt.Sprintf("%s_orderer%d", o.Uscc, cureentOrderer.SerialNumber+uint(i)+uint(j)+1)

			//ordererServe := container.NewPeerService()
			orderer := Orderer{
				Domain:         domain,
				DueTime:        param.DueTime,
				RestartTime:    param.RestartTime,
				NodeBandwidth:  v.NodeBandwidth,
				NodeCore:       v.NodeCore,
				NodeDisk:       v.NodeDisk,
				NodeMemory:     v.NodeMemory,
				Name:           name,
				SerialNumber:   cureentOrderer.SerialNumber + uint(i) + uint(j) + 1,
				HostId:         hostid,
				OrganizationId: o.ID,
				Port:           uint(port),
				OrgPackageId:   param.OrgPackageId,
				Status:         0,
			}
			go func(wg *sync.WaitGroup) {
				defer wg.Done()

				// ssh申请节点证书
				sshcli, err := remotessh.ConnectToHost(o.Host.Username, o.Host.Pw, o.Host.UseIp, 22)
				if err != nil {
					orderer.Status = 0
					orderer.Error = err.Error()
					ordererCh <- orderer
					return
				}
				defer func(sshcli *goph.Client) {
					err := sshcli.Close()
					if err != nil {
						log.Println(err)
					}
				}(sshcli)
				// ssh申请orderer证书
				if err := remotessh.RegisterOrderer(sshcli, o.Uscc, name, domain, strconv.Itoa(int(o.CaServerPort))); err != nil {
					orderer.Status = 0
					orderer.Error = err.Error()
					ordererCh <- orderer
					return
				}
				orderer.Status = 1
				ordererCh <- orderer
				return
			}(wg)
			if err := orderer.Create(tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return nil
}

// GroupList 将nodelist分组,分为peer组,orderer组
func GroupList(list []organizations.NodeList) (peerList []organizations.NodeList, ordererList []organizations.NodeList, goroNun int, err error) {
	for _, v := range list {
		switch v.NodeType {
		case 2:
			peerList = append(peerList, v)
			goroNun += int(v.NodeNumber)
		case 1:
			ordererList = append(ordererList, v)
			goroNun += int(v.NodeNumber)
		default:
			return peerList, ordererList, goroNun, errors.New("invalid node type")
		}
	}
	return peerList, ordererList, goroNun, nil
}
