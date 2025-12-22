package websocket

import (
	"encoding/json"
	"my-chat/internal/service"
	"my-chat/pkg/zlog"
	"sync"

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
}

func NewClientManager(chatService *service.ChatService) *ClientManager {
	return &ClientManager{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan []byte),
		Clients:     make(map[string]*Client),
		chatService: chatService,
	}
}
func (manager *ClientManager) Start() {
	zlog.Info("Websocket Client Manager Started")
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
	var rawMsg Message
	if err := json.Unmarshal(message, &rawMsg); err != nil {
		zlog.Error("Failed to unmarshal message", zap.String("message", string(message)), zap.Error(err))
		return
	}
	if rawMsg.Action == ActionChatMessage {
		var chatData ChatMessageContent
		if err := json.Unmarshal(rawMsg.Content, &chatData); err != nil {
			zlog.Error("Content Parse Error", zap.Error(err))
			return
		}
		//调用Service存消息， 获取标准发送格式
		jsonBytes, err := manager.chatService.SaveAndFactory(
			chatData.SendId,
			chatData.ReceiverId,
			chatData.Content,
			chatData.Type,
			1,
		)
		if err != nil {
			zlog.Error("SaveAndFactory Error", zap.Error(err))
			return
		}
		//单聊
		if chatData.Type == 1 {
			manager.sendToUser(chatData.ReceiverId, jsonBytes)
			manager.sendToUser(chatData.SendId, jsonBytes)
		} else if chatData.Type == 2 {
			zlog.Info("暂未实现")
		}
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
