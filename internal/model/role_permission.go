package model

import (
	"time"

	"gorm.io/gorm"
)

// 角色权限关联模型
type RolePermission struct {
	ID           uint           `json:"id"`
	RoleID       uint           `json:"role_id"`
	PermissionID uint           `json:"permission_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}

// TableName 自定义表名
func (r *RolePermission) TableName() string {
	return "role_permissions"
}
