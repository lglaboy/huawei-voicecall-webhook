package config

import (
	yaml "gopkg.in/yaml.v2"
	"os"
)

var Cfg *Config

type Config struct {
	Lead        []Lead        `yaml:"lead"`
	Huawei      Huawei        `yaml:"huawei"`
	Server      Server        `yaml:"server"`
	Common      Common        `yaml:"common"`
	Datastore   Datastore     `yaml:"datastore"`
	Schedule    Schedule      `yaml:"schedule"`
	WhiteList   []string      `yaml:"white_list"`
	BlackList   []string      `yaml:"black_list"`
	EnvDirector []EnvDirector `yaml:"env_director"`
}

type Lead struct {
	Name  string `yaml:"name"`
	Phone string `yaml:"phone"`
}

type Huawei struct {
	Voicecall Voicecall `yaml:"voicecall"`
}

type Voicecall struct {
	BaseUrl    string `yaml:"base_url"`
	AppKey     string `yaml:"appKey"`
	AppSecret  string `yaml:"appSecret"`
	DisplayNbr string `yaml:"displayNbr"`
	TemplateId string `yaml:"templateId"`
}

type Server struct {
	Port int `yaml:"port"`
}

type Common struct {
	StatusUrl       string          `yaml:"statusUrl"`
	Model           string          `yaml:"model"`
	ForbiddenPeriod ForbiddenPeriod `yaml:"forbiddenPeriod"`
	DailyCheckTime  string          `yaml:"dailyCheckTime"`
}

type ForbiddenPeriod struct {
	StartTime string `yaml:"startTime"`
	EndTime   string `yaml:"endTime"`
}

type Datastore struct {
	TTL             string `yaml:"ttl"`
	CleanupInterval string `yaml:"cleanup_interval"`
}

type Schedule struct {
	File string `yaml:"file"`
}

type EnvDirector struct {
	Env   string `yaml:"env"`
	Name  string `yaml:"name"`
	Phone string `yaml:"phone"`
}

// LoadConfig 从指定的 YAML 文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func InitConfig(filePath string) error {
	cfg, err := LoadConfig(filePath)
	if err != nil {
		return err
	}
	Cfg = cfg
	return nil
}
