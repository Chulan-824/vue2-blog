package middleware

import (
	"strings"

	"vue2-blog-server/internal/auth"
	"vue2-blog-server/internal/pkg"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header或Cookie中获取token
		token := extractToken(c)
		if token == "" {
			pkg.AuthError(c, "未提供认证令牌")
			c.Abort()
			return
		}

		// 解析token
		claims, err := jwtManager.ParseToken(token)
		if err != nil {
			pkg.AuthError(c, "无效的认证令牌")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("admin", claims.Admin)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		admin, exists := c.Get("admin")
		if !exists {
			pkg.AuthError(c, "未找到用户权限信息")
			c.Abort()
			return
		}

		if !admin.(bool) {
			pkg.AuthError(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求登录）
func OptionalAuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token != "" {
			claims, err := jwtManager.ParseToken(token)
			if err == nil {
				// 将用户信息存入上下文
				c.Set("userID", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("admin", claims.Admin)
			}
		}
		c.Next()
	}
}

// extractToken 提取token
func extractToken(c *gin.Context) string {
	// 首先从Authorization Header中获取
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:6]) == "BEARER" {
		return bearerToken[7:]
	}

	// 从Cookie中获取
	token, err := c.Cookie("token")
	if err == nil {
		return token
	}

	// 从query参数中获取
	return c.Query("token")
}