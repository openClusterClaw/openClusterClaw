.PHONY: help dev dev-backend dev-frontend build build-frontend prepare-embed build-backend cleanup-embed test clean docker db-shell install-tools happy-restart

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}

dev: ## 同时启动后端和前端开发服务器
	@echo "启动开发服务器..."
	@make dev-frontend &
	@sleep 2
	@make dev-backend

dev-backend: ## 启动后端开发服务器 (Air 热重载)
	@echo "启动后端服务..."
	@which air || go install github.com/air-verse/air@latest
	@air

dev-frontend: ## 启动前端开发服务器
	@echo "启动前端服务..."
	@cd frontend && pnpm install --silent && pnpm run dev

build-frontend: ## 构建前端静态文件
	@echo "构建前端..."
	@cd frontend && pnpm install --silent
	@cd frontend && pnpm run build

prepare-embed: ## 准备嵌入目录（复制前端文件）
	@echo "准备嵌入目录..."
	@mkdir -p internal/embed/ui/dist
	@rm -rf internal/embed/ui/dist/*
	@cp -r frontend/dist/* internal/embed/ui/dist/
	@echo "前端文件已复制到 internal/embed/ui/dist/"

build-backend: ## 构建后端（包含嵌入的前端）
	@echo "构建后端应用..."
	@go build -o bin/controlplane cmd/controlplane/main.go

build: build-frontend prepare-embed build-backend cleanup-embed ## 完整构建：前端+后端
	@echo "构建完成！"

cleanup-embed: ## 清理嵌入目录
	@echo "清理嵌入目录..."
	@rm -rf internal/embed/ui/dist/*
	@touch internal/embed/ui/dist/.keep
	@echo "嵌入目录已清理"

test: ## 运行测试
	go test -v ./...

clean: ## 清理构建产物
	rm -rf bin/
	rm -rf data/
	rm -rf frontend/dist/
	rm -rf frontend/node_modules/.vite 2>/dev/null || true
	rm -rf internal/embed/ui/dist/*
	@touch internal/embed/ui/dist/.keep

docker-build: ## 构建 Docker 镜像
	docker build -t open-cluster-claw:latest .

docker-run: ## 运行 Docker 容器
	docker run -p 8080:8080 open-cluster-claw:latest

db-shell: ## 打开数据库命令行
	sqlite3 data/clusterclaw.db

install-tools: ## 安装开发工具
	go install github.com/air-verse/air@latest

deps: ## 下载依赖
	go mod download
	cd frontend && pnpm install

tidy: ## 整理依赖
	go mod tidy

happy-restart: ## 重启 Happy
	@./scripts/restart-happy.sh