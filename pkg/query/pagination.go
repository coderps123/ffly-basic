package query

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Total *int64 `json:"total"`
}

// GetPage 获取分页页码
func GetPage(c *gin.Context) int {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	return page
}

// GetSize 获取分页大小
func GetSize(c *gin.Context) int {
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if size < 1 {
		size = 10
	}

	return size
}

// GetPageInfo 获取分页信息
func GetPageInfo(c *gin.Context) *Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	return &Pagination{
		Page: page,
		Size: size,
	}
}

// GetListByPage 获取分页列表
func GetListByPage(db *gorm.DB, model any, p *Pagination) error {
	// 计算偏移量
	offset := (p.Page - 1) * p.Size

	// 计算每页数量
	limit := p.Size

	if err := db.Model(model).Count(p.Total).Offset(offset).Limit(limit).Find(model).Error; err != nil {
		// 处理错误
		return err
	}

	return nil
}

// GetListByContext 获取分页列表 --- 主要的
func GetListByContext(db *gorm.DB, model any, c *gin.Context) (*Pagination, error) {
	// 设置分页参数
	p := GetPageInfo(c)

	// 查询权限列表
	if err := GetListByPage(db, model, p); err != nil {
		return nil, err
	}

	return p, nil
}
