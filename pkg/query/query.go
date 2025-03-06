package query

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Query 用于查询
// 本系统默认所有的查询都是get请求
type Query struct {
	Page     int           `json:"page"`     // 页码
	Size     int           `json:"size"`     // 每页显示的数量
	Params   []SearchParam `json:"params"`   // 查询参数
	Simple   bool          `json:"simple"`   // 是否为简单查询(仅返回id，name，value字段)
	Complete bool          `json:"complete"` // 是否为完全查询
}

// SimpleQuerier 定义简单查询的接口
type SimpleQuerier interface {
	SimpleQueryFields() []string // 返回简单查询的字段
}

var simpleCommonFields = []string{"id", "name"}

// GetQuerySQL 获取查询SQL以及分页信息
func GetQuerySQL[T any](c *gin.Context, db *gorm.DB, model T) (*gorm.DB, *Pagination, error) {
	// 获取查询语句
	query, err := GetQuery(c, db)
	if err != nil {
		return nil, nil, err
	}

	// 设置查询模型
	query = query.Model(model)
	if query.Error != nil {
		return nil, nil, query.Error
	}

	// 判断是否为简单查询
	if c.DefaultQuery("simple", "false") == "true" {
		// 简单查询 (仅返回id，name，value字段)

		// 如果没有指定字段，则使用默认字段(id, name)
		fields := simpleCommonFields
		// 如果实现了SimpleQuerier接口, 则使用接口返回的字段
		// 使用 model 而不是 T 进行类型断言，因为 T 是类型参数而不是值
		// 使用 any(model) 将 model 转换为 interface{} 类型再进行类型断言
		if simpleQuerier, ok := any(model).(SimpleQuerier); ok {
			simpleFields := simpleQuerier.SimpleQueryFields()
			if len(simpleFields) > 0 {
				fields = simpleFields
			}
		}

		query = query.Select(fields)
		if query.Error != nil {
			return nil, nil, query.Error
		}
	}

	var p *Pagination
	// 判断是否为完全查询 (是否分页)
	if c.DefaultQuery("complete", "false") != "true" {
		// 获取分页参数 page, size
		p = new(Pagination)
		p.Page = GetPage(c)
		p.Size = GetSize(c)
		p.Total = new(int64)

		// 计算偏移量
		offset := (p.Page - 1) * p.Size
		// 计算每页数量
		limit := p.Size

		// 先计算总数
		if err := query.Count(p.Total).Error; err != nil {
			return nil, nil, err
		}

		// 再设置分页
		query = query.Offset(offset).Limit(limit)
		if query.Error != nil {
			return nil, nil, query.Error
		}
	}

	return query, p, nil
}

// GetQueryData 获取查询数据
// result 必须是一个指针类型, 且元素类型为切片类型
func GetQueryData[T any](db *gorm.DB, c *gin.Context) (*[]*T, *Pagination, error) {
	// 创建一个新的元素实例用于结果输出
	result := new([]*T) // []*model.User, T == model.User

	// 创建一个新的元素实例用于Model
	model := new(T)

	// 获取查询SQL
	query, p, err := GetQuerySQL(c, db, model)
	if err != nil {
		return nil, nil, err
	}

	// 执行查询
	if err := query.Find(result).Error; err != nil {
		return nil, nil, err
	}

	return result, p, nil
}
