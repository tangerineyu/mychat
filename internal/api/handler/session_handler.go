package handler

import (
	"my-chat/internal/service"
	"my-chat/pkg/errno"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}
func (h *SessionHandler) List(c *gin.Context) {
	userId := c.GetString("userId")
	if userId == "" {
		SendResponse(c, errno.ErrTokenInvalid, nil)
		return
	}
	list, err := h.sessionService.GetUserSessions(userId)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, list)
}
