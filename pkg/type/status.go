package types

import (
	"encoding/json"
	"fmt"
)

type Status int

const (
	StatusEnabled  Status = iota + 1 // 启用
	StatusDisabled                   // 禁用
)

var statusNames = map[Status]string{
	StatusEnabled:  "启用",
	StatusDisabled: "禁用",
}

// String 获取状态名称
func (s Status) String() string {
	name, ok := statusNames[s]
	if ok {
		return name
	}

	return "未知状态"
}

// Valid 验证状态是否有效
func (s Status) Valid() bool {
	_, ok := statusNames[s]
	return ok
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
// 用于将 json 反序列化为 Status 类型
func (s *Status) UnmarshalJSON(data []byte) error {
	var value int32
	// 将 json 解析到 int32 变量中
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	// 将 int32 变量转换为 Status 类型
	*s = Status(value)
	// 验证状态是否有效
	if !s.Valid() {
		return fmt.Errorf("invalid status value: %d", value)
	}
	return nil
}
