package model

import (
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/localtime"
	"gorm.io/gorm"
)

type Peer struct {
	Domain         string              `json:"Domain" gorm:"column:Domain;unique"`    // 域名 ${uscc}.example.com
	DueTime        localtime.LocalTime `json:"dueTime" gorm:"column:dueTime"`         // 到期时间
	RestartTime    localtime.LocalTime `json:"restartTime" gorm:"column:restartTime"` // 重启时间
	NodeCore       uint                `json:"nodeCore" gorm:"column:nodeCore"`
	NodeMemory     uint                `json:"nodeMemory" gorm:"column:nodeMemory"`
	NodeBandwidth  uint                `json:"nodeBandwidth" gorm:"column:nodeBandwidth"`
	NodeDisk       uint                `json:"nodeDisk" gorm:"column:nodeDisk"`
	HostId         uint                `json:"hostId" gorm:"column:hostId"` // 所在主机
	Host           Host                `json:"host" gorm:"foreignKey:HostId" `
	Port           uint                `json:"port" gorm:"column:port"`                     // 占用端口号
	Name           string              `json:"name" gorm:"column:name"`                     // 节点名ex: peer1
	SerialNumber   uint                `json:"serialNumber" gorm:"column:serialNumber"`     // 序列号
	OrganizationId uint                `json:"organizationId" gorm:"column:organizationId"` // 所属组织
	Organization   Organization        `json:"organization" gorm:"foreignKey:OrganizationId" `
	Base
}

func (Peer) TableName() string {
	return "baas_peer"
}

func (p *Peer) Create(tx *gorm.DB) error {
	if tx == nil {
		tx = db.Session(&gorm.Session{
			SkipDefaultTransaction: true,
		})
	}
	if err := tx.Create(p).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (p *Peer) GetMaxSerial(tx *gorm.DB, id uint) error {
	if tx == nil {
		tx = db.Session(&gorm.Session{
			SkipDefaultTransaction: true,
		})
	}
	if err := tx.Where("organizationId = ?", id).First(&p).Error; err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			p.ID = 0
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (p *Peer) GetByID(id uint) error {
	if err := db.Preload("Host").Where("id = ?", id).First(&p).Error; err != nil {
		return err
	}
	return nil
}
