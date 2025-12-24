package service

import (
	"encoding/json"
	"my-chat/internal/model"
	"my-chat/internal/repo"

	"github.com/google/uuid"
)

type ChatService struct {
	msgRepo   repo.MessageRepository
	groupRepo repo.GroupRepository
}

func NewChatService(msgRepo repo.MessageRepository, groupRepo repo.GroupRepository) *ChatService {
	return &ChatService{
		msgRepo:   msgRepo,
		groupRepo: groupRepo,
	}
}

type MsgPayload struct {
	Uuid       string `json:"uuid"`
	FromUserId string `json:"form_user_id"`
	ToId       string `json:"to_id"`
	Content    string `json:"content"`
	Type       int    `json:"type"`
	MediaType  int    `json:"media_type"`
	CreatedAt  string `json:"created_at"`
}

func (s *ChatService) SaveAndFactory(fromId, toId, content string, chatType, mediaType int) ([]byte, error) {
	msgModel := &model.Message{
		Uuid:       "M" + uuid.New().String(),
		FromUserId: fromId,
		ToId:       toId,
		Content:    content,
		Type:       chatType,
		MediaType:  mediaType,
	}
	err := s.msgRepo.CreateMessage(msgModel)
	if err != nil {
		return nil, err
	}
	payload := MsgPayload{
		Uuid:       msgModel.Uuid,
		FromUserId: msgModel.FromUserId,
		ToId:       msgModel.ToId,
		Content:    msgModel.Content,
		Type:       msgModel.Type,
		MediaType:  msgModel.MediaType,
		CreatedAt:  msgModel.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return json.Marshal(&payload)
}
func (s *ChatService) GetGroupMemberIDs(groupId string) ([]string, error) {
	//未来可以使用redis
	return s.groupRepo.GetMemberIDs(groupId)
}
func (s *ChatService) GetHistory(userId, targetId string, chatType int) ([]MsgPayload, error) {
	messages, err := s.msgRepo.GetMessages(userId, targetId, chatType, 0, 50)
	if err != nil {
		return nil, err
	}
	var result []MsgPayload
	for _, msg := range messages {
		result = append(result, MsgPayload{
			Uuid:       msg.Uuid,
			FromUserId: msg.FromUserId,
			ToId:       msg.ToId,
			Content:    msg.Content,
			Type:       msg.Type,
			MediaType:  msg.MediaType,
			CreatedAt:  msg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}
