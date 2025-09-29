package model

import (
	types "ffly-baisc/pkg/type"
)

// Permission 权限模型
type Permission struct {
	Title string `json:"title"` // 权限标题 --- menu
	Name  string `json:"name"`
	// Type      string        `json:"type"`      // 权限类型，数据库层面限制用户输入必须为 “menu / button"
	Path string `json:"path"` // 路由路径 --- menu
	// Code      string        `json:"code"`      // 权限码 --- button
	Component string        `json:"component"` // 路由组件名称 --- menu
	Redirect  string        `json:"redirect"`  // 重定向路径 --- menu
	Visible   bool          `json:"visible"`   //
	Icon      string        `json:"icon"`      // 图标 --- menu
	Sort      int           `json:"sort"`      // 排序 --- menu
	ParentID  uint          `json:"parentId"`  // 父级权限ID --- menu
	Remark    string        `json:"remark"`
	Status    types.Status  `json:"status"` // 1:启用 2:禁用
	Buttons   string        `json:"buttons"`
	Params    string        `json:"params"`
	Children  []*Permission `json:"children,omitempty" gorm:"-"` // 子权限列表
	BaseModel
}

// PermissionCreatedRequest 权限创建请求
type PermissionCreatedRequest struct {
	Title string `json:"title" binding:"required"` // 权限标题 --- menu
	Name  string `json:"name"`
	// Type      *string      `json:"type" binding:"required"` // 权限类型，数据库层面限制用户输入必须为 “menu / button"
	Path string `json:"path"` // 路由路径 --- menu
	// Code      *string      `json:"code"`                    // 权限码 --- button
	Redirect  string       `json:"redirect"`  // 重定向路径 --- menu
	Visible   bool         `json:"visible"`   //
	Component string       `json:"component"` // 路由组件名称 --- menu
	Icon      string       `json:"icon"`      // 图标 --- menu
	Sort      int          `json:"sort"`      // 排序 --- menu
	ParentID  uint         `json:"parentId" ` // 父级权限ID --- menu
	Remark    string       `json:"remark"`
	Buttons   string       `json:"buttons"`
	Params    string       `json:"params"`
	Status    types.Status `json:"status"` // 1:启用 2:禁用
	BaseModel
}

// PermissionPatchRequest 权限更新请求
type PermissionPatchRequest struct {
	Title string  `json:"title"` // 权限标题 --- menu
	Name  *string `json:"name"`
	// Type      *string      `json:"type"`      // 权限类型，数据库层面限制用户输入必须为 “menu / button"
	Path *string `json:"path"` // 路由路径 --- menu
	// Code      *string      `json:"code"`      // 权限码 --- button
	Component *string      `json:"component"` // 路由组件名称 --- menu
	Redirect  string       `json:"redirect"`  // 重定向路径 --- menu
	Visible   bool         `json:"visible"`   // 1: 可见 2: 不可见 ---
	Icon      *string      `json:"icon"`      // 图标 --- menu
	Sort      int          `json:"sort"`      // 排序 --- menu
	ParentID  uint         `json:"parentId"`  // 父级权限ID --- menu
	Remark    *string      `json:"remark"`
	Buttons   string       `json:"buttons"`
	Params    string       `json:"params"`
	Status    types.Status `json:"status"` // 1:启用 2:禁用
	BaseModel
}

// TableName 表名
func (p *Permission) TableName() string {
	return "permissions"
}
