package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"ffly-baisc/internal/model"
	"ffly-baisc/internal/service"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter               // 原始的响应写入器
	body               *bytes.Buffer // 用于保存响应体的缓冲区
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)                  //  将数据写入缓冲区
	return w.ResponseWriter.Write(b) //  将数据写入原始的 ResponseWriter
}

// ApiLog 记录API访问日志
func ApiLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求体
		var reqBodyBytes []byte
		if c.Request.Body != nil {
			// 1. 一次性读取全部内容
			reqBodyBytes, _ = io.ReadAll(c.Request.Body)
			// 2. 创建新的缓冲区， 将已读取的请求体数据保存到一个可重复读取的缓冲区
			newReader := bytes.NewBuffer(reqBodyBytes)
			// 3. 替换原来的请求体
			c.Request.Body = io.NopCloser(newReader)
		}

		responseBodyWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,        // 原始的 ResponseWriter
			body:           &bytes.Buffer{}, // 初始化缓冲区
		}
		c.Writer = responseBodyWriter // 替换为自定义的 ResponseWriter

		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算请求处理时间
		duration := time.Since(startTime)

		var typeStr string
		if strings.Contains(c.Request.URL.Path, "/login") {
			typeStr = "login"
		} else {
			typeStr = "operate"
		}

		// 创建日志记录
		apiLog := &model.ApiLog{
			UserID:       c.GetUint("userID"),
			Username:     c.GetString("username"),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Query:        c.Request.URL.RawQuery,
			Body:         string(reqBodyBytes),
			ResponseBody: responseBodyWriter.body.String(), // 获取响应体
			ClientIP:     c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			Type:         typeStr,
			StatusCode:   c.Writer.Status(),
			Duration:     duration.Milliseconds(), // 转换为毫秒
		}

		// 异步保存日志
		go func(log *model.ApiLog) {
			var service service.ApiLogService
			// 创建日志
			if err := service.CreateApiLog(log); err != nil {
				// 记录日志失败
			}
		}(apiLog)
	}
}
