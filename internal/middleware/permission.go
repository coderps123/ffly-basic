package middleware

import (
	"ffly-baisc/internal/service"
	"ffly-baisc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequirePermission 权限检查中间件
// 使用方式：router.Use(RequirePermission("user:read", "user:write"))
// func RequirePermission(permissions ...string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// 获取用户ID
// 		userID, exists := c.Get("userID")
// 		if !exists {
// 			response.Error(c, http.StatusUnauthorized, "用户未认证", nil)
// 			c.Abort()
// 			return
// 		}

// 		// 检查权限
// 		var authService service.AuthPermissionService
// 		hasPermission, err := authService.HasAnyPermission(userID.(uint), permissions)
// 		if err != nil {
// 			response.Error(c, http.StatusInternalServerError, "权限检查失败", err)
// 			c.Abort()
// 			return
// 		}

// 		if !hasPermission {
// 			response.Error(c, http.StatusForbidden, "权限不足", nil)
// 			c.Abort()
// 			return
// 		}

// 		c.Next()
// 	}
// }

// RequireRole 角色检查中间件
// 使用方式：router.Use(RequireRole(1, 2, 3))
func RequireRole(roleIDs ...uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("userID")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "用户未认证", nil)
			c.Abort()
			return
		}

		// 获取用户角色
		var authService service.AuthPermissionService
		userRoles, err := authService.GetUserRoles(userID.(uint))
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "角色查询失败", err)
			c.Abort()
			return
		}

		// 检查是否有任意一个角色
		hasRole := false
		for _, userRole := range userRoles {
			for _, requiredRole := range roleIDs {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			response.Error(c, http.StatusForbidden, "角色权限不足", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
