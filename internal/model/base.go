package model

import (
	"time"

	"gorm.io/gorm"
)

// Base 基础模型
type Base struct {
	ID        uint64          `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"-"` // json 中隐藏删除时间
}
