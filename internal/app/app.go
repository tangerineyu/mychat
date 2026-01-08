package app

import (
	"context"
	"errors"
	"fmt"
	"my-chat/internal/bootstrap"
	"my-chat/pkg/zlog"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type App struct {
	HTTPServer *http.Server
	WSStart    func() // starts ws manager loops
	Deps       *bootstrap.Deps
}

func (a *App) Run(ctx context.Context) error {
	if a.HTTPServer == nil {
		return fmt.Errorf("nil http server")
	}
	if a.WSStart != nil {
		go a.WSStart()
	}

	errCh := make(chan error, 1)
	go func() {
		zlog.Info("HTTP server starting", zap.String("addr", a.HTTPServer.Addr))
		if err := a.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (a *App) Shutdown(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	// ensure there is a timeout
	_, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	if a.HTTPServer != nil {
		if err := a.HTTPServer.Shutdown(ctx); err != nil {
			zlog.Warn("http shutdown failed", zap.Error(err))
		}
	}
	bootstrap.CloseDeps(ctx, a.Deps, zlog.L)
}
