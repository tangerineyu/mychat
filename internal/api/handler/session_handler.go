package handler

import (
	"github.com/gin-gonic/gin"
	"my-chat/internal/service"
	"my-chat/pkg/errno"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}
func (h *SessionHandler) List(c *gin.Context) {
	userId := c.Query("uid")
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
