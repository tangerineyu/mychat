package model

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Uuid     string `gorm:"type:varchar(64);uniqueIndex;not null;comment:群唯一标识"`
	Name     string `gorm:"type:varchar(64);comment:群名称"`
	Notice   string `gorm:"type:varchar(255);comment:群公告"`
	OwnerId  string `gorm:"type:varchar(64);index;comment:群主UUID"` // 冗余字段，方便快速查询
	Avatar   string `gorm:"type:varchar(255);comment:群头像"`
	LastMsg  string `gorm:"type:text"`
	LastTime int64  `gorm:"index"`
}

func (Group) TableName() string {
	return "groups"
}
