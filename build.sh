#!/bin/bash

# 版本号
VERSION="v1.0.0"
# 程序名称
BINARY_NAME="cfspeedtest"

# 创建构建目录
mkdir -p build

# 构建各平台二进制
GOOS=windows GOARCH=amd64 go build -o build/${BINARY_NAME}_${VERSION}_windows_amd64.exe
GOOS=windows GOARCH=386 go build -o build/${BINARY_NAME}_${VERSION}_windows_386.exe
GOOS=linux GOARCH=amd64 go build -o build/${BINARY_NAME}_${VERSION}_linux_amd64
GOOS=linux GOARCH=386 go build -o build/${BINARY_NAME}_${VERSION}_linux_386
GOOS=darwin GOARCH=amd64 go build -o build/${BINARY_NAME}_${VERSION}_darwin_amd64
GOOS=darwin GOARCH=arm64 go build -o build/${BINARY_NAME}_${VERSION}_darwin_arm64

# 复制配置文件
cp -r configs build/
cp ip.txt build/ 