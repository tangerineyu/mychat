package model

import (
	"gorm.io/gorm"
)

const (
	RoleMember = 0
	RoleAdmin  = 1
	RoleOwner  = 2
)

type GroupMember struct {
	gorm.Model
	GroupId  string `gorm:"type:varchar(64);not null;index:idx_group_member;comment:群组UUID"`
	UserId   string `gorm:"type:varchar(64);not null;index:idx_group_member;comment:用户UUID"`
	Nickname string `gorm:"type:varchar(64);comment:群内昵称"`
	Role     int    `gorm:"type:tinyint;default:0;comment:角色 0:成员 1:管理"`
}

func (GroupMember) TableName() string {
	return "group_members"
}
