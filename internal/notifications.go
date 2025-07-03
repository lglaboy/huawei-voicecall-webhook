package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"huawei-voicecall-webhook/config"
	"huawei-voicecall-webhook/datastore"
	"huawei-voicecall-webhook/utils"
	"strings"
	"time"
)

// ProcessVoiceCallVariables 处理用于语音通知模板中的变量
func ProcessVoiceCallVariables(s string, i int) string {
	// 如果为空，返回字符串 NULL
	// 如果不为空，按照指定字符串长度检查
	if utils.CheckStringUnicodeLength(s) == 0 {
		return "NULL"
	}
	return utils.TruncateStringByByteLength(s, i)
}

func envNameAndAlertNameFromRuleName(s string) (string, string) {
	//空格拆分字符串，取第二个
	// 使用空格分割字符串
	parts := strings.Split(s, " ")

	// 检查是否有足够的部分
	if len(parts) >= 2 {
		return parts[0], parts[1]
	} else {
		slog.Error("字符串格式不符合预期")
		return "", ""
	}
}

// SendFirstVoiceNotify 首次发送语音通知
func SendFirstVoiceNotify(phone string, body []string) {
	phone = utils.AddDefaultCountryCode(phone, "+86")

	ttl, err := utils.ParseDuration(config.Cfg.Datastore.TTL)

	if err != nil {
		slog.Error(fmt.Sprintf("解析错误：%v\n", err))
	}

	// 不存在发送记录才发送，避免该号码位于多个环境，多次发送
	if v, _, b := datastore.GetData(phone); b == false {
		slog.Info(fmt.Sprintf("首次发送告警信息到: %s", phone))
		SendVoiceNotify(phone, body)
		datastore.SetData(phone, v+1, body, ttl)
	} else {
		slog.Info(fmt.Sprintf("该号码正在呼叫中: %s", phone))
	}
}

// SendRepeatVoiceNotify 再次发送语音通知
func SendRepeatVoiceNotify(phone string) {
	phone = utils.AddDefaultCountryCode(phone, "+86")

	ttl, err := utils.ParseDuration(config.Cfg.Datastore.TTL)
	if err != nil {
		slog.Error(fmt.Sprintf("解析错误：%v\n", err))
	}

	// 存在发送记录才发送，并检查是否达到发送次数限制
	if v, body, t := datastore.GetData(phone); t == true {
		if v < 3 {
			slog.Info(fmt.Sprintf("再次发送告警信息到: %s,内容: %s", phone, body))
			SendVoiceNotify(phone, body)
			datastore.SetData(phone, v+1, body, ttl)
		} else {
			// 本次未接通重复三次请求结束
			slog.Info(fmt.Sprintf("该号码: %s 已连续拨打三次，停止重复呼叫。", phone))
			datastore.DeleteData(phone)
		}
	} else {
		slog.Info(fmt.Sprintf("数据已过期/不存在,无法再次发送语音通知."))
	}
}

// SendVoiceNotificationToDuty 向值班人员发送语音通知
func SendVoiceNotificationToDuty(body []string) {
	duty, _ := utils.GetDutyFromCSV(config.Cfg.Schedule.File)
	// 获取当前时间
	currentTime := time.Now()

	// 格式化为 HH:MM:SS
	formattedTime := currentTime.Format("2006-01-02")

	for _, item := range duty {
		if item.Date == formattedTime {
			slog.Info(fmt.Sprintf("发送告警信息到值班人员: %s, 日期: %s", item.Name, item.Date))
			SendFirstVoiceNotify(item.Phone, body)
		}
	}
}

func SendVoiceNotificationToEnvDirector(body []string) {
	if len(config.Cfg.EnvDirector) >= 1 {
		for _, item := range config.Cfg.EnvDirector {
			if strings.HasPrefix(item.Env, body[0]) {
				slog.Info(fmt.Sprintf("发送告警信息到项目负责人: %s", item.Name))
				SendFirstVoiceNotify(item.Phone, body)
			}
		}
	}
}

// SendVoiceNotificationToLeader 向Leader 发送语音通知
func SendVoiceNotificationToLeader(body []string) {
	if len(config.Cfg.Lead) >= 1 {
		for _, lead := range config.Cfg.Lead {
			slog.Info(fmt.Sprintf("发送告警信息到Leader: %s", lead.Name))
			SendFirstVoiceNotify(lead.Phone, body)
		}
	}
}

func SendVoiceNotificationToAll(body []string) {
	SendVoiceNotificationToDuty(body)
	SendVoiceNotificationToEnvDirector(body)
	SendVoiceNotificationToLeader(body)
}

