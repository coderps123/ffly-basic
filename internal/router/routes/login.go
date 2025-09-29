package routes

import (
	"ffly-baisc/internal/handler"

	"github.com/gin-gonic/gin"
)

// 注册路由
func ResigterLoginRouter(group *gin.RouterGroup) {
	// 用户注册
	group.POST("/register", handler.Register)
	// 用户登录
	group.POST("/login", handler.Login)
	// 刷新 Token
	group.POST("/refresh", handler.RefreshToken)
}
