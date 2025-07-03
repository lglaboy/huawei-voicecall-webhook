package main

import (
	"flag"
	"fmt"
	"golang.org/x/exp/slog"
	"huawei-voicecall-webhook/api"
	"huawei-voicecall-webhook/config"
	"huawei-voicecall-webhook/datastore"
	"huawei-voicecall-webhook/internal"
	"huawei-voicecall-webhook/utils"
	"log"
)

var configFile string

func init() {
	// 使用 flag 包定义 -c 参数，指定配置文件路径
	flag.StringVar(&configFile, "c", "/opt/webhook/config.yaml", "Path to the configuration file")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	// 打印配置文件路径，供你后续使用
	slog.Info(fmt.Sprintf("Using config file: %s", configFile))

	// 初始化配置
	if err := config.InitConfig(configFile); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	// 启动后台任务，按照指定间隔时间进行清理
	cleanupInterval, err := utils.ParseDuration(config.Cfg.Datastore.CleanupInterval)
	if err != nil {
		slog.Error(fmt.Sprintf("解析错误：%v\n", err))
	}
	go datastore.StartDataCleanup(cleanupInterval)

	// 启动后台任务，每天指定时间发送未发送告警信息
	// 解析每天的检查时间
	dailyCheckTime, err := utils.ParseConfigTime(config.Cfg.Common.DailyCheckTime)
	if err != nil {
		slog.Error(fmt.Sprintf("起始时间解析错误: %v\n", err))
	}
	// 启动定时器
	go internal.SendRecordedAlertAtScheduledTime(dailyCheckTime)

	// 启动 Webhook 服务
	err = api.StartWebhookServer()
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting server: %v", err))
	}
}
