package model

import "gorm.io/gorm"

type Peer struct {
	Domain         int    `json:"Domain" gorm:"column:Domain"`           // 域名 ${uscc}.example.com
	DueTime        string `json:"dueTime" gorm:"column:dueTime"`         // 到期时间
	RestartTime    string `json:"restartTime" gorm:"column:restartTime"` // 重启时间
	NodeCore       string `json:"nodeCore" gorm:"column:nodeCore"`
	NodeMemory     string `json:"nodeMemory" gorm:"column:nodeMemory"`
	NodeBandwidth  string `json:"nodeBandwidth" gorm:"column:nodeBandwidth"`
	NodeDisk       string `json:"nodeDisk" gorm:"column:nodeDisk"`
	HostId         int    `json:"hostId" gorm:"column:hostId"`                  // 所在主机
	Port           int    `json:"port" gorm:"column:port"`                      // 占用端口号
	Name           string `json:"name" gorm:"column:name"`                      // 节点名ex: peer1
	OrganizationId int    `json:"organizationId" gorm:"column:organizationId" ` // 所属组织

	gorm.Model
}

func (Peer) TableName() string {
	return "baas_peer"
}
