package model

type ConsortiumPeer struct {
	Base
	PeerId       uint       `json:"peerId" gorm:"column:peerId"`               // 节点id
	Peer         Peer       `json:"peer" gorm:"column:PeerId"`                 // peer外键
	ConsortiumId uint       `json:"consortiumId" gorm:"column:consortiumId"`   // 所属联盟
	Consortium   Consortium `json:"consortium" gorm:"foreignKey:ConsortiumId"` // 外键
	SyncBlock    int        `json:"syncBlock" gorm:"column:syncBlock"`         // 是否准入准出
}

func (ConsortiumPeer) TableName() string {
	return "baas_consortium_peer"
}
