package middleware

import (
	"my-chat/internal/api/handler"
	"my-chat/pkg/errno"
	"my-chat/pkg/util/token"
	"my-chat/pkg/zlog"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		zlog.Info("Auth Debug", zap.String("header_token", tokenStr))
		if tokenStr != "" {
			parts := strings.SplitN(tokenStr, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		} else {
			tokenStr = c.Query("token")
			zlog.Info("Auth Debug", zap.String("query_token", tokenStr))
		}
		if tokenStr == "" {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		claims, err := token.ParseAccessToken(tokenStr)
		if err != nil {
			zlog.Error("Token 校验失败",
				zap.String("token_part", tokenStr[0:20]+"..."), // 只打前20位防刷屏
				zap.Error(err))
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		//
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
