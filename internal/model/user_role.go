package model

// 用户角色关联模型
type UserRole struct {
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
	BaseModel
}

// TableName 自定义表名
func (r *UserRole) TableName() string {
	return "user_roles"
}
