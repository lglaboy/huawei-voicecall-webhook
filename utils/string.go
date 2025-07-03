package utils

import (
	"strings"
	"time"
	"unicode/utf8"
)

// CheckStringUnicodeLength Unicode编码字符长度检查
func CheckStringUnicodeLength(s string) int {
	// 使用utf8.RuneCountInString检查字符长度
	return utf8.RuneCountInString(s)
}

// TruncateStringByUnicodeLength 将字符串裁剪成指定 Unicode编码长度
func TruncateStringByUnicodeLength(s string, length int) string {
	if length <= 0 {
		return ""
	}

	// 将字符串转换为rune数组
	runes := []rune(s)

	// 确保不截取超过字符串长度的字符
	if length > len(runes) {
		length = len(runes)
	}

	// 截取前length个字符
	truncatedRunes := runes[:length]

	// 将截取后的字符数组转换回字符串
	return string(truncatedRunes)
}

// TruncateStringByByteLength 将字符串裁剪成指定 UTF-8 编码长度
func TruncateStringByByteLength(s string, length int) string {
	if length <= 0 {
		return ""
	}

	runes := []rune(s)
	result := make([]rune, 0)
	currentByteLength := 0

	for _, r := range runes {
		//runeSize := utf8.RuneLen(r)
		runeByteLength := utf8.RuneLen(r)
		if currentByteLength+runeByteLength > length {
			break
		}
		result = append(result, r)
		currentByteLength += runeByteLength
	}

	return string(result)
}

// AddDefaultCountryCode 检查电话号码是否包含国家码，如果不包含，则添加默认的国家码。
func AddDefaultCountryCode(phoneNumber string, defaultCountryCode string) string {
	// 假设国家码在电话号码中使用加号 "+" 表示
	if !strings.HasPrefix(phoneNumber, "+") {
		// 如果电话号码不以 "+" 开头，添加默认的国家码
		return defaultCountryCode + phoneNumber
	}
	// 否则，电话号码已包含国家码，不做修改
	return phoneNumber
}

func ParseDuration(s string) (time.Duration, error) {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	return duration, nil
}
