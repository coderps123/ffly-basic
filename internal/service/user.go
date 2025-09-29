package service

import (
	"errors"
	"ffly-baisc/internal/db"
	"ffly-baisc/internal/model"
	"ffly-baisc/pkg/query"
	"ffly-baisc/pkg/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct{}

// GetUserList 获取用户列表
func (service *UserService) GetUserList(c *gin.Context) ([]*model.User, *query.Pagination, error) {
	// 开启事务
	tx := db.DB.MySQL.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 回滚事务
		}
	}()

	users, pagination, err := query.GetQueryData[model.User](db.DB.MySQL, c)
	if err != nil {
		tx.Rollback() // 回滚事务
		return nil, nil, err
	}

	// 为每个用户填充角色信息
	var userRoleService UserRoleService
	for _, user := range *users {
		userRoles, err := userRoleService.GetRolesByUserID(db.DB.MySQL, user.ID)
		if err != nil {
			tx.Rollback() // 回滚事务
			return nil, nil, fmt.Errorf("获取用户角色失败: %w", err)
		}
		// 获取 roleIds
		var roleIds []uint
		for _, userRole := range userRoles {
			roleIds = append(roleIds, userRole.RoleID)
		}
		// 获取角色信息
		var roles []*model.Role
		if err := db.DB.MySQL.Model(&model.Role{}).Where("id in (?)", roleIds).Find(&roles).Error; err != nil {
			tx.Rollback() // 回滚事务
			return nil, nil, fmt.Errorf("获取角色信息失败: %w", err)
		}
		// 填充角色信息
		user.Roles = roles
	}

	return *users, pagination, nil
}

// GetUserByID 根据 ID 获取用户信息
func (service *UserService) GetUserByID(id uint) (*model.User, error) {
	// 开启事务
	tx := db.DB.MySQL.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 回滚事务
		}
	}()

	user := &model.User{}
	if err := db.DB.MySQL.First(user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback() // 回滚事务
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	// 用户填充角色信息
	var userRoleService UserRoleService
	userRoles, err := userRoleService.GetRolesByUserID(db.DB.MySQL, user.ID)
	if err != nil {
		tx.Rollback() // 回滚事务
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}
	// 获取 roleIds
	var roleIds []uint
	for _, userRole := range userRoles {
		roleIds = append(roleIds, userRole.RoleID)
	}
	// 获取角色信息
	var roles []*model.Role
	if err := db.DB.MySQL.Model(&model.Role{}).Where("id in (?)", roleIds).Find(&roles).Error; err != nil {
		tx.Rollback() // 回滚事务
		return nil, fmt.Errorf("获取角色信息失败: %w", err)
	}
	// 填充角色信息
	user.Roles = roles

	return user, nil
}

// CreateUser 创建用户
func (service *UserService) CreateUser(userCreateRequest *model.UserCreateRequest) error {
	// 开启事务
	tx := db.DB.MySQL.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 回滚事务
		}
	}()

	// 校验手机号是否合规
	if userCreateRequest.Phone != nil && !utils.IsPhone(*userCreateRequest.Phone) {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("手机号不合规")
	}

	// 加密密码
	if userCreateRequest.Password == nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("密码不能为空")
	}
	hashedPassword, err := utils.EncodePassword(*userCreateRequest.Password)
	if err != nil {
		tx.Rollback() // 回滚事务
		return err
	}
	userCreateRequest.Password = &hashedPassword

	// 将请求数据转换为User模型
	user := &model.User{
		Username:  userCreateRequest.Username,
		Password:  userCreateRequest.Password,
		Nickname:  userCreateRequest.Nickname,
		Email:     userCreateRequest.Email,
		Phone:     userCreateRequest.Phone,
		Status:    userCreateRequest.Status,
		BaseModel: userCreateRequest.BaseModel,
	}

	// 创建用户
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}

	// 设置创建后的ID
	userCreateRequest.ID = user.ID

	// 创建用户角色关联（支持多个角色）
	if len(userCreateRequest.RoleIDs) > 0 {
		var userRoleService UserRoleService
		if err := userRoleService.SaveUserRoles(tx, user.ID, userCreateRequest.RoleIDs); err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		// 事务提交失败，回滚事务
		tx.Rollback() // 回滚事务
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// DeleteUser 删除用户
func (service *UserService) DeleteUser(id uint) error {
	// 开启事务
	tx := db.DB.MySQL.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // 回滚事务
		}
	}()

	// 删除用户
	if err := tx.Delete(&model.User{}, id).Error; err != nil {
		tx.Rollback() // 回滚事务
		return err
	}

	// 删除用户角色关联
	var userRoleService UserRoleService
	if err := userRoleService.SaveUserRole(tx, id, 0); err != nil {
		tx.Rollback() // 回滚事务
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// PatchUser 修改用户状态信息
func (service *UserService) PatchUser(id uint, userPatchRequest *model.UserPatchRequest) error {
	// 开启事务
	tx := db.DB.MySQL.Begin()
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

	// 更新用户角色关联，
	if len(userPatchRequest.RoleIDs) > 0 {
		// 多个的话，
		var userRoleService UserRoleService
		if err := userRoleService.SaveUserRoles(tx, id, userPatchRequest.RoleIDs); err != nil {
			tx.Rollback() // 回滚事务
			return err
		}
	}

	// 更新用户信息
	result := tx.Model(&model.User{}).Where("id = ?", id).Updates(userPatchRequest)
	if result.Error != nil {
		tx.Rollback() // 回滚事务
		return result.Error
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // 回滚事务
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
	if err := db.DB.MySQL.Model(&model.User{}).Where("id = ?", id).Update("password", &hashedPassword).Error; err != nil {
		return err
	}

	return nil
}
