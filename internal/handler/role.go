package handler

import (
	"ffly-baisc/internal/model"
	"ffly-baisc/internal/service"
	"ffly-baisc/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRoleList 获取角色列表
func GetRoleList(c *gin.Context) {
	var roleService service.RoleService

	roles, pagination, err := roleService.GetRoleList(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取角色列表失败", err)
		return
	}

	response.SuccessWithPagination(c, roles, pagination, "角色列表获取成功")
}

// GetRole 获取角色详情
func GetRole(c *gin.Context) {
	var roleService service.RoleService

	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 解析用户ID 10：表示10进制，64：表示64位
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID", err)
		return
	}

	roleInfo, err := roleService.GetRoleByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取角色失败", err)
		return
	}

	response.Success(c, roleInfo, "角色获取成功")
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	var roleService service.RoleService
	var roleCreateRequest model.RoleCreateRequest

	if err := c.ShouldBindJSON(&roleCreateRequest); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效", err)
		return
	}

	if err := roleService.CreateRole(&roleCreateRequest); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建角色失败", err)
		return
	}

	response.Success(c, nil, "角色创建成功")
}

// PatchRole 部分更新角色
func PatchRole(c *gin.Context) {
	var roleService service.RoleService
	var rolePatchRequest model.RolePatchRequest

	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 解析用户ID 10：表示10进制，64：表示64位
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID", err)
		return
	}

	if err := c.ShouldBindJSON(&rolePatchRequest); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数无效", err)
		return
	}

	if err := roleService.PatchRole(uint(id), &rolePatchRequest); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新角色失败", err)
		return
	}

	response.Success(c, nil, "角色更新成功")
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	var roleService service.RoleService

	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 解析用户ID 10：表示10进制，64：表示64位
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID", err)
		return
	}

	if err := roleService.DeleteRole(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除角色失败", err)
		return
	}

	response.Success(c, nil, "角色删除成功")
}
