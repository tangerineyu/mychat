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
	db, err := dao.NewMySQL(&config.GlobalConfig.MySQL)
	if err != nil {
		zlog.Error("数据库初始化失败", zap.Error(err))
		panic(err)
	}
	rdb, err := dao.NewRedis(&config.GlobalConfig.Redis)
	if err != nil {
		zlog.Error("Redis初始化失败", zap.Error(err))
		panic(err)
	}
	kafkaClient := mq.NewKafkaClient(&config.GlobalConfig.Kafka)
	defer kafkaClient.Close()

	userRepo := repo.NewUserRepository(db)
	msgRepo := repo.NewMessageRepository(db)
	groupRepo := repo.NewGroupRepository(db, rdb)
	contactRepo := repo.NewContactRepository(db)
	sessionRepo := repo.NewSessionRepository(db, rdb)
	adminRepo := repo.NewAdminRepository(db)

	userService := service.NewUserService(userRepo)
	chatService := service.NewChatService(msgRepo, groupRepo)
	groupService := service.NewGroupService(groupRepo, userRepo)
	contactService := service.NewContactService(contactRepo, userRepo)
	sessionService := service.NewSessionService(sessionRepo, groupRepo, userRepo)
	adminService := service.NewAdminService(adminRepo)

	wsManager := websocket.NewClientManager(chatService, sessionRepo, kafkaClient)
	go wsManager.Start()
	go wsManager.StartConsumer()
	userHandler := handler.NewUserHandler(userService)
	wsHandler := handler.NewWSHandler(wsManager)
	groupHandler := handler.NewGroupHandler(groupService)
	chatHandler := handler.NewChatHandler(chatService)
	contactHandler := handler.NewContactHandler(contactService)
	sessionHandler := handler.NewSessionHandler(sessionService)
	adminHandler := handler.NewAdminHandler(adminService)

	r := gin.New()
	r.Use(middleware.GinLogger())
	r.Use(gin.Recovery())

	r.Static("/static", "./static")

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.POST("/refresh-token", userHandler.RefreshToken)

	}
	authGroup := v1.Group("/")
	authGroup.Use(middleware.Auth())
	{
		authGroup.GET("/ws", wsHandler.Connect)
		//用户相关
		authGroup.POST("/upload/avatar", userHandler.UploadAvatar)
		authGroup.POST("/user/updateUserInfo", userHandler.UpdateUserInfo)
		//群组相关
		authGroup.POST("/group/create", groupHandler.Create)
		authGroup.POST("/group/join", groupHandler.Join)
		authGroup.POST("/group/getGroupInfo", groupHandler.GetGroupInfo)
		authGroup.POST("/group/getGroupMemberList", groupHandler.GetGroupMemberList)
		authGroup.POST("/group/loadMyGroup", groupHandler.LoadMyJoinedGroup)
		authGroup.POST("/group/leaveGroup", groupHandler.LeaveGroup)
		authGroup.POST("/group/kickGroupMember", groupHandler.KickGroupMember)
		authGroup.POST("/group/dismissGroup", groupHandler.DismissGroup)
		//联系人相关
		authGroup.POST("/contact/add", contactHandler.AddFriend)
		authGroup.POST("/contact/agree", contactHandler.AgreeFriend)
		authGroup.POST("/contact/list", contactHandler.GetContactList)
		authGroup.POST("/contact/applyList", contactHandler.GetApplyList)
		authGroup.POST("/contact/refuseContactApply", contactHandler.RefuseApply)
		authGroup.POST("/contact/deleteContact", contactHandler.DeleteContact)
		authGroup.POST("/contact/blackContact", contactHandler.BlackContact)
		authGroup.POST("/contact/cancelBlackContact", contactHandler.UnBlackContact)
		//聊天历史记录
		authGroup.POST("/chat/history", chatHandler.History)
		//会话接口
		authGroup.POST("/session/list", sessionHandler.List)
		// Admin User
		authGroup.POST("/user/getUserInfoList", adminHandler.GetUserList)
		authGroup.POST("/user/disableUsers", adminHandler.DisableUser)
		authGroup.POST("/user/ableUsers", adminHandler.AbleUser)
		// Admin Group
		authGroup.POST("/group/getGroupInfoList", adminHandler.GetGroupList)
		authGroup.POST("/group/disableGroup", adminHandler.DisableGroup)
		authGroup.POST("/group/ableGroups", adminHandler.AbleGroup)
	}
	zlog.Info("服务器启动成功", zap.String("port", "8080"))
	if err := r.Run(":8080"); err != nil {
		zlog.Error("服务启动失败", zap.String("error", err.Error()))
		panic(err)
	}
}
