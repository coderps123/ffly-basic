package service

import (
	"ffly-baisc/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type RolePermissionService struct{}

// SaveRolePermission 更新角色权限
func (service *RolePermissionService) SaveRolePermission(tx *gorm.DB, id uint, permissionIDs []uint) error {
	// 如果传入的权限ID列表为空，则清空该角色的所有权限
	if len(permissionIDs) == 0 {
		if err := tx.Model(&model.RolePermission{}).Where("role_id = ?", id).Delete(&model.RolePermission{}).Error; err != nil {
			return fmt.Errorf("删除角色权限关联失败: %w", err)
		}

		return nil
	}

	// 验证所有的权限ID是否存在
	var count int64
	if err := tx.Model(&model.Permission{}).Where("id in ?", permissionIDs).Count(&count).Error; err != nil {
		return fmt.Errorf("验证权限ID是否存在失败: %w", err)
	}

	// 验证ID列表中存在不存在的权限ID
	if count != int64(len(permissionIDs)) {
		return fmt.Errorf("权限ID列表中存在不存在的权限ID")
	}

	// 删除该角色的所有权限 需要硬删除
	if err := tx.Where("role_id = ?", id).Unscoped().Delete(&model.RolePermission{}).Error; err != nil {
		return fmt.Errorf("删除角色权限关联失败: %w", err)
	}

	// 批量插入角色权限关系
	rolePermissions := make([]model.RolePermission, 0, len(permissionIDs))
	for _, permissionID := range permissionIDs {
		rolePermissions = append(rolePermissions, model.RolePermission{
			RoleID:       id,
			PermissionID: permissionID,
		})
	}
	if err := tx.Create(&rolePermissions).Error; err != nil {
		return fmt.Errorf("创建角色权限关联失败: %w", err)
	}

	return nil
}

// GetRolePermissions 根据角色ID获取权限ids
func (service *RolePermissionService) GetRolePermissionIds(tx *gorm.DB, roleID uint) ([]uint, error) {
	var rolePermissions []model.RolePermission

	if err := tx.Model(&model.RolePermission{}).Where("role_id = ?", roleID).Find(&rolePermissions).Error; err != nil {
		return nil, fmt.Errorf("查询角色权限失败: %w", err)
	}

	rolePermissionsIDs := make([]uint, 0, len(rolePermissions))
	for _, rolePermission := range rolePermissions {
		rolePermissionsIDs = append(rolePermissionsIDs, rolePermission.PermissionID)
	}

	return rolePermissionsIDs, nil
}
