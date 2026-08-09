# Makefile for goat

# 项目信息
APP_NAME = goat
MODULE = github.com/tang980129/goat
VERSION ?= 1.0.0

# 构建时注入变量
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -X '$(MODULE)/cmd.version=$(VERSION)' \
           -X '$(MODULE)/cmd.commit=$(COMMIT)' \
           -X '$(MODULE)/cmd.buildDate=$(BUILD_DATE)'

# 默认目标
.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) .

# Windows 构建
.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(APP_NAME).exe .

# 清理生成文件
.PHONY: clean
clean:
	rm -f $(APP_NAME) $(APP_NAME).exe

# 运行测试
.PHONY: test
test:
	go test ./...

# 直接运行（开发模式）
.PHONY: run
run:
	go run .