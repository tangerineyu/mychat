package main

import (
	"my-chat/internal/api/handler"
	"my-chat/internal/api/middleware"
	"my-chat/internal/config"
	"my-chat/internal/dao"
	"my-chat/internal/mq"
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
	defer func() {
		_ = zlog.L.Sync()
	}()
	zlog.Info("my chat 正在启动")
	dao.InitDB()
	dao.InitRedis()
	mq.InitKafka()

	userRepo := repo.NewUserRepository(dao.DB)
	msgRepo := repo.NewMessageRepository(dao.DB)
	groupRepo := repo.NewGroupRepository(dao.DB, dao.RDB)
	contactRepo := repo.NewContactRepository(dao.DB)
	sessionRepo := repo.NewSessionRepository(dao.DB)

	userService := service.NewUserService(userRepo)
	chatService := service.NewChatService(msgRepo, groupRepo)
	groupService := service.NewGroupService(groupRepo, userRepo)
	contactService := service.NewContactService(contactRepo, userRepo)
	sessionService := service.NewSessionService(sessionRepo, groupRepo, userRepo)

	wsManager := websocket.NewClientManager(chatService, sessionRepo)
	go wsManager.Start()
	go wsManager.StartConsumer()
	userHandler := handler.NewUserHandler(userService)
	wsHandler := handler.NewWSHandler(wsManager)
	groupHandler := handler.NewGroupHandler(groupService)
	chatHandler := handler.NewChatHandler(chatService)
	contactHandler := handler.NewContactHandler(contactService)
	sessionHandler := handler.NewSessionHandler(sessionService)

	r := gin.New()
	r.Use(middleware.GinLogger())
	r.Use(gin.Recovery())

	r.Static("/static", "./static")

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.POST("/refresh-token", userHandler.RefreshToken)
		v1.GET("/ws", wsHandler.Connect)
		v1.POST("/group/create", groupHandler.Create)
		v1.POST("/group/join", groupHandler.Join)
		v1.POST("/chat/history", chatHandler.History)
		v1.POST("/upload/avatar", userHandler.UploadAvatar)
		v1.POST("/user/updateUserInfo", userHandler.UpdateUserInfo)
		v1.POST("/group/getGroupInfo", groupHandler.GetGroupInfo)
		v1.POST("/group/getGroupMemberList", groupHandler.GetGroupMemberList)
		v1.POST("/contact/add", contactHandler.AddFriend)
		v1.POST("/contact/agree", contactHandler.AgreeFriend)
		v1.POST("/contact/list", contactHandler.GetContactList)
		v1.POST("/contact/applyList", contactHandler.GetApplyList)
		v1.POST("/session/list", sessionHandler.List)
	}
	zlog.Info("服务器启动成功", zap.String("port", "8080"))
	if err := r.Run(":8080"); err != nil {
		zlog.Error("服务启动失败", zap.String("error", err.Error()))
		panic(err)
	}
}
