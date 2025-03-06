package model

import (
	types "ffly-baisc/pkg/type"
)

// User 用户模型 -- 查询 只用于查询
type User struct {
	Username  *string      `json:"username,omitempty"`
	Password  *string      `json:"-"` // 不返回给前端, 但是也不从前端接收了
	Nickname  *string      `json:"nickname,omitempty"`
	Email     *string      `json:"email,omitempty"`
	Phone     *string      `json:"phone,omitempty"`
	Status    types.Status `json:"status,omitempty"`
	RoleID    uint         `json:"roleId,omitempty"`
	BaseModel              // 嵌入基础模型
}

// UserCreateRequest 用户创建请求模型 -- 请求入参
type UserCreateRequest struct {
	Username  *string      `json:"username" binding:"required,min=3,max=50"`
	Password  *string      `json:"password" binding:"required"`
	Nickname  *string      `json:"nickname" binding:"omitempty,min=2,max=50"`
	Email     *string      `json:"email" binding:"omitempty,email"`
	Phone     *string      `json:"phone" binding:"omitempty"`
	Status    types.Status `json:"status" gorm:"default:1" binding:"omitempty,oneof=1 2"` // 使用指针以区分是否需要更新
	RoleID    *uint        `json:"roleId"`
	BaseModel              // 嵌入基础模型
}

// UserPatchRequest 用户更新请求模型 -- 部分更新 请求入参
type UserPatchRequest struct {
	Username  *string      `json:"username" binding:"omitempty,min=3,max=50"`
	Nickname  *string      `json:"nickname" binding:"omitempty,min=2,max=50"`
	Email     *string      `json:"email" binding:"omitempty,email"`
	Phone     *string      `json:"phone" binding:"omitempty"`
	Status    types.Status `json:"status" binding:"omitempty,oneof=1 2"` // 使用指针以区分是否需要更新
	RoleID    *uint        `json:"roleId" binding:"omitempty"`
	BaseModel              // 嵌入基础模型
}

// UpdatePasswordRequest 更新密码请求模型 -- 更新密码 请求入参
type UpdatePasswordRequest struct {
	Password        *string `json:"password" binding:"omitempty"`
	NewPassword     *string `json:"newPassword" binding:"omitempty"`
	PasswordConfirm *string `json:"passwordConfirm" binding:"omitempty"`
	BaseModel               // 嵌入基础模型
}

// SimpleQueryFields 简单查询器
func (u *User) SimpleQueryFields() []string {
	return []string{"id", "username"}
}

// TableName 自定义表名
func (u *User) TableName() string {
	return "users"
}
