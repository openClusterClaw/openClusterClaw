.PHONY: help install dev build preview lint format

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## 安装依赖
	pnpm install

dev: ## 启动开发服务器
	pnpm dev

build: ## 构建生产版本
	pnpm build

preview: ## 预览生产构建
	pnpm preview

lint: ## 代码检查
	pnpm lint

format: ## 代码格式化
	pnpm format

clean: ## 清理构建产物
	rm -rf dist