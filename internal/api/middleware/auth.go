package middleware

import (
	"github.com/gin-gonic/gin"
	"my-chat/internal/api/handler"
	"my-chat/pkg/errno"
	"my-chat/pkg/util/token"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr != "" {
			parts := strings.SplitN(tokenStr, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		} else {
			tokenStr = c.Query("token")
		}
		if tokenStr == "" {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		claims, err := token.ParseAccessToken(tokenStr)
		if err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		c.Set("userId", claims)
		c.Next()
	}
}
