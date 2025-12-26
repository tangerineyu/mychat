package websocket

import (
	"context"
	"encoding/json"
	"my-chat/internal/model"
	"my-chat/internal/mq"
	"my-chat/pkg/zlog"
	"time"

	"go.uber.org/zap"
)

// StartConsumer 启动消费者
func (manager *ClientManager) StartConsumer() {
	zlog.Info("Kafka Consumer Started...")

	for {
		//阻塞读取 Kafka 消息
		m, err := mq.GlobalKafka.Reader.ReadMessage(context.Background())
		if err != nil {
			zlog.Error("Kafka Read Error", zap.Error(err))
			continue
		}
		//解析消息
		var rawMsg Message
		if err := json.Unmarshal(m.Value, &rawMsg); err != nil {
			zlog.Error("Consumer Parse Error", zap.Error(err))
			continue
		}
		if rawMsg.Action == ActionChatMessage {
			var chatData ChatMessageContent
			if err := json.Unmarshal(rawMsg.Content, &chatData); err != nil {
				continue
			}

			jsonBytes, err := manager.chatService.SaveAndFactory(
				chatData.SendId,
				chatData.ReceiverId,
				chatData.Content,
				chatData.Type,
				1,
			)
			if err != nil {
				zlog.Error("Save Message Error", zap.Error(err))
				continue
			}
			currentTs := time.Now().Unix()
			if chatData.Type == 1 {
				manager.sendToUser(chatData.ReceiverId, jsonBytes)
				manager.sendToUser(chatData.SendId, jsonBytes)

				manager.sessionRepo.UpsertSession(&model.Session{
					UserId:   chatData.SendId,
					TargetId: chatData.ReceiverId,
					Type:     1,
					LastMsg:  chatData.Content,
					LastTime: currentTs,
				})
				manager.sessionRepo.UpsertSession(&model.Session{
					UserId:    chatData.ReceiverId,
					TargetId:  chatData.SendId,
					Type:      1,
					LastMsg:   chatData.Content,
					LastTime:  currentTs,
					UnreadCnt: 1,
				})

			} else if chatData.Type == 2 {
				memberIds, _ := manager.chatService.GetGroupMemberIDs(chatData.ReceiverId)
				for _, memberId := range memberIds {
					manager.sendToUser(memberId, jsonBytes)
				}
				manager.sessionRepo.UpsertSession(&model.Session{
					UserId:   chatData.SendId,
					TargetId: chatData.ReceiverId,
					Type:     2,
					LastMsg:  "[群消息]" + chatData.Content,
					LastTime: currentTs,
				})
			}
		}
	}
}
