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

	// 创建用户
	if err := mysql.DB.Model(&model.User{}).Create(userCreateRequest).Error; err != nil {
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

// UpdateUser 修改用户
func (service *UserService) UpdateUser(id uint, userUpdateRequest *model.UserUpdateRequest) error {
	// 直接更新并检查是否存在
	result := mysql.DB.Model(&model.User{}).Where("id = ?", id).Updates(userUpdateRequest)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// PatchUserStatus 修改用户状态信息
func (service *UserService) PatchUserStatus(id uint, userPatchStatusRequest *model.UserPatchStatusRequest) error {
	// 直接更新并检查是否存在
	result := mysql.DB.Model(&model.User{}).Where("id = ?", id).Updates(userPatchStatusRequest)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
