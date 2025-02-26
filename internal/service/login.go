package service

import (
	"errors"
	"ffly-baisc/internal/db"
	"ffly-baisc/internal/model"
	"ffly-baisc/pkg/auth"
	"ffly-baisc/pkg/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type LoginService struct {
	Username string `json:"username" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

// LoginLimiter 用户登录限流（一分钟内最多登录5次）
func LoginLimiter(loginService *LoginService) (bool, error) {
	// 构造 redis key
	key := fmt.Sprintf("login_attempts:%s", loginService.Username)

	// 使用 Redis 的 INCR 命令增加计数
	count, err := db.DB.Redis.Incr(key).Result()
	if err != nil {
		return false, err
	}

	// 如果是第一次登录，则设置过期时间为 1 分钟
	if count == 1 {
		db.DB.Redis.Expire(key, time.Minute)
	}

	// 检查是否超过限制
	if count > 5 {
		return true, fmt.Errorf("登录次数过多，请稍后再试")
	}

	return false, nil
}

func (service *LoginService) Login() (string, error) {
	// 用户登录限流
	isLimit, err := LoginLimiter(service)
	if err != nil {
		return "", err
	}
	if isLimit {
		return "", err
	}

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
	Username        *string `json:"username" binding:"required,min=2,max=20"`
	Password        *string `json:"password" binding:"required,min=6,max=20"`
	ConfirmPassword *string `json:"confirmPassword" binding:"required,min=6,max=20"`
	Nickname        *string `json:"nickname"`
	Email           *string `json:"email" binding:"omitempty,email"` // omitempty 允许为空
	Phone           *string `json:"phone" binding:"omitempty,e164"`  // omitempty 允许为空
}

func (service *RegisterService) Register() error {
	// 密码存在并且检查密码是否一致
	if *service.Password != "" && *service.Password != *service.ConfirmPassword {
		return errors.New("两次密码输入不一致")
	}

	// 创建用户
	userCreateRequest := &model.UserCreateRequest{
		Username: service.Username,
		Password: service.Password,
		Nickname: service.Nickname,
		Email:    service.Email,
		Phone:    service.Phone,
	}
	var userService UserService
	if err := userService.CreateUser(userCreateRequest); err != nil {
		return err
	}

	return nil
}
