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
func (service *RoleService) CreateRole(role *model.Role) error {
	// 检查 status 是否合规
	if *role.Status != model.RoleStatusActive && *role.Status != model.RoleStatusInactive {
		return fmt.Errorf("status 必须为 %d 或 %d", model.RoleStatusActive, model.RoleStatusInactive)
	}

	if err := mysql.DB.Create(role).Error; err != nil {
		return fmt.Errorf("创建角色失败: %w", err)
	}
	return nil
}

// UpdateRole 更新角色
func (service *RoleService) UpdateRole(id uint, role *model.Role) error {
	// 检查 status 是否合规
	if *role.Status != model.RoleStatusActive && *role.Status != model.RoleStatusInactive {
		return fmt.Errorf("status 必须为 %d 或 %d", model.RoleStatusActive, model.RoleStatusInactive)
	}

	// 使用 Select 指定要更新的字段，避免零值更新问题
	if err := mysql.DB.Model(&model.Role{}).Where("id = ?", id).Select("name", "code", "remark", "status").Updates(role).Error; err != nil {
		return fmt.Errorf("更新角色失败: %w", err)
	}
	return nil
}

// PatchRoleStatus 修改角色状态信息
func (service *RoleService) PatchRoleStatus(id uint, patchRoleRequest *model.PatchRoleRequest) error {
	// 检查 status 是否合规
	if *patchRoleRequest.Status != model.RoleStatusActive && *patchRoleRequest.Status != model.RoleStatusInactive {
		return fmt.Errorf("status 必须为 %d 或 %d", model.RoleStatusActive, model.RoleStatusInactive)
	}

	if err := mysql.DB.Model(&model.Role{}).Where("id = ?", id).Update("status", patchRoleRequest.Status).Error; err != nil {
		return fmt.Errorf("更新角色状态失败: %w", err)
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
