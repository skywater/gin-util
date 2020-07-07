package util

import "strings"

// IsBlank 字符串是否为空
func IsBlank(str string) bool {
	return strings.TrimSpace(str) == ""
}

// IsNotBlank 字符串是否不为空
func IsNotBlank(str string) bool {
	return !IsBlank(str)
}
