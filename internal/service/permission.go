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

// // GetPermissionByID 根据 ID 获取菜单信息
// func (service *PermissionService) GetPermissionByID(id uint) (*model.Permission, error) {
// 	permission := &model.Permission{}
// 	if err := mysql.DB.First(&permission, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return permission, nil
// }

// CreatePermission 创建菜单
func (service *PermissionService) CreatePermission(permission *model.Permission) error {
	// 状态验证是自动的，通过 UnmarshalJSON 实现
	if err := mysql.DB.Create(permission).Error; err != nil {
		return fmt.Errorf("创建权限失败: %w", err)
	}
	return nil
}

// // DeletePermission 删除菜单
// func (service *PermissionService) DeletePermission(id uint) error {
// 	if err := mysql.DB.Delete(&model.Permission{}, id).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// UpdatePermission 修改菜单
func (service *PermissionService) UpdatePermission(id uint, permission *model.Permission) error {
	// 直接更新并检查是否存在
	// 状态验证是自动的，通过 UnmarshalJSON 实现
	if err := mysql.DB.Where("id = ?", id).Updates(permission).Error; err != nil {
		return fmt.Errorf("更新菜单失败: %w", err)
	}

	return nil
}

// // PatchPermissionStatus 修改菜单状态信息
// func (service *PermissionService) PatchPermissionStatus(id uint, patchPermissionStatus *model.PatchPermissionStatusRequest) error {
// 	// 检查 status 是否合规
// 	if *patchPermissionStatus.Status != model.PermissionStatusActive && *patchPermissionStatus.Status != model.PermissionStatusInactive {
// 		return fmt.Errorf("status 必须为 %d 或 %d", model.PermissionStatusActive, model.PermissionStatusInactive)
// 	}

// 	// 直接更新并检查是否存在
// 	result := mysql.DB.Model(&model.Permission{}).Where("id = ?", id).Updates(patchPermissionStatus)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }
