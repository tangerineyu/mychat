package model

import "gorm.io/gorm"

const (
	ContactTypeFriend = 1 //好友
	ContactTypeGroup  = 2 //群组
	ContactTypeBlack  = 3 //拉黑
)

type Contact struct {
	gorm.Model
	OwnerId  string `gorm:"type:varchar(64);index;not null;comment:谁的通讯录"`
	TargetId string `gorm:"type:varchar(64);index;not null;comment:好友Id"`
	Type     int    `gorm:"type:tinyint;default:1;comment:好友类型 1:单聊好友 2:群聊群组"`
	Desc     string `gorm:"type:varchar(255);comment:备注"`
}

func (Contact) TableName() string {
	return "contacts"
}

type ContactApply struct {
	gorm.Model
	UserId   string `gorm:"type:varchar(64);index;not null;comment:申请人Id"`
	TargetId string `gorm:"type:varchar(64);index;not null;comment:被申请人Id"`
	Msg      string `gorm:"type:varchar(255);comment:申请理由"`
	Status   int    `gorm:"type:tinyint;default:0;comment:申请状态 0:待处理 1:已同意 2:已拒绝"`
}

func (ContactApply) TableName() string {
	return "contact_applies"
}
