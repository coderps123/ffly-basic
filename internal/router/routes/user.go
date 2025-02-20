package routes

import (
	"ffly-baisc/internal/handler"

	"github.com/gin-gonic/gin"
)

func ResigterUserRouter(g *gin.RouterGroup) {
	group := g.Group("/user")
	{
		// 如果要这样写，那么 /info 就必须在 /:id 之前，否则会匹配到 /:id 路由
		group.GET("/info", handler.GetCurrentUserInfo)
		// 修改密码
		group.PATCH("/:id/password", handler.UpdateUserPassword)

		group.GET("", handler.GetUserList)
		group.GET("/:id", handler.GetUser)
		group.POST("", handler.CreateUser)
		group.PATCH("/:id", handler.PatchUser)
		group.DELETE("/:id", handler.DeleteUser)
	}
}
