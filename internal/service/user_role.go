package service

import (
	"errors"
	"ffly-baisc/internal/model"
	types "ffly-baisc/pkg/type"
	"fmt"

	"gorm.io/gorm"
)

type UserRoleService struct{}

// SaveUserRole 创建或更新用户角色关联 roleID 为 0 表示删除用户角色关联
func (service *UserRoleService) SaveUserRole(tx *gorm.DB, userID uint, roleID uint) error {
	if roleID == 0 {
		// 如果 roleID 为 0，则表示要删除用户角色关联
		if err := tx.Where("user_id = ?", userID).First(&model.UserRole{}).Error; err != nil {
			// 先查询是否有已存在用户角色关联
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 记录不存在，则忽略
				return nil
			}
			tx.Rollback()
			return fmt.Errorf("查询用户角色关联失败: %v", err)
		}

		// 存在则删除
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("删除用户角色关联失败: %v", err)
		}

		return nil
	}

	var roleService RoleService
	// 查询角色
	role, err := roleService.GetRoleByID(roleID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 判定角色是否可用
	if role.Status == types.StatusDisabled {
		return fmt.Errorf("角色不可用")
	}

	// 查询是否有已存在用户角色关联
	if err := tx.Where("user_id = ?", userID).First(&model.UserRole{}).Error; err != nil {
		// 没有已存在用户角色关联，则创建
		var userRole = &model.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
		if err := tx.Model(model.UserRole{}).Create(userRole).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("创建用户角色关联失败: %v", err)
		}
	} else {
		// 有已存在用户角色关联，则更新
		if err := tx.Model(model.UserRole{}).Where("user_id = ?", userID).Update("role_id", roleID).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("更新用户角色关联失败: %v", err)
		}
	}

	return nil
}
