package service

import (
	"errors"
	"ffly-baisc/internal/model"
	"ffly-baisc/internal/mysql"
	"ffly-baisc/pkg/pagination"
	"ffly-baisc/pkg/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct{}

// GetUserList 获取用户列表
func (service *UserService) GetUserList(c *gin.Context) ([]*model.User, *pagination.Pagination, error) {
	var users []*model.User

	// 查询权限列表
	pagination, err := pagination.GetListByContext(mysql.DB, &users, c)
	if err != nil {
		return nil, nil, err
	}

	return users, pagination, nil
}

// GetUserByID 根据 ID 获取用户信息
func (service *UserService) GetUserByID(id uint) (*model.User, error) {
	user := &model.User{}
	if result := mysql.DB.First(user, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, result.Error
	}
	return user, nil
}

// CreateUser 创建用户
func (service *UserService) CreateUser(userCreateRequest *model.UserCreateRequest) error {
	// 校验手机号是否合规
	if userCreateRequest.Phone != nil && !utils.IsPhone(*userCreateRequest.Phone) {
		return fmt.Errorf("手机号不合规")
	}

	// 加密密码
	if userCreateRequest.Password == nil {
		return fmt.Errorf("密码不能为空")
	}
	hashedPassword, err := utils.EncodePassword(*userCreateRequest.Password)
	if err != nil {
		return err
	}
	userCreateRequest.Password = &hashedPassword

	var user = &model.User{
		Username: userCreateRequest.Username,
		Password: userCreateRequest.Password,
		Nickname: userCreateRequest.Nickname,
		Email:    userCreateRequest.Email,
		Phone:    userCreateRequest.Phone,
		Status:   userCreateRequest.Status,
		RoleID:   userCreateRequest.RoleID,
	}

	// 创建用户
	if err := mysql.DB.Model(&model.User{}).Create(user).Error; err != nil {
		return err
	}

	// 创建用户角色关联
	var userRoleService UserRoleService
	if err := userRoleService.CreateUserRole(user); err != nil {
		return err
	}

	return nil
}

// DeleteUser 删除用户
func (service *UserService) DeleteUser(id uint) error {
	if err := mysql.DB.Delete(&model.User{}, id).Error; err != nil {
		return err
	}

	return nil
}

// PatchUser 修改用户状态信息
func (service *UserService) PatchUser(id uint, userPatchRequest *model.UserPatchRequest) error {
	// 开启事务
	tx := mysql.DB.Begin()
	defer func() {
		if r := recover(); r != nil { // 遇到异常回滚事务
			tx.Rollback() // 回滚事务
		}
	}()

	// 校验手机号是否合规
	if userPatchRequest.Phone != nil && !utils.IsPhone(*userPatchRequest.Phone) {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("手机号不合规")
	}

	// 更新用户信息
	result := tx.Model(&model.User{}).Where("id = ?", id).Updates(userPatchRequest)
	if result.Error != nil {
		tx.Rollback() // 回滚事务
		return result.Error
	}

	// 如果用户传入 roleID
	if userPatchRequest.RoleID != nil {
		// 更新用户角色关联，
		var userRoleService UserRoleService
		if err := userRoleService.UpdateUserRole(tx, id, *userPatchRequest.RoleID); err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// UpdatePassword 修改用户密码
func (service *UserService) UpdatePassword(id uint, updatePasswordRequest *model.UpdatePasswordRequest) error {
	// 查询用户
	user, err := service.GetUserByID(id)
	if err != nil {
		return err
	}
	// 校验密码与确认密码是否一致
	if updatePasswordRequest.NewPassword != nil && updatePasswordRequest.PasswordConfirm != nil &&
		*updatePasswordRequest.NewPassword != *updatePasswordRequest.PasswordConfirm {
		return fmt.Errorf("新密码和确认密码不匹配")
	}

	// 校验密码是否正确
	if !utils.CheckPassword(*user.Password, *updatePasswordRequest.Password) {
		return fmt.Errorf("旧密码错误")
	}

	// 加密密码
	hashedPassword, err := utils.EncodePassword(*updatePasswordRequest.NewPassword)
	if err != nil {
		return err
	}

	// 更新密码
	if err := mysql.DB.Model(&model.User{}).Where("id = ?", id).Update("password", &hashedPassword).Error; err != nil {
		return err
	}

	return nil
}
