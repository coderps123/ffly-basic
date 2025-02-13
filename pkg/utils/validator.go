package utils

import "regexp"

func IsPhone(phone string) bool {
	// 使用正则表达式验证手机号码是否合法
	reg := `^1[34578]\d{9}$`
	if ok, _ := regexp.MatchString(reg, phone); !ok {
		return false
	}
	return true
}
