package repo

import (
	"my-chat/internal/model"

	"gorm.io/gorm"
)

type GroupRepository interface {
	GetMemberIDs(groupId string) ([]string, error)
}
type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}
func (r *groupRepository) GetMemberIDs(groupId string) ([]string, error) {
	var userIds []string
	//只查user_id字段， Pluck是grom专门查单列数据的
	err := r.db.Model(&model.GroupMember{}).
		Where("group_id = ?", groupId).
		Pluck("user_id", &userIds).Error
	return userIds, err
}
