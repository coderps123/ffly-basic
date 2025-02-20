package service

import (
	"ffly-baisc/internal/model"
	"ffly-baisc/internal/mysql"
	"ffly-baisc/pkg/pagination"
	"fmt"

	"github.com/gin-gonic/gin"
)

type PermissionService struct{}

// GetPermissionList 获取权限列表
func (service *PermissionService) GetPermissionList(c *gin.Context) ([]*model.Permission, *pagination.Pagination, error) {
	var permissions []*model.Permission

	// 查询权限列表
	pagination, err := pagination.GetListByContext(mysql.DB, &permissions, c)
	if err != nil {
		return nil, nil, err
	}

	return permissions, pagination, nil
}

// GetPermissionByID 根据 ID 获取菜单信息
func (service *PermissionService) GetPermissionByID(id uint) (*model.Permission, error) {
	permission := &model.Permission{}
	if err := mysql.DB.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

// CreatePermission 创建菜单
func (service *PermissionService) CreatePermission(permissionCreatedRequest *model.PermissionCreatedRequest) error {
	// 状态验证是自动的，通过 UnmarshalJSON 实现
	if err := mysql.DB.Model(&model.Permission{}).Create(permissionCreatedRequest).Error; err != nil {
		return fmt.Errorf("创建权限失败: %w", err)
	}
	return nil
}

// DeletePermission 删除菜单
func (service *PermissionService) DeletePermission(id uint) error {
	if err := mysql.DB.Delete(&model.Permission{}, id).Error; err != nil {
		return fmt.Errorf("删除权限失败: %w", err)
	}

	return nil
}

// PatchPermission 修改菜单
func (service *PermissionService) PatchPermission(id uint, permissionPatchRequest *model.PermissionPatchRequest) error {
	// 直接更新并检查是否存在
	// 状态验证是自动的，通过 UnmarshalJSON 实现
	if err := mysql.DB.Model(&model.Permission{}).Where("id = ?", id).Updates(permissionPatchRequest).Error; err != nil {
		return fmt.Errorf("更新菜单失败: %w", err)
	}

	return nil
}
