package model

import (
	types "ffly-baisc/pkg/type"
)

// Role 角色模型 -- 只用于查询
type Role struct {
	Name          string       `json:"name"`
	Code          string       `json:"code"`
	Remark        string       `json:"remark"`
	Status        types.Status `json:"status"`
	PermissionIDs []uint       `json:"permissionIds,omitempty" gorm:"-"` // 权限ID列表，不存储在数据库中
	BaseModel
}

// RoleCreateRequest 创建角色请求模型 -- 请求入参
type RoleCreateRequest struct {
	Name   *string      `json:"name" binding:"required"`
	Code   *string      `json:"code" binding:"required"`
	Remark *string      `json:"remark"`
	Status types.Status `json:"status" gorm:"default:1" binding:"omitempty,oneof=1 2"`
	BaseModel
}

// RolePatchRequest 部分更新角色请求模型 -- 请求入参
type RolePatchRequest struct {
	Name   *string      `json:"name"`
	Code   *string      `json:"code"`
	Remark *string      `json:"remark"`
	Status types.Status `json:"status"`
	BaseModel
}

// TableName 自定义表名
func (r *Role) TableName() string {
	return "roles"
}
