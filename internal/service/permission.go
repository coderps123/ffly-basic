package service

import (
	"ffly-baisc/internal/db"
	"ffly-baisc/internal/model"
	"ffly-baisc/pkg/file"
	"ffly-baisc/pkg/query"
	"fmt"

	"github.com/gin-gonic/gin"
)

type PermissionService struct{}

// BuildPermissionTree 构建权限树
func (service *PermissionService) BuildPermissionTree(permissions []*model.Permission, parentID uint) []*model.Permission {
	var trees []*model.Permission
	for _, permission := range permissions {
		// 父级权限ID等于当前权限ID，则为子权限
		if permission.ParentID == parentID {
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
func (service *PermissionService) GetPermissionList(c *gin.Context) ([]*model.Permission, *query.Pagination, error) {
	var permissions []*model.Permission

	// 查询权限列表
	if err := db.DB.MySQL.Find(&permissions).Error; err != nil {
		return nil, nil, fmt.Errorf("获取权限列表失败: %w", err)
	}

	// 获取分页信息
	pageination := query.GetPageInfo(c)

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

	return permissionTree[start:end], &query.Pagination{
		Page:  pageination.Page,
		Size:  pageination.Size,
		Total: &total,
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

// DeletePermission 删除菜单
func (service *PermissionService) DeletePermission(id uint) error {
	if err := db.DB.MySQL.Delete(&model.Permission{}, id).Error; err != nil {
		return fmt.Errorf("删除权限失败: %w", err)
	}

	return nil
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

// ExportPermission 导出菜单
func (service *PermissionService) ExportPermission(c *gin.Context) error {
	// 获取所有的权限
	var permissions []*model.Permission
	if err := db.DB.MySQL.Find(&permissions).Error; err != nil {
		return fmt.Errorf("获取权限列表失败: %w", err)
	}

	// 构建权限树
	permissionTree := service.BuildPermissionTree(permissions, 0)

	columns := []file.ColumnConfig{
		{Title: "ID", Field: "ID", Width: 20},
		{Title: "权限名称", Field: "Name", Width: 20, Prefix: "-->"},
		{Title: "权限类型", Field: "Type", Width: 20},
		{Title: "路由路径", Field: "Path", Width: 20},
		{Title: "权限码", Field: "Code", Width: 20},
		{Title: "组件名称", Field: "Component", Width: 20},
		{Title: "图标", Field: "Icon", Width: 20},
		{Title: "排序", Field: "Sort", Width: 20},
		{Title: "状态", Field: "Status", Width: 20},
		{Title: "父级ID", Field: "ParentID", Width: 20},
		{Title: "备注", Field: "Remark", Width: 20},
	}

	err := file.ExportExcel(c, permissionTree, columns, "权限列表", file.Options{})
	if err != nil {
		return fmt.Errorf("生成excel并返回字节流失败: %w", err)
	}

	return nil
}
