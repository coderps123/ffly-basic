package model

// 角色权限关联模型
type RolePermission struct {
	RoleID       uint `json:"roleId"`
	PermissionID uint `json:"permission_id"`
	BaseModel
}

// TableName 自定义表名
func (r *RolePermission) TableName() string {
	return "role_permissions"
}

// RolePermissionUpdateRequest 角色权限更新请求模型
type RolePermissionUpdateRequest struct {
	PermissionIDs []uint `json:"permissionIds" binding:"required"`
}
