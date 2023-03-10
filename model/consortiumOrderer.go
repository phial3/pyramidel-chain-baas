package model

type ConsortiumOrderer struct {
	Base
	OrganizationId uint `json:"organizationId" gorm:"column:organizationId"` // 组织
	OrdererId      uint `json:"ordererId" gorm:"column:ordererId"`           // orderer id
}
