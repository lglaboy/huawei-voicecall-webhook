## 说明

本项目通过创建 Webhook 接口，将 Grafana 告警信息对接至华为云平台 语音通话_VoiceCall 服务，实现语音通知功能。

**功能**

1. 用于Grafana，作为webhook类型告警方式，通过电话通知异常事件。
2. 拨打方式: 无人接通(连续拨打三次 结束)，有人接通(通知完成后结束)。
3. 黑白名单策略，针对指定环境禁止或允许发送通知。
4. 可指定禁止发送通知的时间段，比如：每天 02:00:00 - 06:00:00 出现的异常，不进行通知
5. 指定时间点检查是否存在异常通知，进行告警。
   
   如果指定每天 06:00:00 检查，禁止发送时间（02:00:00 - 06:00:00）内的告警在06:00:00 前恢复，则不告警。
6. 可根据值班表，给指定值班人员发送通知。动态获取。
7. 区分 领导电话、项目负责人电话、值班人员电话。

## 编译

编译命令

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/app.go
```

启动选项

| 选项 | 默认值                        | 用途     |
|----|----------------------------|--------|
| -c | `/opt/webhook/config.yaml` | 指定配置文件 |

配置文件默认采用 `/opt/webhook/config.yaml`，可使用 `-c` 指定配置文件，如:

```shell
./app -c config.yaml
```

## 封装镜像

```shell
docker build -t huawei-voice-notification:v5 .
```

# 启动

```shell
# 创建本地挂载目录
mkdir -p /opt/huawei-voicecall-webhook/

# 创建配置文件
vim /opt/huawei-voicecall-webhook/config.yml

# 创建值班文件
vim /opt/huawei-voicecall-webhook/schedule.csv

# 启动服务
docker run -itd --name huawei-voicecall-webhook \
--cpus 1 -m 2G \
--log-opt max-size=512m \
--log-opt max-file=3 \
--restart=unless-stopped \
-v /etc/localtime:/etc/localtime \
-v /etc/hosts:/etc/hosts \
-v /opt/huawei-voicecall-webhook/config.yml:/opt/webhook/config.yaml \
-v /opt/huawei-voicecall-webhook/schedule.csv:/opt/webhook/schedule.csv \
-p 8092:8080 \
huawei-voice-notification:v5 -c /opt/webhook/config.yaml
```

## 配置Grafana，添加Webhook语音告警

grafana -> Alerting -> Contact points -> 添加一个联络点

联络点配置：

- Name: 自定义
- Integration: Webhook
- URL: `http://192.168.*.*:8092`
