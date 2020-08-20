package middleware

import (
	"fmt"
	"github.com/wlcy/tradehome-service/common/errno"
	"github.com/wlcy/tradehome-service/common/token"
	"github.com/wlcy/tradehome-service/handler"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := token.ParseRequest(c); err != nil {
			fmt.Printf("AuthMiddleware error: %s\n", err)
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
