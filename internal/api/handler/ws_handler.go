package handler

import (
	"my-chat/internal/websocket"
	"my-chat/pkg/errno"
	"my-chat/pkg/util/token"
	"my-chat/pkg/zlog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	manager       *websocket.ClientManager
	HeartbeatTime int64
}

func NewWSHandler(manager *websocket.ClientManager) *WSHandler {
	return &WSHandler{manager: manager}
}
func (h *WSHandler) Connect(c *gin.Context) {
	userId := c.GetString("userId")
	if userId == "" {
		// WS 连接可能没有经过 middleware.Auth（比如独立挂载），这里做一次兜底鉴权。
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.GetHeader("Authorization")
			if tokenStr != "" {
				parts := strings.SplitN(tokenStr, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenStr = parts[1]
				}
			}
		}
		claims, err := token.ParseAccessToken(tokenStr)
		if err != nil || claims == nil || claims.UserId == "" {
			SendResponse(c, errno.ErrTokenInvalid, nil)
			return
		}
		userId = claims.UserId
		c.Set("userId", userId)
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zlog.Error("webSocket upgrade failed", zap.Error(err))
		return
	}
	client := &websocket.Client{
		Manager:       h.manager,
		Conn:          conn,
		UserId:        userId,
		Send:          make(chan []byte, 256),
		HeartbeatTime: time.Now().Unix(),
	}
	h.manager.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
