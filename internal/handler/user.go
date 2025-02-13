package handler

import (
	"ffly-baisc/internal/model"
	"ffly-baisc/internal/service"
	"ffly-baisc/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	var userService service.UserService

	users, pagination, err := userService.GetUserList(c)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户列表失败", err)
		return
	}

	response.SuccessWithPagination(c, users, pagination, "用户列表获取成功")
}

// GetUser 获取用户信息
func GetUser(c *gin.Context) {
	var userService service.UserService

	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 解析用户ID 10：表示10进制，64：表示64位
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	user, err := userService.GetUserByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户信息失败", err)
		return
	}

	response.Success(c, user, "获取成功")
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var userService service.UserService

	// 解析请求参数
	var userCreateRequest model.UserCreateRequest
	if err := c.ShouldBindJSON(&userCreateRequest); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 创建用户
	if err := userService.CreateUser(&userCreateRequest); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建用户失败", err)
		return
	}

	response.Success(c, nil, "创建成功")
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	var userService service.UserService

	// 解析请求参数
	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 解析用户ID 10：表示10进制，64：表示64位
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var userUpdateRequest model.UserUpdateRequest
	if err := c.ShouldBindJSON(&userUpdateRequest); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 更新用户信息
	if err := userService.UpdateUser(uint(id), &userUpdateRequest); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新用户信息失败", err)
		return
	}

	response.Success(c, nil, "更新成功")
}

// PatchUser 更新部分用户信息
func PatchUser(c *gin.Context) {
	var userService service.UserService
	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 解析用户ID 10：表示10进制，64：表示64位
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	var UserPatchStatusRequest = model.UserPatchStatusRequest{}
	if err := c.ShouldBindJSON(&UserPatchStatusRequest); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if err := userService.PatchUserStatus(uint(id), &UserPatchStatusRequest); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新用户信息失败", err)
		return
	}

	response.Success(c, nil, "更新成功")
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	var userService service.UserService

	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 解析用户ID 10：表示10进制，64：表示64位
	if err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	if err := userService.DeleteUser(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除用户失败", err)
		return
	}

	response.Success(c, nil, "删除成功")
}

// GetCurrentUserInfo 获取当前用户信息
func GetCurrentUserInfo(c *gin.Context) {
	var userService service.UserService
	userID := c.GetUint("userID")

	user, err := userService.GetUserByID(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户信息失败", err)
		return
	}

	response.Success(c, user, "获取成功")
}
