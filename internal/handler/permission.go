package handler

import (
	"ffly-baisc/internal/model"
	"ffly-baisc/internal/service"
	"ffly-baisc/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPermissionList 获取权限列表
func GetPermissionList(c *gin.Context) {
	var permissionService service.PermissionService

	permissions, pagination, err := permissionService.GetPermissionList(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取权限列表失败", err)
		return
	}

	response.SuccessWithPagination(c, permissions, pagination, "菜单列表获取成功")
}

// GetPermission 获取菜单信息
func GetPermission(c *gin.Context) {
	var permissionService service.PermissionService

	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 解析用户ID 10：表示10进制，64：表示64位
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	permission, err := permissionService.GetPermissionByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取菜单信息失败", err)
		return
	}

	response.Success(c, permission, "获取成功")
}

// CreatePermission 创建菜单
func CreatePermission(c *gin.Context) {
	var permissionService service.PermissionService

	// 解析请求参数
	var permissionCreatedRequest model.PermissionCreatedRequest
	if err := c.ShouldBindJSON(&permissionCreatedRequest); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 创建菜单
	if err := permissionService.CreatePermission(&permissionCreatedRequest); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建菜单失败", err)
		return
	}

	response.Success(c, nil, "创建成功")
}

// PatchPermission 更新部分菜单信息
func PatchPermission(c *gin.Context) {
	var permissionService service.PermissionService
	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 解析用户ID 10：表示10进制，64：表示64位
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var permissionPatchRequest = model.PermissionPatchRequest{}
	if err := c.ShouldBindJSON(&permissionPatchRequest); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if err := permissionService.PatchPermission(uint(id), &permissionPatchRequest); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新菜单信息失败", err)
		return
	}

	response.Success(c, nil, "更新成功")
}

// DeletePermission 删除菜单
func DeletePermission(c *gin.Context) {
	var permissionService service.PermissionService

	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 解析用户ID 10：表示10进制，64：表示64位
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if err := permissionService.DeletePermission(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除菜单失败", err)
		return
	}

	response.Success(c, nil, "删除成功")
}
