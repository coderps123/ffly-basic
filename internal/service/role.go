package service

import (
	"ffly-baisc/internal/model"
	"ffly-baisc/internal/mysql"
	"ffly-baisc/pkg/pagination"
	"fmt"

	"github.com/gin-gonic/gin"
)

type RoleService struct{}

// GetRoleList 获取角色列表
func (service *RoleService) GetRoleList(c *gin.Context) ([]*model.Role, *pagination.Pagination, error) {
	var roles []*model.Role

	// 查询权限列表
	pagination, err := pagination.GetListByContext(mysql.DB, &roles, c)
	if err != nil {
		return nil, nil, err
	}

	return roles, pagination, nil
}

// GetRoleByID 获取角色
func (service *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := mysql.DB.First(&role, id).Error; err != nil {
		return nil, fmt.Errorf("获取角色失败: %w", err)
	}
	return &role, nil
}

// CreateRole 创建角色
func (service *RoleService) CreateRole(roleCreateRequest *model.RoleCreateRequest) error {
	if err := mysql.DB.Model(&model.Role{}).Create(roleCreateRequest).Error; err != nil {
		return fmt.Errorf("创建角色失败: %w", err)
	}
	return nil
}

// PatchRole 部分更新角色
func (service *RoleService) PatchRole(id uint, rolePatchRequest *model.RolePatchRequest) error {
	if err := mysql.DB.Model(&model.Role{}).Where("id = ?", id).Updates(rolePatchRequest).Error; err != nil {
		return fmt.Errorf("更新角色失败: %w", err)
	}
	return nil
}

// DeleteRole 删除角色
func (service *RoleService) DeleteRole(id uint) error {
	if err := mysql.DB.Delete(&model.Role{}, id).Error; err != nil {
		return fmt.Errorf("删除角色失败: %w", err)
	}
	return nil
}
