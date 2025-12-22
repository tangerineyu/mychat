package main

import (
	"my-chat/internal/api/handler"
	"my-chat/internal/api/middleware"
	"my-chat/internal/config"
	"my-chat/internal/dao"
	"my-chat/internal/repo"
	"my-chat/internal/service"
	"my-chat/internal/websocket"
	"my-chat/pkg/zlog"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config.InitConfig()
	zlog.Init(config.GlobalConfig.Log)
	defer zlog.L.Sync()
	zlog.Info("my chat 正在启动")
	dao.InitDB()

	userRepo := repo.NewUserRepository(dao.DB)
	msgRepo := repo.NewMessageRepository(dao.DB)
	groupRepo := repo.NewGroupRepository(dao.DB)

	userService := service.NewUserService(userRepo)
	chatService := service.NewChatService(msgRepo, groupRepo)
	wsManager := websocket.NewClientManager(chatService)
	go wsManager.Start()
	userHandler := handler.NewUserHandler(userService)
	wsHandler := handler.NewWSHandler(wsManager)

	r := gin.New()
	r.Use(middleware.GinLogger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.POST("/refresh-token", userHandler.RefreshToken)
		v1.GET("/ws", wsHandler.Connect)
	}
	zlog.Info("服务器启动成功", zap.String("port", "8080"))
	if err := r.Run(":8080"); err != nil {
		zlog.Error("服务启动失败", zap.String("error", err.Error()))
		panic(err)
	}
}
