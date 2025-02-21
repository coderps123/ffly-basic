package service

import (
	"errors"
	"ffly-baisc/internal/db"
	"ffly-baisc/internal/model"
	"ffly-baisc/pkg/auth"
	"ffly-baisc/pkg/utils"

	"gorm.io/gorm"
)

type LoginService struct {
	Username string `json:"username" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

func (service *LoginService) Login() (string, error) {
	// 检查用户名是否存在
	var user model.User
	if err := db.DB.MySQL.Where("username = ?", service.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// gorm.ErrRecordNotFound 是 gorm 的一个错误类型，表示没有找到记录
			// Is 用于判断错误是否为 gorm.ErrRecordNotFound
			return "", errors.New("用户名不存在")
		}
		return "", err
	}

	// 验证密码
	if !utils.CheckPassword(*user.Password, service.Password) {
		return "", errors.New("密码错误")
	}

	// 生成 token
	token, err := auth.GenerateToken(user.ID, *user.Username, user.RoleID)
	if err != nil {
		return "", err
	}

	return token, nil
}

type RegisterService struct {
	Username        string `json:"username" binding:"required,min=2,max=20"`
	Password        string `json:"password" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=6,max=20"`
	Nickname        string `json:"nickname"`
	Email           string `json:"email" binding:"omitempty,email"` // omitempty 允许为空
	Phone           string `json:"phone" binding:"omitempty,e164"`  // omitempty 允许为空
}

func (service *RegisterService) Register() error {
	// 检查用户名是否存在
	var count int64
	if err := db.DB.MySQL.Model(&model.User{}).Where("username = ?", service.Username).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("用户名已存在")
	}

	if service.Password != service.ConfirmPassword {
		return errors.New("两次密码输入不一致")
	}

	// 加密密码
	hashedPassword, err := utils.EncodePassword(service.Password)
	if err != nil {
		return err
	}

	// 创建用户
	user := &model.User{
		Username: &service.Username,
		Password: &hashedPassword,
		Nickname: &service.Nickname,
		Email:    &service.Email,
		Phone:    &service.Phone,
	}

	return db.DB.MySQL.Create(user).Error
}
