package utils

import (
	"fmt"
	"golang.org/x/exp/slog"
	"huawei-voicecall-webhook/config"
	"time"
)

// ParseConfigTime 解析配置文件中的 时分秒 时间字符串
func ParseConfigTime(timeString string) (time.Time, error) {
	// 解析时分秒
	t, err := time.Parse("15:04:05", timeString)
	if err != nil {
		return time.Time{}, err
	}

	// 获取当天的日期
	now := time.Now()
	year, month, day := now.Date()

	// 设置检查时间的日期为当天
	t = time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), 0, now.Location())

	return t, nil
}

// IsSendingAllowed 检查是否发送通知
func IsSendingAllowed() bool {
	// 从配置文件中获取起始和结束时间的字符串
	startTimeStr := config.Cfg.Common.ForbiddenPeriod.StartTime // 08:00:00
	endTimeStr := config.Cfg.Common.ForbiddenPeriod.EndTime     // 22:00:00

	// 解析每天的禁止发送消息起始时间
	startTime, err := ParseConfigTime(startTimeStr)
	if err != nil {
		slog.Error(fmt.Sprintf("起始时间解析错误: %v\n", err))
	}

	// 解析每天的禁止发送消息结束时间
	endTime, err := ParseConfigTime(endTimeStr)
	if err != nil {
		slog.Error(fmt.Sprintf("起始时间解析错误: %v\n", err))
	}

	// 获取当前时间
	currentTime := time.Now()

	// 检查当前时间是否在范围内
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		// 在这里可以执行请求的逻辑
		return false
	} else {
		// 在这里可以返回错误或信息给请求方
		return true
	}

}
