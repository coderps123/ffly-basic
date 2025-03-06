package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint            `json:"id,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `json:"-"` // json 中隐藏删除时间
}
