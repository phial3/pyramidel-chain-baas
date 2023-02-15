package model

import "gorm.io/gorm"

type Orderer struct {
	Domain        int    `json:"Domain" gorm:"column:Domain"` // 域名 ${uscc}.example.com
	DueTime       string `json:"dueTime" gorm:"column:dueTime"`
	RestartTime   string `json:"restartTime" gorm:"column:restartTime"`
	NodeCore      string `json:"nodeCore" gorm:"column:nodeCore"`
	NodeMemory    string `json:"nodeMemory" gorm:"column:nodeMemory"`
	NodeBandwidth string `json:"nodeBandwidth" gorm:"column:nodeBandwidth"`
	NodeDisk      string `json:"nodeDisk" gorm:"column:nodeDisk"`
	// 到期时间
	// 重启时间
	gorm.Model
}
