package websocket

import (
	"context"
	"encoding/json"
	"my-chat/internal/model"
	"my-chat/internal/mq"
	"my-chat/pkg/zlog"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// 批量处理参数
const (
	BatchSize    = 100
	BatchTimeout = 1 * time.Second
)

// StartConsumer 启动消费者
func (manager *ClientManager) StartConsumer() {
	zlog.Info("Kafka Consumer Started...")
	msgChan := make(chan []byte, 256)
	go func() {
		for {
			m, err := mq.GlobalKafka.Reader.ReadMessage(context.Background())
			if err != nil {
				zlog.Error("Kafka.Reader Error", zap.Error(err))
				time.Sleep(100 * time.Millisecond)
				continue
			}
			msgChan <- m.Value
		}
	}()
	var batchBuffer []*model.Message
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	flushDB := func() {
		if len(batchBuffer) == 0 {
			return
		}
		err := manager.chatService.BatchSave(batchBuffer)
		if err != nil {
			zlog.Error("Batch save to DB error", zap.Error(err))
		} else {
			zlog.Info("Batch save to DB success",
				zap.Int("count", len(batchBuffer)))
		}
		batchBuffer = make([]*model.Message, 0, BatchSize)
	}
	for {
		select {
		case rawMsg := <-msgChan:
			/**m, err := mq.GlobalKafka.Reader.ReadMessage(context.Background())
			if err != nil {
				zlog.Error("Kafka Read Error", zap.Error(err))
			}**/
			var kafkaMsg Message
			if err := json.Unmarshal(rawMsg, &kafkaMsg); err != nil {
				return
			}
			if kafkaMsg.Action != ActionChatMessage {
				return
			}
			var chatData ChatMessageContent
			if err := json.Unmarshal(kafkaMsg.Content, &chatData); err != nil {
				return
			}
			jsonBytes, _ := json.Marshal(chatData)
			if chatData.Type == 1 {
				manager.sendToUser(chatData.ReceiverId, jsonBytes)
				manager.sendToUser(chatData.SendId, jsonBytes)
			} else if chatData.Type == 2 {
				memberIds, _ := manager.chatService.GetGroupMemberIDs(chatData.ReceiverId)
				for _, memberId := range memberIds {
					manager.sendToUser(memberId, jsonBytes)
				}

			}
			//now := time.Now()
			msgModel := &model.Message{
				Uuid:       "M" + uuid.New().String(),
				FromUserId: chatData.SendId,
				ToId:       chatData.ReceiverId,
				Content:    chatData.Content,
				Type:       chatData.Type,
				MediaType:  1, // 暂时写死文本

			}
			batchBuffer = append(batchBuffer, msgModel)
			if len(batchBuffer) >= BatchSize {
				flushDB()
				ticker.Reset(BatchTimeout) // 重置定时器
			}
		case <-ticker.C:
			flushDB()
		}
	}
}
