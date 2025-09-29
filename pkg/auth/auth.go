package auth

import (
	"ffly-baisc/internal/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"` // "access" 或 "refresh"
	jwt.RegisteredClaims
}

// TokenPair Token 对
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // Access Token 过期时间（秒）
}

// GenerateTokenPair 生成 Access Token 和 Refresh Token
func GenerateTokenPair(userID uint, username string) (*TokenPair, error) {
	// Access Token - 短期有效
	accessToken, err := generateToken(userID, username, "access", 60*60*30) // 30分钟
	if err != nil {
		return nil, err
	}

	// Refresh Token -
	refreshToken, err := generateToken(userID, username, "refresh", 2*24*60*60) // 2天
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    30 * 60,
	}, nil
}

// generateToken 生成指定类型的 Token
func generateToken(userID uint, username string, tokenType string, expiresIn int64) (string, error) {
	claims := Claims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.App.JWTSecret))
}

// RefreshAccessToken 使用 Refresh Token 刷新 Access Token
func RefreshAccessToken(refreshTokenString string) (*TokenPair, error) {
	// 解析 Refresh Token
	claims, err := ParseToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("无效的 Refresh Token: %w", err)
	}

	// 验证是否为 Refresh Token
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("Token 类型错误，需要 Refresh Token")
	}

	// 生成新的 Token 对
	return GenerateTokenPair(claims.UserID, claims.Username)
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.App.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GetUserRoles 根据用户ID获取用户角色列表
func GetUserRoles(userID uint) ([]uint, error) {
	// 这里需要导入数据库相关的包，为了避免循环依赖，建议在service层实现
	// 或者创建一个专门的权限服务
	return nil, nil
}

// HasPermission 检查用户是否有指定权限
func HasPermission(userID uint, permission string) (bool, error) {
	// 实现权限检查逻辑
	// 1. 根据用户ID查询用户角色
	// 2. 根据角色查询权限
	// 3. 检查是否包含指定权限
	return false, nil
}
