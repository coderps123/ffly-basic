package service

import (
	"ffly-baisc/internal/db"
	"ffly-baisc/internal/model"
	"ffly-baisc/pkg/pagination"
	"fmt"

	"github.com/gin-gonic/gin"
)

type PermissionService struct{}

// BuildPermissionTree 构建权限树
func (service *PermissionService) BuildPermissionTree(permissions []*model.Permission, id uint) []*model.Permission {
	var trees []*model.Permission
	for _, permission := range permissions {
		// 父级权限ID等于当前权限ID，则为子权限
		if permission.ParentID == id {
			// 递归构建子权限树
			children := service.BuildPermissionTree(permissions, permission.ID)
			if len(children) > 0 {
				permission.Children = children
			}
			// 加入树中
			trees = append(trees, permission)
		}
	}

	return trees
}

// GetPermissionList 获取权限列表
func (service *PermissionService) GetPermissionList(c *gin.Context) ([]*model.Permission, *pagination.Pagination, error) {
	var permissions []*model.Permission

	// 查询权限列表
	if err := db.DB.MySQL.Find(&permissions).Error; err != nil {
		return nil, nil, fmt.Errorf("获取权限列表失败: %w", err)
	}

	// 获取分页信息
	pageination := pagination.GetPageInfo(c)

	// 获取权限树
	permissionTree := service.BuildPermissionTree(permissions, 0)

	total := int64(len(permissionTree))

	start := (pageination.Page - 1) * pageination.Size
	end := start + pageination.Size
	// 判断是否越界
	if end > len(permissionTree) {
		end = len(permissionTree)
	}
	if start > len(permissionTree) {
		start = len(permissionTree)
	}

	return permissionTree[start:end], &pagination.Pagination{
		Page:  pageination.Page,
		Size:  pageination.Size,
		Total: total,
	}, nil
}

// GetPermissionByID 根据 ID 获取菜单信息
func (service *PermissionService) GetPermissionByID(id uint) (*model.Permission, error) {
	permission := &model.Permission{}
	if err := db.DB.MySQL.First(&permission, id).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

// CreatePermission 创建菜单
func (service *PermissionService) CreatePermission(permissionCreatedRequest *model.PermissionCreatedRequest) error {
	// 状态验证是自动的，通过 UnmarshalJSON 实现
	if err := db.DB.MySQL.Model(&model.Permission{}).Create(permissionCreatedRequest).Error; err != nil {
		return fmt.Errorf("创建权限失败: %w", err)
	}
	return nil
}

// DeletePermission 删除菜单及其所有子菜单
func (service *PermissionService) DeletePermission(id uint) error {
	var permissions []*model.Permission

	// 查询所有的菜单
	if err := db.DB.MySQL.Find(&permissions).Error; err != nil {
		return fmt.Errorf("获取权限列表失败: %w", err)
	}

	// 获取待删除的权限ID列表
	permissionIDs := service.getPermissionIDsToDelete(permissions, id)

	// 删除所有的相关权限
	if err := db.DB.MySQL.Where("id IN ?", permissionIDs).Delete(&model.Permission{}).Error; err != nil {
		return fmt.Errorf("删除权限失败: %w", err)
	}

	return nil
}

// getPermissionIDsToDelete 获取待删除的权限ID列表
func (service *PermissionService) getPermissionIDsToDelete(permissions []*model.Permission, id uint) []uint {
	var ids []uint
	// 先添加待删除的权限ID
	ids = append(ids, id)
	for _, permission := range permissions {
		// 如果当前权限的父级ID等于待删除的ID，则添加到列表中
		if permission.ParentID == id {
			ids = append(ids, permission.ID)
			// 递归获取子节点的ID
			childIDs := service.getPermissionIDsToDelete(permissions, permission.ID)
			ids = append(ids, childIDs...)
		}
	}
	return ids
}

// PatchPermission 修改菜单
func (service *PermissionService) PatchPermission(id uint, permissionPatchRequest *model.PermissionPatchRequest) error {
	// 直接更新并检查是否存在
	// 状态验证是自动的，通过 UnmarshalJSON 实现
	if err := db.DB.MySQL.Model(&model.Permission{}).Where("id = ?", id).Updates(permissionPatchRequest).Error; err != nil {
		return fmt.Errorf("更新菜单失败: %w", err)
	}

	return nil
}
