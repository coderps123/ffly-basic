package middleware

import (
	"ffly-baisc/pkg/auth"
	"ffly-baisc/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "未登录或非法访问", nil)
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, "请求头中 Authorization 格式有误", nil)
			c.Abort()
			return
		}

		// 解析 token
		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "无效的 Token", err)
			c.Abort()
			return
		}

		// 验证是否为 Access Token
		if claims.TokenType != "access" {
			response.Error(c, http.StatusUnauthorized, "Token 类型错误，需要 Access Token", nil)
			c.Abort()
			return
		}

		// 将当前请求的用户信息保存到请求的上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
