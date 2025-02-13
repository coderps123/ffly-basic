package model

import (
	"time"

	"gorm.io/gorm"
)

type RoleStatusType int

const (
	RoleStatusActive   RoleStatusType = 1 // 1: active
	RoleStatusInactive RoleStatusType = 0 // 0: inactive
)

// PatchRoleRequest 更新角色状态请求
type PatchRoleRequest struct {
	Status *RoleStatusType `json:"status" binding:"required,oneof=1 0"` // 1:启用 0:禁用
}

// Role 角色模型
type Role struct {
	gorm.Model
	ID        uint            `json:"id"`
	Name      string          `json:"name" binding:"required"`
	Code      string          `json:"code" binding:"required"`
	Remark    string          `json:"remark"`
	Status    *RoleStatusType `json:"status"` // 1:启用 0:禁用
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `json:"-"`
}

// TableName 自定义表名
func (r *Role) TableName() string {
	return "roles"
}
