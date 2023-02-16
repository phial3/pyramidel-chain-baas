package model

import "gorm.io/gorm"

type Organization struct {
	Uscc           string `json:"uscc" gorm:"column:uscc;unique"`                     // 统一信用代码
	Domain         string `json:"domain" gorm:"column:domain;unique"`                 // 域名 ${uscc}.example.com
	CaHostId       Host   `json:"caHostId" gorm:"foreignKey:ID;reference:caHostId"`   // ca服务运行服务器id
	CaUser         string `json:"caUser" gorm:"column:caUser"`                        // ca服务root用户
	CaPassword     string `json:"caPassword" gorm:"column:caPassword"`                // ca服务root密码
	CaServerPort   int    `json:"caServerPort" gorm:"column:caServerPort"`            // ca服务运行端口
	CaServerDomain string `json:"caServerDomain" gorm:"column:caServerDomain;unique"` // ca服务域名 容器名
	CaServerName   string `json:"caServerName" gorm:"column:caServerName;unique"`     // ca服务名 FABRIC_CA_SERVER_CA_NAME
	gorm.Model
}

func (Organization) TableName() string {
	return "baas_organization"
}

func (o *Organization) Create() error {
	tx := db.Session(&gorm.Session{
		SkipDefaultTransaction: true,
	})
	if err := tx.Create(o).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
