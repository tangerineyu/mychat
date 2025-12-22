package repo

import (
	"my-chat/internal/model"

	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(message *model.Message) error
}
type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}
func (r *messageRepository) CreateMessage(message *model.Message) error {
	return r.db.Create(message).Error
}
