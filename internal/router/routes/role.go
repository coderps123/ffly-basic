package routes

import (
	"ffly-baisc/internal/handler"

	"github.com/gin-gonic/gin"
)

func ResigterRoleRouter(g *gin.RouterGroup) {
	group := g.Group("/role")
	{
		group.GET("", handler.GetRoleList)
		group.GET("/:id", handler.GetRole)
		group.POST("", handler.CreateRole)
		group.PATCH("/:id", handler.PatchRole)
		// 更新角色权限
		group.PATCH("/:id/permissions", handler.PatchRolePermissions)
		group.DELETE("/:id", handler.DeleteRole)
	}
}
