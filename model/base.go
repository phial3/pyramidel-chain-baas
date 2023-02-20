package model

import (
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/localtime"
	"gorm.io/gorm"
)

type Base struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt localtime.LocalTime
	UpdatedAt localtime.LocalTime
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
