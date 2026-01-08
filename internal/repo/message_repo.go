package repo

import (
	"my-chat/internal/model"

	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(message *model.Message) error
	GetMessages(userId, targetId string, chatType int, offset, limit int) ([]*model.Message, error)
	BatchCreate(messages []*model.Message) error
}
type messageRepository struct {
	db *gorm.DB
}

func (r *messageRepository) BatchCreate(messages []*model.Message) error {
	if len(messages) == 0 {
		return nil
	}
	return r.db.Create(messages).Error
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}
func (r *messageRepository) CreateMessage(message *model.Message) error {
	return r.db.Create(message).Error
}
func (r *messageRepository) GetMessages(userId, targetId string, chatType int, offset, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	db := r.db.Model(&model.Message{})
	if chatType == 1 {
		db = db.Where("type = 1 AND ((from_user_id = ? AND to_id = ?) OR (from_user_id = ? AND to_id = ?))",
			userId, targetId, targetId, userId)
	} else {
		db = db.Where("type = 2 AND to_id = ?", targetId)
	}
	err := db.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error
	return messages, err
}