func SendVoiceNotification(envName, alertName, time string) error {
	// 创建一个字符串切片
	strList := []string{envName, alertName, time}

	// 使用encoding/json包将切片转换为JSON格式的字符串
	jsonStr, err := json.Marshal(strList)
	if err != nil {
		slog.Error(fmt.Sprintln("JSON marshal error:", err))
		return err
	}
	// 输出JSON格式的字符串
	slog.Info(string(jsonStr))

	// 检查当前时间是否允许发送
	if utils.IsSendingAllowed() {
		SendVoiceNotificationToAll(strList)
	} else {
		slog.Info("当前时间不在允许发送范围内.")
		datastore.AddRecord(strList)
	}
	return nil
}

func checkRecordData(envName, alertName, time string) error {
	// 处理数据
	envName = ProcessVoiceCallVariables(envName, 10)
	alertName = ProcessVoiceCallVariables(alertName, 12)
	// 创建一个字符串切片
	strList := []string{envName, alertName, time}

	// 使用encoding/json包将切片转换为JSON格式的字符串
	jsonStr, err := json.Marshal(strList)
	if err != nil {
		slog.Error(fmt.Sprintln("JSON marshal error:", err))
		return err
	}
	// 输出JSON格式的字符串
	slog.Info(string(jsonStr))

	// 删除recordData
	datastore.RemoveRecord(strList)

	return nil
}

// isInBlackList 检查envName是否存在于black_list中
func isInBlackList(envName string) bool {
	for _, item := range config.Cfg.BlackList {
		if strings.HasPrefix(item, envName) {
			return true
		}
	}
	return false
}

// isInWhiteList 检查envName是否存在于white_list中
func isInWhiteList(envName string) bool {
	for _, item := range config.Cfg.WhiteList {
		if strings.HasPrefix(item, envName) {
			return true
		}
	}
	return false
}

func CheckAndRunByMode(envName, alertName, time string) error {
	//# 发送告警
	//# 报警通知, ${TXT_10}环境的${TXT_12}服务于${TIME}出现紧急异常.
	//# 参考：https://support.huaweicloud.com/VoiceCall_faq/VoiceCall_faq_00020.html
	//# 模板采用UTF-8编码格式，汉字和中文符号为3个字节，字母、数字和英文符号为1个字节。

	slog.Info("提取内容",
		slog.String("envName", envName),
		slog.String("alertName", alertName))

	envName = ProcessVoiceCallVariables(envName, 10)
	alertName = ProcessVoiceCallVariables(alertName, 12)
	slog.Info("最终内容",
		slog.String("envName", envName),
		slog.String("alertName", alertName),
	)

	switch config.Cfg.Common.Model {
	case "blacklist":
		if !isInBlackList(envName) {
			if err := SendVoiceNotification(envName, alertName, time); err != nil {
				slog.Error(fmt.Sprintln("SendVoiceNotification error:", err))
				return err
			}
		} else {
			slog.Info(fmt.Sprintf("该环境: %s 存在于黑名单列表中，不进行通知。", envName))
		}

	case "whitelist":
		if isInWhiteList(envName) {
			if err := SendVoiceNotification(envName, alertName, time); err != nil {
				slog.Error(fmt.Sprintln("SendVoiceNotification error:", err))
				return err
			}
		} else {
			slog.Info(fmt.Sprintf("该环境: %s 不存在于白名单列表中，不进行通知。", envName))
		}
	default:
		slog.Error(fmt.Sprintln("不支持的模式，请配置 blacklist 或者 whitelist"))
		return errors.New("不支持的模式，请配置 blacklist 或者 whitelist")
	}
	return nil
}

func HandleOldVersionNotification(alertData *map[string]interface{}) error {
	slog.Info("按照旧版消息结构处理")
	//# 发送告警
	//# 报警通知, ${TXT_10}环境的${TXT_12}服务于${TIME}出现紧急异常.
	//# 参考：https://support.huaweicloud.com/VoiceCall_faq/VoiceCall_faq_00020.html
	//# 模板采用UTF-8编码格式，汉字和中文符号为3个字节，字母、数字和英文符号为1个字节。

	ruleName := (*alertData)["ruleName"].(string)
	state := (*alertData)["state"].(string)

	// 获取当前时间
	currentTime := time.Now()

	// 格式化为 HH:MM:SS
	formattedTime := currentTime.Format("15:04:05")

	envName, alertName := envNameAndAlertNameFromRuleName(ruleName)

	// 发送消息
	switch state {
	case "alerting":
		if err := CheckAndRunByMode(envName, alertName, formattedTime); err != nil {
			return err
		}
	case "ok":
		// 检查当前时间是否禁止发送
		if !utils.IsSendingAllowed() {
			slog.Info("当前时间不在允许发送范围内，警报恢复，查找并删除记录")
			// 处理数据
			envName = ProcessVoiceCallVariables(envName, 10)
			alertName = ProcessVoiceCallVariables(alertName, 12)
			// 创建一个字符串切片
			strList := []string{envName, alertName, formattedTime}
			slog.Info(fmt.Sprintf("删除前记录警报为：%s", datastore.RecordedData))
			// 根据前两个值，删除recordData中所有匹配的数据
			datastore.RemoveRecordsByPrefix(strList)
			slog.Info(fmt.Sprintf("删除后记录警报为：%s", datastore.RecordedData))
		}
	}
	return nil
}

