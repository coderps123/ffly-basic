package handler

import (
	"ffly-baisc/internal/service"
	"ffly-baisc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetApiLogList 获取用户列表
func GetApiLogList(c *gin.Context) {
	var apiLogService service.ApiLogService

	apiLogs, pagination, err := apiLogService.GetApiLogList(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取日志列表失败", err)
		return
	}

	response.SuccessWithPagination(c, apiLogs, pagination, "日志列表获取成功")
}
