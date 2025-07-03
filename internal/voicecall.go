package internal

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"huawei-voicecall-webhook/config"
	"io"
	"net/http"
	"strings"
	"time"
)

func buildAKSKHeader(appKey, appSecret string) string {
	now := time.Now().Format("2006-01-02T15:04:05Z")           // Created
	nonce := strings.Replace(uuid.New().String(), "-", "", -1) // Nonce

	data := []byte(nonce + now)
	digest := hmac.New(sha256.New, []byte(appSecret))
	digest.Write(data)
	digestBytes := digest.Sum(nil)
	digestBase64 := base64.StdEncoding.EncodeToString(digestBytes) // PasswordDigest

	return fmt.Sprintf(`UsernameToken Username="%s", PasswordDigest="%s", Nonce="%s", Created="%s"`, appKey, digestBase64, nonce, now)
}

func voiceNotifyAPI(displayNbr string, calleeNbr string, playInfoList []PlayInfo) {
	if len(displayNbr) < 1 || len(calleeNbr) < 1 || len(playInfoList) == 0 {
		return
	}

	apiUri := "/rest/httpsessions/callnotify/v2.0" // v1.0 or v2.0
	requestURL := config.Cfg.Huawei.Voicecall.BaseUrl + apiUri

	jsonData := map[string]interface{}{
		"displayNbr":   displayNbr,
		"calleeNbr":    calleeNbr,
		"playInfoList": playInfoList,
		// 设置SP接收状态上报的URL,要求使用BASE64编码
		"statusUrl": base64.StdEncoding.EncodeToString([]byte(config.Cfg.Common.StatusUrl)),
	}

	// 编码JSON数据为字节数组
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		slog.Error(fmt.Sprintln("编码JSON数据时出错:", err))
		return
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonBytes))

	if err != nil {
		slog.Error(fmt.Sprintln("创建请求时出错:", err))
		return
	}

	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Authorization", `AKSK realm="SDP",profile="UsernameToken",type="Appkey"`)
	req.Header.Add("X-AKSK", buildAKSKHeader(config.Cfg.Huawei.Voicecall.AppKey, config.Cfg.Huawei.Voicecall.AppSecret))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error(fmt.Sprintln("发送请求时出错:", err))
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//return "", err
		slog.Error(fmt.Sprintln(err))
	}
	slog.Info(fmt.Sprintln(string(body)))
}

func getPlayInfoList(templateID string, templateParas []string) []PlayInfo {
	playInfo := PlayInfo{
		TemplateID:    templateID,
		TemplateParas: templateParas,
	}

	return []PlayInfo{playInfo}
}

type PlayInfo struct {
	//NotifyVoice   string   `json:"notifyVoice"`
	TemplateID    string   `json:"templateId"`
	TemplateParas []string `json:"templateParas"`
	// 你可以添加更多字段
}

func SendVoiceNotify(calleeNbr string, body []string) {
	playInfoList := getPlayInfoList(config.Cfg.Huawei.Voicecall.TemplateId, body)
	voiceNotifyAPI(config.Cfg.Huawei.Voicecall.DisplayNbr, calleeNbr, playInfoList)
}
