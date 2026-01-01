package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"my-chat/internal/model"
	"my-chat/pkg/zlog"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GroupRepository interface {
	GetMemberIDs(groupId string) ([]string, error)
	CreateGroup(group *model.Group, ownerMember *model.GroupMember) error
	AddMember(member *model.GroupMember) error
	FindGroup(groupId string) (*model.Group, error)
	FindGroupsByIds(groupIds []string) (map[string]*model.Group, error)
	IsMember(groupId, userId string) (bool, error)
	GetGroupMembers(groupId string) ([]*model.GroupMember, error)
	GetUserJoinedGroups(userId string) ([]*model.Group, error)
	RemoveMember(groupId, userId string) error
	DeleteGroup(groupId string) error
	UpdateGroupLastMsg(groupId string, content string, time int64) error
}
type groupRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

// 实现更新群最新消息
func (r *groupRepository) UpdateGroupLastMsg(groupId string, content string, lastTime int64) error {
	return r.db.Model(&model.Group{}).
		Where("uuid = ?", groupId).
		Updates(map[string]interface{}{
			"last_msg":  content,
			"last_time": lastTime,
		}).Error
}

func (r *groupRepository) FindGroupsByIds(groupIds []string) (map[string]*model.Group, error) {
	var groups []*model.Group
	if len(groupIds) == 0 {
		return make(map[string]*model.Group), nil
	}
	err := r.db.Where("id IN (?)", groupIds).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	groupMap := make(map[string]*model.Group)
	for _, group := range groups {
		groupMap[group.Uuid] = group
	}
	return groupMap, nil
}

func (r *groupRepository) DeleteGroup(groupId string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("uuid = ?", groupId).Delete(&model.Group{}).Error; err != nil {
			return err
		}
		if err := tx.Where("group_id = ?", groupId).Delete(&model.GroupMember{}).Error; err != nil {
			return err
		}
		cacheKey := fmt.Sprintf("im:group:members:%s", groupId)
		r.rdb.Del(context.Background(), cacheKey)
		return nil
	})
}

func (r *groupRepository) RemoveMember(groupId, userId string) error {
	err := r.db.Where("group_id = ? and user_id = ?", groupId, userId).Delete(&model.GroupMember{}).Error
	if err != nil {
		return err
	}
	cacheKey := fmt.Sprintf("im:group:members:%s", groupId)
	r.rdb.Del(context.Background(), cacheKey)
	return nil
}

func (r *groupRepository) GetUserJoinedGroups(userId string) ([]*model.Group, error) {
	var groups []*model.Group
	//select g.* from groups g join group_members m on g.uuid = m.group_id where m.userid = ?
	err := r.db.Table("groups").
		Select("groups.*").
		Joins("JOIN group_members ON group_members.group_id = group.uuid").
		Where("group_members.user_id = ?", userId).
		Find(&groups).Error
	return groups, err

}

func (r *groupRepository) GetGroupMembers(groupId string) ([]*model.GroupMember, error) {
	var members []*model.GroupMember
	err := r.db.Where("group_id = ?", groupId).Find(&members).Error
	return members, err
}

func (r *groupRepository) CreateGroup(group *model.Group, ownerMember *model.GroupMember) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(group).Error; err != nil {
			return err
		}
		if err := tx.Create(ownerMember).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *groupRepository) AddMember(member *model.GroupMember) error {
	err := r.db.Create(member).Error
	if err != nil {
		return err
	}
	cacheKey := fmt.Sprintf("im:group:members:%s", member.GroupId)
	//数据变动直接删除缓存
	if err := r.rdb.SAdd(context.Background(), cacheKey, member.UserId).Err(); err != nil {
		zlog.Error("Failed to delete cache", zap.String("key", cacheKey), zap.Error(err))
	}
	return nil
}

func (r *groupRepository) FindGroup(groupId string) (*model.Group, error) {
	var group model.Group
	err := r.db.Where("uuid = ?", groupId).First(&group).Error
	return &group, err
}

// 用户是不是在群里
func (r *groupRepository) IsMember(groupId, userId string) (bool, error) {
	var count int64
	err := r.db.Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupId, userId).
		Count(&count).Error
	return count > 0, err
}

func NewGroupRepository(db *gorm.DB, rdb *redis.Client) GroupRepository {
	return &groupRepository{
		db:  db,
		rdb: rdb,
	}
}
func (r *groupRepository) GetMemberIDs(groupId string) ([]string, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("im:group:members:%s", groupId)
	//先查缓存
	val, err := r.rdb.SMembers(ctx, cacheKey).Result()
	if err == nil {
		var userIds []string
		//将[]string转为[]byte再反序列化
		var byteVals []byte
		for _, v := range val {
			byteVals = append(byteVals, []byte(v)...)
		}
		if err := json.Unmarshal(byteVals, &userIds); err == nil {
			return userIds, nil
		}
	} else if err != redis.Nil {
		zlog.Error("Redis Get Error", zap.Error(err))
	}

	//缓存未命中
	var userIds []string
	//只查user_id字段， Pluck是grom专门查单列数据的
	err = r.db.Model(&model.GroupMember{}).
		Where("group_id = ?", groupId).
		Pluck("user_id", &userIds).Error
	if err != nil {
		return nil, err
	}
	jsonBytes, _ := json.Marshal(userIds)
	err = r.rdb.Set(ctx, cacheKey, jsonBytes, 1*time.Hour).Err()
	if err != nil {
		zlog.Error("Redis Set Error", zap.Error(err))
	}
	return userIds, err
}
