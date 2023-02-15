package model

import "gorm.io/gorm"

type Organization struct {
	Uscc   string `json:"uscc" gorm:"column:uscc"`     // 统一信用代码
	Domain int    `json:"Domain" gorm:"column:Domain"` // 域名 ${uscc}.example.com
	// 到期时间
	// 重启时间
	gorm.Model
}
