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
			zlog.Info("Consumer处理消息",
				zap.Int("收到Type", chatData.Type),
				zap.String("发送者", chatData.SendId),
				zap.String("接收者", chatData.ReceiverId),
				zap.String("内容", chatData.Content))
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
			//Session更新逻辑,每条消息都数据库更新
			currentTs := time.Now().Unix()
			if chatData.Type == 1 {
				//私聊，更新发送者
				_ = manager.sessionRepo.UpsertSession(&model.Session{
					UserId:    chatData.SendId,
					TargetId:  chatData.ReceiverId,
					Type:      1,
					LastMsg:   chatData.Content,
					LastTime:  currentTs,
					UnreadCnt: 0,
				})
				//删除发送者的缓存，使每一次拉列表都会重新加载数据
				err := manager.sessionRepo.DeleteSessionCache(chatData.SendId)
				if err != nil {
					return
				}
				//私聊，更新接收者
				_ = manager.sessionRepo.UpsertSession(&model.Session{
					UserId:    chatData.ReceiverId,
					TargetId:  chatData.SendId,
					Type:      1,
					LastMsg:   chatData.Content,
					LastTime:  currentTs,
					UnreadCnt: 1,
				})
				err = manager.sessionRepo.DeleteSessionCache(chatData.ReceiverId)
				if err != nil {
					return
				}
			} else if chatData.Type == 2 {
				memberIds, err := manager.chatService.GetGroupMemberIDs(chatData.ReceiverId)
				if err != nil {
					zlog.Error("Get group member id error", zap.Error(err))
					return
				}
				//直接遍历更新，这里有写扩散问题，如果在一个500人群聊里发消息，
				//要遍历499次，也就是要操作数据库500次（1次存，499次更新），如果一秒10个人发消息
				//就要操作5000次写入
				//优化方案是改为读扩散，不懂，暂时没实现
				for _, memberId := range memberIds {
					unread := 0
					if memberId != chatData.SendId {
						unread = 1
					}
					_ = manager.sessionRepo.UpsertSession(&model.Session{
						UserId:    chatData.SendId,
						TargetId:  chatData.ReceiverId, //群Id
						Type:      2,
						LastMsg:   "群消息:" + chatData.Content,
						LastTime:  currentTs,
						UnreadCnt: unread,
					})
					manager.sendToUser(memberId, jsonBytes)
					_ = manager.sessionRepo.DeleteSessionCache(memberId)
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
