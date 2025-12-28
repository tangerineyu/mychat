package handler

import (
	"my-chat/internal/service"
	"my-chat/pkg/errno"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

type HistoryReq struct {
	TargetId string `json:"target_id" binding:"required"`
	ChatType int    `json:"type" binding:"required"` //1-私聊 2-群聊
}

func (h *ChatHandler) History(c *gin.Context) {
	var req HistoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	//currentUserId := c.Query("uid")
	val, exists := c.Get("userId")
	if !exists {
		SendResponse(c, errno.ErrTokenInvalid, nil)
		return
	}
	currentUserId := val.(string)
	if currentUserId == "" {
		SendResponse(c, errno.ErrTokenInvalid, nil)
		return
	}
	messages, err := h.chatService.GetHistory(currentUserId, req.TargetId, req.ChatType)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, messages)
}
