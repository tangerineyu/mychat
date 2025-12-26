package repo

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"my-chat/internal/model"
)

type SessionRepository interface {
	GetList(userId string) ([]*model.Session, error)
	UpsertSession(session *model.Session) error
}
type sessionRepository struct {
	db *gorm.DB
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
		DoUpdates: clause.AssignmentColumns([]string{"last_msg", "last_time", "updated_at"}),
	}).Create(session).Error
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}
