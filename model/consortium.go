package model

type Consortium struct {
	Base
	Name          string `json:"name" gorm:"column:name"`                     // 展示名
	ManagerPolicy string `json:"manager_policy" gorm:"column:manager_policy"` // 联盟管理策略
	LedgerStatus  bool   `json:"ledger_query" gorm:"column:ledger_query"`     // 账本状态
	ChannelName   string `json:"channel_name" gorm:"column:channel_name"`     // 格式,fabric通道名{ch}-{uuid}
	Scenes        string `json:"scenes" gorm:"column:scenes"`                 // 场景
	Version       string `json:"version" gorm:"column:version"`               // fabric通道配置版本建议累加
	Description   string `json:"description" gorm:"column:description"`       // 简介
	OrdererType   string `json:"orderer-type" gorm:"column:orderer-type"`     // 共识算法kafkaraft,etcdraft
	// TODO:添加出块规则相关字段
}

func (Consortium) TableName() string {
	return "consortium"
}
