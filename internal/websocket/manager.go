package websocket

import (
	"encoding/json"
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
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		Clients:    make(map[string]*Client),
	}
}
func (manager *ClientManager) Start() {
	zlog.Info("Websocket Client Manager Started")
	for {
		select {
		case client := <-manager.Register:
			zlog.Info("New connection",
				zap.String("uuid", client.UserId))
			if oldClient, ok := manager.Clients[client.UserId]; ok {
				close(oldClient.Send)
				delete(manager.Clients, client.UserId)
			}
			manager.Clients[client.UserId] = client
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client.UserId]; ok {
				zlog.Info("Disconnect", zap.String("uuid", client.UserId))
				delete(manager.Clients, client.UserId)
				close(client.Send)
			}
		case message := <-manager.Broadcast:
			var msgObj Message
			json.Unmarshal(message, &msgObj)
			//		manager.sendToAll(message)
		}
	}
}
