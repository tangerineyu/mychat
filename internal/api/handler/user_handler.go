package handler

import (
	"my-chat/pkg/errno"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"my-chat/internal/service"
	"my-chat/pkg/zlog"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type RegisterRequest struct {
	Telephone string `json:"telephone" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Nickname  string `json:"nickname" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.userService.Register(req.Telephone, req.Password, req.Nickname); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "注册成功"})
}

type LoginRequest struct {
	Telephone string `json:"telephone" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Warn("Login bind failed", zap.Error(err))
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	accessToken, refreshToken, user, err := h.userService.Login(req.Telephone, req.Password)
	if err != nil {
		zlog.Warn("Login failed",
			zap.String("telephone", req.Telephone),
			zap.Error(err),
		)
		SendResponse(c, err, nil)
		return
	}
	zlog.Info("Login success", zap.String("uuid", user.Uuid))
	SendResponse(c, nil, gin.H{
		"token":         accessToken,
		"refresh_token": refreshToken,
		"nickname":      user.Nickname,
	})
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Warn("RefreshToken bind failed", zap.Error(err))
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	newAccess, newRefresh, err := h.userService.RefreshToken(req.RefreshToken)
	if err != nil {
		//刷新失败，前端应强制跳转回登录页
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": err.Error()})
		
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "刷新成功",
		"data": gin.H{
			"access_token":  newAccess,
			"refresh_token": newRefresh,
		},
	})
}
