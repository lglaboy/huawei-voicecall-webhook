package internal

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slog"
	"huawei-voicecall-webhook/datastore"
)

// StatusInfo 包含呼叫状态事件信息
type StatusInfo struct {
	SessionID string `json:"sessionId"`
	Timestamp string `json:"timestamp"`
	Caller    string `json:"caller"`
	Called    string `json:"called"`
	UserData  string `json:"userData"`
	DigitInfo string `json:"digitInfo"`
	StateCode int    `json:"stateCode"`
	StateDesc string `json:"stateDesc"`
	// 其他字段...
}

// CallEvent 包含通知事件类型和状态信息
type CallEvent struct {
	EventType  string     `json:"eventType"`
	StatusInfo StatusInfo `json:"statusInfo"`
}

// OnCallEvent 处理呼叫事件
func OnCallEvent(jsonBody []byte) error {
	var callEvent CallEvent

	err := json.Unmarshal(jsonBody, &callEvent)
	if err != nil {
		fmt.Println("JSON解析错误：", err)
		return err
	}

	eventType := callEvent.EventType
	statusInfo := callEvent.StatusInfo

	slog.Info(fmt.Sprintln("eventType:", eventType)) // 打印通知事件类型

	switch eventType {
	case "callout":
		// 呼出事件处理
		if statusInfo.SessionID != "" {
			fmt.Println("sessionId:", statusInfo.SessionID)
		}
	case "alerting":
		// 振铃事件处理
		if statusInfo.SessionID != "" {
			fmt.Println("sessionId:", statusInfo.SessionID)
		}
	case "answer":
		// 应答事件处理
		if statusInfo.SessionID != "" {
			fmt.Println("sessionId:", statusInfo.SessionID)
		}
	case "collectInfo":
		// 放音收号结果事件处理
		if statusInfo.DigitInfo != "" {
			fmt.Println("digitInfo:", statusInfo.DigitInfo)
		}
	case "disconnect":
		// 挂机事件处理
		if statusInfo.SessionID != "" {
			fmt.Println("sessionId:", statusInfo.SessionID)
		}

		if statusInfo.StateCode == 0 {
			datastore.DeleteData(statusInfo.Called)
		} else {
			SendRepeatVoiceNotify(statusInfo.Called)
		}
	default:
		fmt.Println("EventType错误:", eventType)
	}
	return nil
}
