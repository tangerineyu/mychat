package handler

import (
	"my-chat/internal/websocket"
	"my-chat/pkg/errno"
	"my-chat/pkg/zlog"
	"net/http"

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
	manager *websocket.ClientManager
}

func NewWSHandler(manager *websocket.ClientManager) *WSHandler {
	return &WSHandler{manager: manager}
}
func (h *WSHandler) Connect(c *gin.Context) {
	userId := c.GetString("userId")
	if userId == "" {
		SendResponse(c, errno.ErrTokenInvalid, nil)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zlog.Error("webSocket upgrade failed", zap.Error(err))
		return
	}
	client := &websocket.Client{
		Manager: h.manager,
		Conn:    conn,
		UserId:  userId,
		Send:    make(chan []byte, 256),
	}
	h.manager.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
