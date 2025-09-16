package middleware

import (
	"net/http"

	"vue2-blog-server/internal/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("Panic recovered",
				zap.String("error", err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
		} else {
			logger.Error("Panic recovered",
				zap.Any("error", recovered),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
		}

		pkg.ServerError(c, "服务器内部错误")
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}