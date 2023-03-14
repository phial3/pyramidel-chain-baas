package model

type ConsortiumOrderer struct {
	Base
	OrdererId    uint       `json:"ordererId" gorm:"column:ordererId"`       // orderer id
	Orderer      Orderer    `json:"orderer" gorm:"foreignKey:OrdererId"`     // 外键
	ConsortiumId uint       `json:"consortiumId" gorm:"column:consortiumId"` // 所属联盟
	Consortium   Consortium `json:"consortium" gorm:"foreignKey:ConsortiumId"`
}

func (ConsortiumOrderer) TableName() string {
	return "baas_consortium_orderer"
}
