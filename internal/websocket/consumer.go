package websocket

import (
	"context"
	"encoding/json"
	"my-chat/internal/model"
	"my-chat/pkg/util/snowflake"
	"my-chat/pkg/zlog"
	"time"

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
	//msgChan := make(chan []byte, 256)
	go func() {
		for {
			m, err := manager.mqClient.Reader.ReadMessage(context.Background())
			if err != nil {
				zlog.Error("Kafka.Reader Error", zap.Error(err))
				time.Sleep(100 * time.Millisecond)
				continue
			}
			var kafkaMsg Message
			if err := json.Unmarshal(m.Value, &kafkaMsg); err != nil {
				zlog.Error("Kafka.Msg Unmarshal Error", zap.Error(err))
				continue
			}
			if kafkaMsg.Action != ActionChatMessage {
				continue
			}
			var chatData ChatMessageContent
			if err := json.Unmarshal(kafkaMsg.Content, &chatData); err != nil {
				zlog.Error("Unmarshal Chat data failed", zap.Error(err))
				continue
			}

			// 4. 【核心步骤】生成 ID 并 落库 (DB Persistence)
			// 确保消息有 ID，如果没有则生成一个 (防止旧版本消息导致空 ID)
			if chatData.Uuid == "" {
				chatData.Uuid = snowflake.GenStringID()
			}

			msgModel := &model.Message{
				Uuid:       chatData.Uuid,
				FromUserId: chatData.SendId,
				ToId:       chatData.ReceiverId,
				Content:    chatData.Content,
				Type:       chatData.Type,
				MediaType:  1, // 暂时写死文本，后续可在 ChatMessageContent 中透传
			}

			// 同步写入 MySQL
			err = manager.chatService.InsertMessage(msgModel)
			if err != nil {
				// 【重要】落库失败处理
				// 这是一个严重问题，意味着消息丢了。
				// 生产环境通常会：1. 重试 N 次  2. 放入死信队列 (DLQ)
				// 这里为了简化，我们记录 Error 日志，并且跳过后续推送（保证数据一致性：没落库就不让用户看见）
				zlog.Error("Insert Message to DB failed!!!",
					zap.String("uuid", msgModel.Uuid),
					zap.Error(err))
				continue
			}

			// ---------------------------------------------------------
			// 程序执行到这里，说明消息已经安全躺在 MySQL 里了。
			// 下面的操作（Session更新、推送）即使失败，用户也可以通过拉取历史记录看到消息。
			// ---------------------------------------------------------

			// 5. 更新 Session (会话列表)
			// 直接复用你原来的逻辑，但放在了落库之后
			currentTs := time.Now().Unix()
			if chatData.Type == 1 {
				// 私聊：更新发送者会话
				_ = manager.sessionRepo.UpsertSession(&model.Session{
					UserId:    chatData.SendId,
					TargetId:  chatData.ReceiverId,
					Type:      1,
					LastMsg:   chatData.Content,
					LastTime:  currentTs,
					UnreadCnt: 0,
				})
				_ = manager.sessionRepo.DeleteSessionCache(chatData.SendId)

				// 私聊：更新接收者会话
				_ = manager.sessionRepo.UpsertSession(&model.Session{
					UserId:    chatData.ReceiverId,
					TargetId:  chatData.SendId,
					Type:      1,
					LastMsg:   chatData.Content,
					LastTime:  currentTs,
					UnreadCnt: 1, // 接收者未读 +1 (这里简单逻辑先写死 1，更完善的是查出来+1 或者 redis incr)
				})
				_ = manager.sessionRepo.DeleteSessionCache(chatData.ReceiverId)

			} else if chatData.Type == 2 {
				// 群聊：更新群信息的 LastMsg
				err := manager.groupRepo.UpdateGroupLastMsg(chatData.ReceiverId, "群消息:"+chatData.Content, currentTs)
				if err != nil {
					zlog.Error("update group last msg failed", zap.Error(err))
				}
				// 注意：群聊没有给每个成员更新 Session 表，因为那会造成写扩散。
				// 通常群聊列表的 LastMsg 都是直接查群信息的，或者用户上线时拉取。
			}

			// 6. WebSocket 广播 (Push)
			zlog.Info("Consumer处理消息成功，准备推送",
				zap.String("uuid", chatData.Uuid),
				zap.String("sender", chatData.SendId),
				zap.String("receiver", chatData.ReceiverId))

			// 准备推送的数据
			jsonBytes, _ := json.Marshal(chatData)

			if chatData.Type == 1 {
				// 私聊：推给接收方和发送方（多端同步）
				manager.sendToUser(chatData.ReceiverId, jsonBytes)
				manager.sendToUser(chatData.SendId, jsonBytes)
			} else if chatData.Type == 2 {
				// 群聊：推给所有群成员
				memberIds, _ := manager.chatService.GetGroupMemberIDs(chatData.ReceiverId)
				for _, memberId := range memberIds {
					manager.sendToUser(memberId, jsonBytes)
				}
			}
		}
	}()

}
