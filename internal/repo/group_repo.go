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
}
type groupRepository struct {
	db  *gorm.DB
	rdb *redis.Client
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
