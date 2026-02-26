.PHONY: help dev dev-backend dev-frontend build test clean docker db-shell install-tools happy-restart

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

build: ## 构建应用
	go build -o bin/controlplane cmd/controlplane/main.go

test: ## 运行测试
	go test -v ./...

clean: ## 清理构建产物
	rm -rf bin/
	rm -rf data/
	cd frontend && rm -rf dist/ node_modules/.vite 2>/dev/null || true

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