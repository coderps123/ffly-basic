package router

import (
	"ffly-baisc/internal/config"
	"ffly-baisc/internal/middleware"
	"ffly-baisc/internal/router/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(gin.DebugMode) // 设置运行模式
	r := gin.Default()

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 公开路由
		public := v1.Group("")
		// 注册登录路由
		routes.ResigterLoginRouter(public)

		// 需要认证的路由
		authGroup := v1.Group("")
		// 注册中间件
		authGroup.Use(middleware.Auth())

		// 注册用户路由
		routes.ResigterUserRouter(authGroup)
		// 注册角色路由
		routes.ResigterRoleRouter(authGroup)
		// 注册权限路由
		routes.ResigterPermissionRouter(authGroup)
	}

	r.Run(fmt.Sprintf(":%d", config.GlobalConfig.App.Port)) // 监听端口
}
