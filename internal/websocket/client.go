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
	Manager *ClientManager  //客户端管理器，读到消息后广播， 断开注销
	Conn    *websocket.Conn //实际的ws连接
	UserId  string          //用户ID，这个连接属于谁
	Send    chan []byte     //发送缓冲通道
}

// ReadPump负责从WebSocket连接中读取消息，检查客户端是不是活着
func (c *Client) ReadPump() {
	defer func() {
		c.Manager.Unregister <- c
		//关闭底层的websocket连接
		c.Conn.Close()
	}()
	//设置最大消息大小和读超时
	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return
	}
	//收到pong消息后更新读超时
	c.Conn.SetPongHandler(func(string) error {
		if err = c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			return err
		}
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

// 把Send通道里的数据写给客户端 定时发送ping消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		//c.Send有消息要发送
		case message, ok := <-c.Send:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if !ok {
				//忽略err，因为本来就要退出了
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			//w.Write(message)
			if _, err := w.Write(message); err != nil {
				return
			}
			//如果 Send 通道里堆积了好几条消息，与其发好几次 TCP 包，
			//不如一次性全部读出来，合并到一个 Writer 里发出去，减少网络开销。
			n := len(c.Send)
			for i := 0; i < n; i++ {
				//w.Write(<-c.Send)
				if _, err := w.Write(<-c.Send); err != nil {
					return
				}
			}
			if err := w.Close(); err != nil {
				return
			}
		//发送ping消息
		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			//客户端收到后会回复一个pong消息，触发ReadPump里的pong处理器
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
