package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"my-chat/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SessionRepository interface {
	GetList(userId string) ([]*model.Session, error)
	UpsertSession(session *model.Session) error

	GetListFromCache(userId string) ([]*model.SessionCache, error)
	SaveSessionToCache(userId string, session *model.SessionCache) error
	DeleteSessionCache(userId string) error
}
type sessionRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func (s *sessionRepository) GetListFromCache(userId string) ([]*model.SessionCache, error) {
	ctx := context.Background()
	//存顺序
	keySeq := fmt.Sprintf("im:session:seq:%s", userId)
	//存数据
	keyData := fmt.Sprintf("im:session:data:%s", userId)
	//拿到最新的100个会话id
	targetIds, err := s.rdb.ZRevRange(ctx, keySeq, 0, 100).Result()
	if err != nil {
		return nil, err
	}
	if len(targetIds) == 0 {
		return nil, nil
	}
	jsonList, err := s.rdb.HMGet(ctx, keyData, targetIds...).Result()
	if err != nil {
		return nil, err
	}
	var result []*model.SessionCache
	for _, v := range jsonList {
		if v == nil {
			continue
		}
		var item model.SessionCache
		if str, ok := v.(string); ok {
			_ = json.Unmarshal([]byte(str), &item)
			result = append(result, &item)
		}
	}
	return result, nil
}

// 写入缓存
func (s *sessionRepository) SaveSessionToCache(userId string, session *model.SessionCache) error {
	ctx := context.Background()
	//ZSET存顺序
	keySeq := fmt.Sprintf("im:session:seq:%s", userId)
	//HASH存数据
	keyData := fmt.Sprintf("im:session:data:%s", userId)
	dataBytes, _ := json.Marshal(session)
	//管道，打包发送批量执行
	pipe := s.rdb.Pipeline()
	//ZSET，用时间戳作为分数
	pipe.ZAdd(ctx, keySeq, redis.Z{Score: float64(session.LastTime), Member: session.TargetId})
	//HASH， 存具体内容
	pipe.HSet(ctx, keyData, session.TargetId, string(dataBytes))
	//设置过期时间
	pipe.Expire(ctx, keySeq, 168*time.Hour)
	pipe.Expire(ctx, keyData, 168*time.Hour)
	//pipe批量执行
	_, err := pipe.Exec(ctx)
	return err
}

// 删除缓存，有新消息时，强制让缓存失效
func (s *sessionRepository) DeleteSessionCache(userId string) error {
	ctx := context.Background()
	keySeq := fmt.Sprintf("im:session:seq:%s", userId)
	keyData := fmt.Sprintf("im:session:data:%s", userId)
	return s.rdb.Del(ctx, keySeq, keyData).Err()
}

func (s *sessionRepository) GetList(userId string) ([]*model.Session, error) {
	var list []*model.Session
	err := s.db.Where("user_id = ?", userId).
		Order("last_time DESC").
		Find(&list).Error
	return list, err
}

func (s *sessionRepository) UpsertSession(session *model.Session) error {
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "target_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_msg", "last_time", "unread_cnt", "updated_at", "type"}),
	}).Create(session).Error
}

func NewSessionRepository(db *gorm.DB, rdb *redis.Client) SessionRepository {
	return &sessionRepository{
		db:  db,
		rdb: rdb,
	}
}