func HandleNewVersionNotification(alertData *map[string]interface{}) error {
	// 处理新版告警消息结构
	slog.Info("按照新版消息结构处理")

	alertsSlice, ok := (*alertData)["alerts"].([]interface{})
	if !ok {
		// 处理类型断言失败的情况
		return errors.New("无法将 alerts 转换为 []interface{}")
	}

	for i, v := range alertsSlice {
		var ruleName string
		alert, ok := v.(map[string]interface{})
		if !ok {
			// 处理类型断言失败的情况
			slog.Error(fmt.Sprintln("无法将 alert 转换为 map[string]interface{}"))
			continue
		}
		if i == 0 {
			labels := alert["labels"].(map[string]interface{})
			status := alert["status"].(string)

			if labels["alertname"].(string) == "DatasourceNoData" || labels["alertname"].(string) == "DatasourceError" {
				slog.Info("DatasourceNoData|DatasourceError alert 不发送通知。")
				return nil
			}

			if labels["alertname"].(string) == "DatasourceNoData" {
				ruleName = labels["rulename"].(string)
			} else if labels["alertname"].(string) == "DatasourceError" {
				ruleName = labels["rulename"].(string)
			} else {
				ruleName = labels["alertname"].(string)
			}

			envName, alertName := envNameAndAlertNameFromRuleName(ruleName)
			// 取 startsAt 作为告警时间,转换为 HH:MM:SS 格式
			startsAt := alert["startsAt"].(string)

			// 解析时间字符串
			parsedTime, err := time.Parse(time.RFC3339, startsAt)
			if err != nil {
				slog.Error(fmt.Sprintf("解析时间字符串出错: %s", err))
				return err
			}

			// 格式化为 HH:MM:SS
			formattedTime := parsedTime.Format("15:04:05")

			// 发送消息
			switch status {
			case "firing":
				if err := CheckAndRunByMode(envName, alertName, formattedTime); err != nil {
					return err
				}
			case "resolved":
				// 检查当前时间是否禁止发送
				if !utils.IsSendingAllowed() {
					slog.Info("当前时间不在允许发送范围内，警报恢复，查找并删除记录")
					slog.Info(fmt.Sprintf("删除前记录警报为：%s", datastore.RecordedData))
					if err := checkRecordData(envName, alertName, formattedTime); err != nil {
						slog.Error(fmt.Sprintf("checkRecordData : %s", err))
						return err
					}
					slog.Info(fmt.Sprintf("删除后记录警报为：%s", datastore.RecordedData))
				}
			}
		}
	}

	return nil
}

func SelectSendNotifications(s []byte) error {
	// 判断警报信息结构属于老的还是新的，进行处理
	// alertData 动态结构体处理不同版本警报消息结构
	var alertData map[string]interface{}

	err := json.Unmarshal(s, &alertData)
	if err != nil {
		// 处理解析错误
		slog.Error(fmt.Sprintln("JSON marshal error:", err))
		return err
	}

	version := alertData["version"]

	if version == "1" {
		if err := HandleNewVersionNotification(&alertData); err != nil {
			slog.Error(fmt.Sprintln("HandleNewVersionNotification:", err))
			return err
		}
	}

	if version == nil {
		if err := HandleOldVersionNotification(&alertData); err != nil {
			slog.Error(fmt.Sprintln("HandleOldVersionNotification:", err))
			return err
		}
	}

	return nil
}

// SendRecordedAlertAtScheduledTime 定时发送记录的Alert
func SendRecordedAlertAtScheduledTime(dailyCheckTime time.Time) {
	nextClearTime := dailyCheckTime
	for {
		// 计算距离下一次清空的时间间隔
		now := time.Now()

		if nextClearTime.Before(now) {
			nextClearTime = nextClearTime.Add(24 * time.Hour)
		}
		timeUntilNextClear := nextClearTime.Sub(now)

		// 创建定时器，到达指定时间后执行清空操作
		timer := time.NewTimer(timeUntilNextClear)
		startTime := <-timer.C

		slog.Info("start a clear recorded alert")
		// 定时发送记录的未恢复警报
		for i, v := range datastore.RecordedData {
			slog.Info(fmt.Sprintf("%d %s", i, v))
			SendVoiceNotificationToAll(v)
		}
		datastore.RemoveRecordAll()
		slog.Info(fmt.Sprintf("Recorded alerting cleared at %s", startTime))
	}
}
