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
	IsMember(groupId, userId string) (bool, error)
	GetGroupMembers(groupId string) ([]*model.GroupMember, error)
}
type groupRepository struct {
	db  *gorm.DB
	rdb *redis.Client
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
