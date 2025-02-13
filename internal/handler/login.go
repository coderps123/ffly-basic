package handler

import (
	"ffly-baisc/internal/service"
	"ffly-baisc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var login service.LoginService

	if err := c.ShouldBindJSON(&login); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	token, err := login.Login()
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "用户名或密码错误", err)
		return
	}

	response.Success(c, token, "登录成功")
}

func Register(c *gin.Context) {
	var register service.RegisterService

	if err := c.ShouldBindJSON(&register); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err := register.Register()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "注册失败", err)
		return
	}

	response.Success(c, nil, "注册成功")
}
