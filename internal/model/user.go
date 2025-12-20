package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uuid      string `gorm:"type:varchar(64);uniqueIndex;not null;comment:用户标识"`
	Telephone string `gorm:"type:varchar(20);index;not null;comment:手机号"`
	Password  string `gorm:"type:varchar(128);not null;comment:密码"`
	Nickname  string `gorm:"type:varchar(64);comment:昵称"`
	Avatar    string `gorm:"type:varchar(255);comment:头像"`
	Status    int    `gorm:"default:1;comment:状态1：正常 2：禁用"`
}

func (User) TableName() string {
	return "users"
}
