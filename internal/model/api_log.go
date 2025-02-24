package model

type ApiLog struct {
	UserID       uint   `json:"userId"`
	Username     string `json:"username"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	Query        string `json:"query"`
	Body         string `json:"body"`
	UserAgent    string `json:"userAgent"`
	ClientIP     string `json:"clientIp"`
	StatusCode   int    `json:"statusCode"`
	Duration     int64  `json:"duration"`
	ResponseBody string `json:"responseBody"`
	Type         string `json:"type"` // operate: 操作日志, login: 登录日志
	BaseModel
}

func (log *ApiLog) TableName() string {
	return "api_logs"
}
