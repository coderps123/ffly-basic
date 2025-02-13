package model

// System 系统信息模型
type System struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	StartTime string `json:"start_time"` // 启动时间
}

// TableName 表名
func (s *System) TableName() string {
	return "system"
}
