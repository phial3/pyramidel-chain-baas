package model

type ConsortiumOrg struct {
	Base
	OrganizationId uint         `json:"organizationId" gorm:"column:organizationId"` // 所属组织
	Organization   Organization `json:"organization" gorm:"foreignKey:OrganizationId"`
	ConsortiumId   uint         `json:"consortiumId" gorm:"column:consortiumId"` // 所属联盟
	Consortium     Consortium   `json:"consortium" gorm:"foreignKey:ConsortiumId"`
	CommitTx       bool         `json:"commitTx" gorm:"column:commitTx"`       // 是否支持提交事务
	TxSignature    bool         `json:"txSignature" gorm:"column:TxSignature"` // 是否支持签名
}

func (ConsortiumOrg) TableName() string {
	return "baas_consortium_org"
}
