package service

import (
	"ffly-baisc/internal/db"
	"ffly-baisc/internal/model"
	"fmt"
)

// AuthPermissionService 认证权限服务
type AuthPermissionService struct{}

// GetUserRoles 根据用户ID获取用户角色列表
func (s *AuthPermissionService) GetUserRoles(userID uint) ([]uint, error) {
	var userRoles []model.UserRole
	if err := db.DB.MySQL.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, fmt.Errorf("查询用户角色失败: %w", err)
	}

	var roleIDs []uint
	for _, userRole := range userRoles {
		roleIDs = append(roleIDs, userRole.RoleID)
	}

	return roleIDs, nil
}

// GetUserPermissions 根据用户ID获取用户权限列表
func (s *AuthPermissionService) GetUserPermissions(userID uint) ([]*model.Permission, error) {
	// 1. 获取用户角色
	roleIDs, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	if len(roleIDs) == 0 {
		return []*model.Permission{}, nil
	}

	// 2. 根据角色获取权限
	var rolePermissions []*model.RolePermission
	if err := db.DB.MySQL.Where("role_id IN ?", roleIDs).Find(&rolePermissions).Error; err != nil {
		return nil, fmt.Errorf("查询角色权限失败: %w", err)
	}

	var permissionIDs []uint
	for _, rolePermission := range rolePermissions {
		permissionIDs = append(permissionIDs, rolePermission.PermissionID)
	}

	if len(permissionIDs) == 0 {
		return []*model.Permission{}, nil
	}

	// 3. 获取权限详情
	var permissions []*model.Permission
	if err := db.DB.MySQL.Where("id IN ? AND status = 1", permissionIDs).Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("查询权限详情失败: %w", err)
	}

	return permissions, nil
}

// // HasPermission 检查用户是否有指定权限
// func (s *AuthPermissionService) HasPermission(userID uint, permissionPath string) (bool, error) {
// 	userPermissions, err := s.GetUserPermissions(userID)
// 	if err != nil {
// 		return false, err
// 	}

// 	for _, permission := range userPermissions {
// 		if permission == permissionPath {
// 			return true, nil
// 		}
// 	}

// 	return false, nil
// }

// // HasAnyPermission 检查用户是否有任意一个权限
// func (s *AuthPermissionService) HasAnyPermission(userID uint, permissionPaths []string) (bool, error) {
// 	userPermissions, err := s.GetUserPermissions(userID)
// 	if err != nil {
// 		return false, err
// 	}

// 	// 将用户权限转换为map，提高查找效率
// 	permissionMap := make(map[string]bool)
// 	for _, permission := range userPermissions {
// 		permissionMap[permission] = true
// 	}

// 	// 检查是否有任意一个权限
// 	for _, permissionPath := range permissionPaths {
// 		if permissionMap[permissionPath] {
// 			return true, nil
// 		}
// 	}

// 	return false, nil
// }
