package routes

import (
	"ffly-baisc/internal/handler"

	"github.com/gin-gonic/gin"
)

func ResigterPermissionRouter(g *gin.RouterGroup) {
	group := g.Group("/permission")
	{
		group.GET("", handler.GetPermissionList)
		group.GET("/:id", handler.GetPermission)
		group.POST("", handler.CreatePermission)
		group.PATCH("/:id", handler.PatchPermission)
		group.DELETE("/:id", handler.DeletePermission)
	}
}
