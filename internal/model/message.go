package model

import (
	"gorm.io/gorm"
	//"time"
)

// 消息类型
const (
	MsgTypeSingle = 1 //单聊
	MsgTypeGroup  = 2 //群聊
)

// 消息内容类型
const (
	MediaTypeText  = 1 //文本
	MediaTypeImage = 2 //图片
	MediaTypeAudio = 3 //语音
)

type Message struct {
	gorm.Model
	Uuid       string `gorm:"type:varchar(64);uniqueIndex;not null;comment:消息唯一标识"`
	FromUserId string `gorm:"type:varchar(64);index;not null;comment:发送者用户UUID"`
	ToId       string `gorm:"type:varchar(64);index;not null;comment:接收者UUID，单聊为用户UUID，群聊为群UUID"`
	Type       int    `gorm:"type:tinyint;default:1;comment:消息类型 1:单聊 2:群聊"`
	MediaType  int    `gorm:"type:tinyint;default:1;comment:消息内容类型 1:文本 2:图片 3:语音"`
	Content    string `gorm:"type:text;comment:消息内容"`

	PicUrl string `gorm:"type:varchar(255);default:''"`
	Url    string `gorm:"type:varchar(255);default:''"`
}

func (Message) TableName() string {
	return "messages"
}
