package main

import (
	"context"
	"my-chat/internal/app"
	"my-chat/internal/config"
	"my-chat/pkg/util/snowflake"
	"my-chat/pkg/zlog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	config.InitConfig()
	zlog.Init(config.GlobalConfig.Log)
	defer func() { _ = zlog.L.Sync() }()

	zlog.Info("my chat 正在启动")

	machineID := config.GlobalConfig.App.Machine
	if machineID == 0 {
		zlog.Warn("MachineID is 0 ensure this is intended")
	}
	snowflake.Init(machineID)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, config.GlobalConfig)
	if err != nil {
		zlog.Error("app init failed", zap.Error(err))
		panic(err)
	}

	err = application.Run(ctx)
	zlog.Info("server exiting", zap.Error(err))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	application.Shutdown(shutdownCtx)
}
