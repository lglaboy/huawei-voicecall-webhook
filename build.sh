#!/usr/bin/env bash
#
# File: build
# Date: 2024/6/19 11:24
# Author: whoami
# Mail: [邮箱地址]
# Description: 编译服务 构建镜像

# 定义变量
IMAGE_NAME="huawei-voice-notification:v5.1"
FILE_NAME="cmd/app.go"

# 编译
# 检查是否存在go环境
if which go >/dev/null 2>&1; then
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app $FILE_NAME
else
  echo "go 编译环境不存在"
  exit 1
fi

# 打包
# 判断docker是否存在
if which docker >/dev/null 2>&1; then
  docker build -t $IMAGE_NAME .
else
  echo "docker 命令不存在"
  exit 1
fi

# 打包格式
# 是否直接上传
