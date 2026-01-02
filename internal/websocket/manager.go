package websocket

import (
	"context"
	"encoding/json"
	"my-chat/internal/mq"
	"my-chat/internal/repo"
	"my-chat/internal/service"
	"my-chat/pkg/zlog"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ClientManager struct {
	Clients    map[string]*Client
	Register   chan *Client //链接请求
	Unregister chan *Client //断开连接请求
	Broadcast  chan []byte  //消息广播

	rwLock sync.RWMutex
	//注入ChatService, 用于存消息
	chatService *service.ChatService
	sessionRepo repo.SessionRepository
	groupRepo   repo.GroupRepository
}

// 超时常量，为了方便测试超时的时间设置的比较长
const (
	HeartbeatInterval = 30 * time.Second
	HeartbeatTimeout  = 300
)

func NewClientManager(chatService *service.ChatService, sessionRepo repo.SessionRepository) *ClientManager {
	return &ClientManager{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan []byte),
		Clients:     make(map[string]*Client),
		chatService: chatService,
		sessionRepo: sessionRepo,
	}
}
func (manager *ClientManager) StartHeartbeat() {
	//定义定时器
	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()
	zlog.Info("Heartbeat checker started...")
	for range ticker.C {
		//加锁，要遍历Clients map
		manager.rwLock.Lock()
		now := time.Now().Unix()
		for userId, client := range manager.Clients {
			if now-client.HeartbeatTime > HeartbeatTimeout {
				zlog.Warn("心跳超时，下线",
					zap.String("userId", userId),
					zap.Int64("last_beat", client.HeartbeatTime))
				client.Conn.Close()
				delete(manager.Clients, userId)
			}
		}
		manager.rwLock.Unlock()
	}
}
func (manager *ClientManager) Start() {
	zlog.Info("Websocket Client Manager Started")
	//启动心跳
	go manager.StartHeartbeat()
	//启动消费者
	go manager.StartConsumer()
	for {
		select {
		case client := <-manager.Register:
			manager.rwLock.Lock()
			//检查有没有旧连接
			var oldClientToClose *Client
			if oldClient, ok := manager.Clients[client.UserId]; ok {
				oldClientToClose = oldClient //记下来，等会儿关闭，不占用锁
			}
			manager.Clients[client.UserId] = client
			manager.rwLock.Unlock()
			zlog.Info("New connection",
				zap.String("uuid", client.UserId))
			if oldClientToClose != nil {
				close(oldClientToClose.Send)
				zlog.Info("Close old connection", zap.String("uuid", client.UserId))
			}

		case client := <-manager.Unregister:
			manager.rwLock.Lock()
			if _, ok := manager.Clients[client.UserId]; ok {
				delete(manager.Clients, client.UserId)
				close(client.Send)
			}
			manager.rwLock.Unlock()
			zlog.Info("Disconnect", zap.String("uuid", client.UserId))

		case message := <-manager.Broadcast:
			manager.dispatch(message)
		}
	}
}

// 处理消息分发
func (manager *ClientManager) dispatch(message []byte) {
	/**var rawMsg Message
	if err := json.Unmarshal(message, &rawMsg); err != nil {
		zlog.Error("Failed to unmarshal message", zap.String("message", string(message)), zap.Error(err))
		return
	}
	if rawMsg.Action == ActionChatMessage {
		ctx := context.Background()
		err := mq.GlobalKafka.Publish(ctx, nil, message)
		if err != nil {
			zlog.Error("Kafka Publish Error", zap.Error(err))
			return
		}
	}**/
	var baseMsg struct {
		Action  Action          `json:"action"`
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(message, &baseMsg); err != nil {
		zlog.Error("Parse message failed", zap.Error(err))
		return
	}
	switch baseMsg.Action {
	case ActionChatMessage:
		ctx := context.Background()
		err := mq.GlobalKafka.Publish(ctx, nil, message)
		if err != nil {
			zlog.Error("kafka publish error", zap.Error(err))
		}

	case ActionHeartbeat:
	case ActionAck:
		var ackData AckMessage
		if err := json.Unmarshal(baseMsg.Content, &ackData); err != nil {
			return
		}
		zlog.Info("收到ACK",
			zap.String("msg_id", ackData.MsgId),
			zap.String("user_id", ackData.UserId))
	}
}
func (manager *ClientManager) sendToUser(targetId string, msg []byte) {
	manager.rwLock.RLock()
	defer manager.rwLock.RUnlock()
	if client, ok := manager.Clients[targetId]; ok {
		select {
		case client.Send <- msg:
		default:
			// 缓冲区满了，直接关闭连接，防止阻塞 Manager
			close(client.Send)
			delete(manager.Clients, client.UserId)
		}
	} else {
		zlog.Debug("User offline, cannot send message", zap.String("target", targetId))
	}
}
