package model

import (
	"gorm.io/gorm"

	types "ffly-baisc/pkg/type"
)

// User 用户模型 -- 创建、查询
type User struct {
	Username   *string      `json:"username" binding:"required,min=3,max=50"`
	Password   *string      `json:"-" binding:"omitempty"` // 不返回给前端, 但是也不从前端接收了
	Nickname   *string      `json:"nickname" binding:"omitempty,min=2,max=50"`
	Email      *string      `json:"email" binding:"omitempty,email"`
	Phone      *string      `json:"phone" binding:"omitempty"`
	Status     types.Status `json:"status" gorm:"default:1" binding:"omitempty,oneof=1 2"` // 使用指针以区分是否需要更新
	RoleID     uint         `json:"roleId"`
	gorm.Model              // 嵌入基础模型
}

// UserCreateRequest 用户模型 -- 创建、查询
type UserCreateRequest struct {
	Username   *string      `json:"username" binding:"required,min=3,max=50"`
	Password   *string      `json:"password" binding:"required"`
	Nickname   *string      `json:"nickname" binding:"omitempty,min=2,max=50"`
	Email      *string      `json:"email" binding:"omitempty,email"`
	Phone      *string      `json:"phone" binding:"omitempty"`
	Status     types.Status `json:"status" gorm:"default:1" binding:"omitempty,oneof=1 2"` // 使用指针以区分是否需要更新
	RoleID     uint         `json:"roleId"`
	gorm.Model              // 嵌入基础模型
}

type UserUpdateRequest struct {
	Username   *string      `json:"username" binding:"required,min=3,max=50"`
	Nickname   *string      `json:"nickname" binding:"omitempty,min=2,max=50"`
	Email      *string      `json:"email" binding:"omitempty,email"`
	Phone      *string      `json:"phone" binding:"omitempty"`
	Status     types.Status `json:"status" gorm:"default:1" binding:"omitempty,oneof=1 2"` // 使用指针以区分是否需要更新
	RoleID     uint         `json:"roleId" binding:"omitempty"`
	gorm.Model              // 嵌入基础模型
}

type UserPatchStatusRequest struct {
	Status *types.Status `json:"status" binding:"oneof=1 2"`
}

// TableName 自定义表名
func (u *User) TableName() string {
	return "users"
}
