# GPanel Makefile

.PHONY: help build run clean test dev docker-build docker-run docker-clean install-deps

# 默认目标
help:
	@echo "GPanel 服务器管理面板"
	@echo ""
	@echo "可用命令:"
	@echo "  help         显示此帮助信息"
	@echo "  dev          启动开发环境"
	@echo "  build        构建应用程序"
	@echo "  run          运行应用程序"
	@echo "  test         运行测试"
	@echo "  clean        清理构建文件"
	@echo "  install-deps 安装依赖"
	@echo "  docker-build 构建Docker镜像"
	@echo "  docker-run   运行Docker容器"
	@echo "  docker-clean 清理Docker资源"
	@echo "  release      构建发布版本"

# 变量定义
BINARY_NAME=gpanel
MAIN_PATH=cmd/server/main.go
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# 安装依赖
install-deps:
	@echo "安装Go依赖..."
	go mod download
	go mod tidy
	@echo "安装前端依赖..."
	cd web && npm install

# 开发环境
dev:
	@echo "启动开发环境..."
	@echo "启动后端服务..."
	go run $(MAIN_PATH) &
	@echo "启动前端服务..."
	cd web && npm run dev

# 构建应用程序
build:
	@echo "构建应用程序..."
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "构建前端..."
	cd web && npm run build
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行应用程序
run:
	@echo "运行应用程序..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# 测试
test:
	@echo "运行测试..."
	go test -v ./...

# 代码格式化
fmt:
	@echo "格式化代码..."
	go fmt ./...
	cd web && npm run format

# 代码检查
lint:
	@echo "检查代码..."
	golangci-lint run
	cd web && npm run lint

# 清理
clean:
	@echo "清理构建文件..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	cd web && rm -rf dist node_modules

# Docker构建
docker-build:
	@echo "构建Docker镜像..."
	docker build -t gpanel:$(VERSION) .
	docker tag gpanel:$(VERSION) gpanel:latest

# Docker运行
docker-run:
	@echo "运行Docker容器..."
	docker-compose up -d

# Docker清理
docker-clean:
	@echo "清理Docker资源..."
	docker-compose down -v
	docker system prune -f

# 发布版本
release: clean
	@echo "构建发布版本..."
	mkdir -p $(BUILD_DIR)/release
	# 构建多平台二进制文件
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	# 构建前端
	cd web && npm run build
	# 复制必要文件
	cp -r web/dist $(BUILD_DIR)/release/
	cp -r templates $(BUILD_DIR)/release/
	cp -r scripts $(BUILD_DIR)/release/
	cp docker-compose.yaml $(BUILD_DIR)/release/
	cp Dockerfile $(BUILD_DIR)/release/
	cp README.md $(BUILD_DIR)/release/
	# 创建压缩包
	cd $(BUILD_DIR)/release && tar -czf ../gpanel-$(VERSION).tar.gz .
	@echo "发布版本构建完成: $(BUILD_DIR)/gpanel-$(VERSION).tar.gz"

# 安装到系统
install: build
	@echo "安装到系统..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	sudo mkdir -p /etc/gpanel
	sudo cp internal/config/config.yaml /etc/gpanel/
	sudo mkdir -p /var/lib/gpanel
	sudo cp -r templates /var/lib/gpanel/
	sudo cp scripts/gpanel.service /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable gpanel
	@echo "安装完成! 使用 'systemctl start gpanel' 启动服务"

# 卸载
uninstall:
	@echo "从系统卸载..."
	sudo systemctl stop gpanel || true
	sudo systemctl disable gpanel || true
	sudo rm -f /etc/systemd/system/gpanel.service
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	sudo rm -rf /etc/gpanel
	sudo rm -rf /var/lib/gpanel
	sudo systemctl daemon-reload
	@echo "卸载完成!"

# 生成API文档
docs:
	@echo "生成API文档..."
	swag init -g cmd/server/main.go

# 数据库迁移
migrate:
	@echo "执行数据库迁移..."
	go run cmd/migrate/main.go

# 备份数据
backup:
	@echo "备份数据..."
	./scripts/backup.sh

# 恢复数据
restore:
	@echo "恢复数据..."
	./scripts/restore.sh

# 监控
monitor:
	@echo "启动监控..."
	./scripts/monitor.sh

# 日志查看
logs:
	docker-compose logs -f gpanel

# 进入容器
shell:
	docker exec -it gpanel sh