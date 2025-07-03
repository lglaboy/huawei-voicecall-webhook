package api

import (
	"fmt"
	"golang.org/x/exp/slog"
	"huawei-voicecall-webhook/config"
	"huawei-voicecall-webhook/internal"
	"io"
	"net/http"
)

func getServerPort() int {
	cfg := config.Cfg

	if cfg.Server.Port == 0 {
		return 8080
	} else {
		return cfg.Server.Port
	}
}

func StartWebhookServer() error {
	port := getServerPort()
	http.HandleFunc("/", handleWebhook)
	http.HandleFunc("/status", handleStatus)
	//port := 8080
	slog.Info(fmt.Sprintf("Webhook server listening on port %d ...", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		return err
	}
	return nil
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading request body: %v", err)
		return
	}

	slog.Info(fmt.Sprintf("Received Webhook Payload:\n%s\n", body))

	// 在这里处理来自 Grafana 的告警数据
	// 您可以解析 JSON 数据，触发相应的操作
	if err := internal.SelectSendNotifications(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "通知发送失败")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Webhook received and processed successfully")
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading request body: %v", err)
		return
	}

	slog.Info(fmt.Sprintf("Received Webhook Payload:\n%s\n", body))

	// 在这里处理来自 Grafana 的告警数据
	// 您可以解析 JSON 数据，触发相应的操作

	if err := internal.OnCallEvent(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "通知状态处理失败")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Webhook received and processed successfully")
}
