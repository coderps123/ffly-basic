package query

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SearchParam 定义查询参数结构体
type SearchParam struct {
	Param string `json:"param"`
	Sign  string `json:"sign"` // 大写 EQ, NEQ, LK, IN, GT, GTE, LT, LTE,
	Val   string `json:"val"`
}

// BuildQuery 构造查询语句
func BuildQuery(db *gorm.DB, searchParamSlice []SearchParam) *gorm.DB {
	for _, s := range searchParamSlice {
		sign := strings.ToUpper(s.Sign)

		switch sign {
		case "EQ":
			// 查询等于某个值
			db = db.Where(s.Param+" = ?", s.Val)
		case "NEQ":
			// 查询不等于某个值
			db = db.Where(s.Param+" <> ?", s.Val)
		case "LK":
			// 模糊查询
			db = db.Where(s.Param+" LIKE ?", "%"+s.Val+"%")
		case "IN":
			// 查询某个值在某个范围内
			db = db.Where(s.Param+" IN (?)", s.Val)
		case "GT":
			// 查询大于某个值
			db = db.Where(s.Param+" > ?", s.Val)
		case "GTE":
			// 查询大于等于某个值
			db = db.Where(s.Param+" >= ?", s.Val)
		case "LT":
			// 查询小于某个值
			db = db.Where(s.Param+" < ?", s.Val)
		case "LTE":
			// 查询小于等于某个值
			db = db.Where(s.Param+" <= ?", s.Val)
		default:
		}
	}
	return db
}

// GetQuery 获取sql查询语句
func GetQuery(c *gin.Context, db *gorm.DB) (*gorm.DB, error) {
	// 判断是否是get请求
	if c.Request.Method != "GET" {
		return nil, fmt.Errorf("请求方式错误，请使用GET请求")
	}

	// 解析搜索参数
	var searchParamSlice []SearchParam
	paramsStr := c.Query("params")
	if paramsStr == "" {
		return db, nil // 如果没有搜索参数，直接返回原始查询
	}

	// 解码 URL 编码的参数
	decodedParams, err := url.QueryUnescape(paramsStr)
	if err != nil {
		return nil, fmt.Errorf("URL解码失败: %v", err)
	}

	if err := json.Unmarshal([]byte(decodedParams), &searchParamSlice); err != nil {
		return nil, fmt.Errorf("搜索参数解析失败: %v, 原始参数: %s", err, decodedParams)
	}

	// 构造查询语句
	query := BuildQuery(db, searchParamSlice)
	if query.Error != nil {
		return nil, query.Error
	}

	return query, nil
}
