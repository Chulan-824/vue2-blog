.PHONY: build run clean test docker-build docker-run

# 变量定义
APP_NAME=vue2-blog-server
DOCKER_IMAGE=$(APP_NAME):latest

# 构建应用
build:
	go build -o bin/$(APP_NAME) ./cmd/api

# 运行应用
run:
	go run ./cmd/api

# 运行开发环境
dev:
	go run ./cmd/api --config=./configs/config.dev.yaml

# 清理构建文件
clean:
	rm -rf bin/

# 运行测试
test:
	go test -v ./...

# 下载依赖
deps:
	go mod download
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...

# 静态检查
lint:
	golangci-lint run

# 构建Docker镜像
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# 运行Docker容器
docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE)

# 生成API文档
docs:
	swag init -g ./cmd/api/main.go

# 初始化项目
init:
	go mod init vue2-blog-server
	go mod tidy

# 安装工具
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest