package response

import (
	"ffly-baisc/pkg/pagination"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type PageResponse struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

func Success(c *gin.Context, data any, message string) {
	if message == "" {
		message = "success"
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, httpCode int, message string, err error) {
	if err != nil {
		message = err.Error()
	}
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
		Data:    nil,
	})
}

func SuccessWithPagination(c *gin.Context, data any, p *pagination.Pagination, message string) {
	if message == "" {
		message = "success"
	}
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data: PageResponse{
			List:  data,
			Total: p.Total,
			Page:  p.Page,
			Size:  p.Size,
		},
	})
}
