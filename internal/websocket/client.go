package websocket

import (
	"my-chat/pkg/zlog"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait      = 10 * time.Second //写超时
	pongWait       = 60 * time.Second //心跳超时
	PingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Client代表一个WebSocket连接用户
type Client struct {
	Manager *ClientManager  //客户端管理器
	Conn    *websocket.Conn //实际的ws连接
	UserId  string          //用户ID
	Send    chan []byte     //发送缓冲通道
}

func (c *Client) ReadPump() {
	defer func() {
		c.Manager.Unregister <- c
		//关闭底层的websocket连接
		c.Conn.Close()
	}()
	//设置最大消息大小和读超时
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zlog.Error("WS read error", zap.Error(err))
			}
			break
		}
		c.Manager.Broadcast <- message
	}
}
func (c *Client) WritePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
