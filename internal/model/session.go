package model

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	Type      int    `gorm:"type:tinyint;default:1;comment:会话类型 1:单聊 2:群聊"`
	Top       bool   `gorm:"default:false;comment:是否置顶"`
	Mute      int    `gorm:"type:tinyint;default:0;comment:是否免打扰 0:否 1:是"`
	UnreadCnt int    `gorm:"default:0;comment:未读消息数"`
	LastMsg   string `gorm:"type:varchar(255);comment:最新消息"`
	LastTime  int64  `gorm:"index;comment:最新消息时间戳"`
	UserId    string `gorm:"type:varchar(255);uniqueIndex:idx_user_target;not null;comment:会话所属用户Id"`
	TargetId  string `gorm:"type:varchar(255);uniqueIndex:idx_user_target;not null;comment:会话目标Id，单聊为好友Id，群聊为群Id"`
}

func (Session) TableName() string {
	return "sessions"
}
