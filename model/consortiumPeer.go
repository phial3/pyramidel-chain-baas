package model

type ConsortiumPeer struct {
	Base
	OrganizationId uint `json:"organizationId" gorm:"column:organizationId"` // 组织
	PeerId         uint `json:"peerId" gorm:"column:peerId"`                 // 节点id
	EAEStatus      bool `json:"status" gorm:"column:status"`                 // 是否准入准出
}
