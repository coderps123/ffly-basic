package handler

import (
	"ffly-baisc/internal/service"
	"ffly-baisc/pkg/auth"
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

	response.Success(c, token, nil, "登录成功")
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

	response.Success(c, nil, nil, "注册成功")
}

func RefreshToken(c *gin.Context) {
	// 从请求头获取 Refresh Token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, http.StatusUnauthorized, "未提供 Refresh Token", nil)
		return
	}

	// 解析 Bearer Token
	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	} else {
		response.Error(c, http.StatusUnauthorized, "Token 格式错误", nil)
		return
	}

	// 刷新 Access Token
	tokenPair, err := auth.RefreshAccessToken(tokenString)
	if err != nil {
		c.Header("refresh_token_expired", "true")
		response.Error(c, http.StatusUnauthorized, "登录超时，请重新登录", err)
		return
	}

	response.Success(c, tokenPair, nil, "Token 刷新成功")
}
