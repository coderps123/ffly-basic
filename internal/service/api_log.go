package service

import (
	"ffly-baisc/internal/db"
	"ffly-baisc/internal/model"
	"ffly-baisc/pkg/pagination"

	"github.com/gin-gonic/gin"
)

type ApiLogService struct{}

// GetApiLogList 获取用户列表
func (service *ApiLogService) GetApiLogList(c *gin.Context) ([]*model.ApiLog, *pagination.Pagination, error) {
	var apiLogs []*model.ApiLog

	// 查询权限列表
	pagination, err := pagination.GetListByContext(db.DB.MySQL, &apiLogs, c)
	if err != nil {
		return nil, nil, err
	}

	return apiLogs, pagination, nil
}

// CreateApiLog 创建api日志
func (service *ApiLogService) CreateApiLog(apiLog *model.ApiLog) error {
	if err := db.DB.MySQL.Create(apiLog).Error; err != nil {
		return err
	}

	return nil
}
