package routes

import (
	"ffly-baisc/internal/handler"

	"github.com/gin-gonic/gin"
)

func ResigterApiLogRouter(g *gin.RouterGroup) {
	group := g.Group("/api_log")
	{
		group.GET("", handler.GetApiLogList)
	}
}
