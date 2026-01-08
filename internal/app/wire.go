package app

import (
	"context"
	"fmt"
	"my-chat/internal/api/handler"
	"my-chat/internal/api/middleware"
	"my-chat/internal/api/router"
	"my-chat/internal/bootstrap"
	"my-chat/internal/config"
	"my-chat/internal/repo"
	"my-chat/internal/service"
	"my-chat/internal/websocket"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	_ = ctx
	if cfg == nil {
		return nil, fmt.Errorf("nil config")
	}

	deps, err := bootstrap.InitDeps(cfg)
	if err != nil {
		return nil, err
	}

	// repos
	userRepo := repo.NewUserRepository(deps.DB)
	msgRepo := repo.NewMessageRepository(deps.DB)
	groupRepo := repo.NewGroupRepository(deps.DB, deps.Redis)
	contactRepo := repo.NewContactRepository(deps.DB)
	sessionRepo := repo.NewSessionRepository(deps.DB, deps.Redis)
	adminRepo := repo.NewAdminRepository(deps.DB)

	// services
	userService := service.NewUserService(userRepo)
	chatService := service.NewChatService(msgRepo, groupRepo)
	groupService := service.NewGroupService(groupRepo, userRepo)
	contactService := service.NewContactService(contactRepo, userRepo)
	sessionService := service.NewSessionService(sessionRepo, groupRepo, userRepo)
	adminService := service.NewAdminService(adminRepo)

	// websocket manager
	wsManager := websocket.NewClientManager(chatService, sessionRepo, deps.Kafka)
	wsStart := func() {
		// Start() already starts consumer/heartbeat internally.
		wsManager.Start()
	}

	// handlers
	userHandler := handler.NewUserHandler(userService)
	wsHandler := handler.NewWSHandler(wsManager)
	groupHandler := handler.NewGroupHandler(groupService)
	chatHandler := handler.NewChatHandler(chatService)
	contactHandler := handler.NewContactHandler(contactService)
	sessionHandler := handler.NewSessionHandler(sessionService)
	adminHandler := handler.NewAdminHandler(adminService)

	// gin engine
	r := gin.New()
	r.Use(middleware.GinLogger())
	r.Use(gin.Recovery())
	r.Static("/static", "./static")
	router.Register(r, userHandler, wsHandler, groupHandler, chatHandler, contactHandler, sessionHandler, adminHandler)

	port := cfg.App.Port
	addr := ":" + strconv.FormatInt(port, 10)

	httpSrv := &http.Server{Addr: addr, Handler: r}

	return &App{HTTPServer: httpSrv, WSStart: wsStart, Deps: deps}, nil
}
