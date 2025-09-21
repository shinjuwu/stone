# GameHub 容器化 Makefile

.PHONY: help dev prod build start stop restart logs clean backup restore health

# 預設目標
.DEFAULT_GOAL := help

# 環境變量
COMPOSE_FILE ?= docker-compose.yml
ENV_FILE ?= .env
PROJECT_NAME ?= gamehub

help: ## 顯示幫助信息
	@echo "GameHub 容器化管理命令："
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 開發環境
dev: ## 啟動開發環境
	@echo "🚀 啟動開發環境..."
	@docker-compose -f docker-compose.dev.yml up -d
	@echo "✅ 開發環境已啟動"
	@echo "🌐 訪問地址："
	@echo "   - 主服務: http://localhost"
	@echo "   - pgAdmin: http://localhost:5050"
	@echo "   - Redis Commander: http://localhost:8081"
	@echo "   - API 文檔: http://localhost:8083"

dev-build: ## 構建並啟動開發環境
	@echo "🔨 構建開發環境..."
	@docker-compose -f docker-compose.dev.yml build
	@make dev

dev-stop: ## 停止開發環境
	@echo "⏹️  停止開發環境..."
	@docker-compose -f docker-compose.dev.yml down
	@echo "✅ 開發環境已停止"

dev-logs: ## 查看開發環境日誌
	@docker-compose -f docker-compose.dev.yml logs -f

# 生產環境
prod: ## 啟動生產環境
	@echo "🚀 啟動生產環境..."
	@docker-compose up -d postgres redis gamehub nginx
	@make health
	@echo "✅ 生產環境已啟動"

prod-build: ## 構建並啟動生產環境
	@echo "🔨 構建生產環境..."
	@docker-compose build
	@make prod

# 基本操作
build: ## 構建所有鏡像
	@echo "🔨 構建 Docker 鏡像..."
	@docker-compose build

start: ## 啟動服務
	@echo "▶️  啟動服務..."
	@docker-compose up -d

stop: ## 停止服務
	@echo "⏹️  停止服務..."
	@docker-compose down

restart: ## 重啟服務
	@echo "🔄 重啟服務..."
	@make stop
	@make start

logs: ## 查看日誌
	@docker-compose logs -f --tail=100

# 狀態和監控
ps: ## 顯示容器狀態
	@docker-compose ps

stats: ## 顯示資源使用統計
	@docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"

health: ## 執行健康檢查
	@echo "🔍 執行健康檢查..."
	@scripts/deploy.sh health

# 數據管理
backup: ## 備份數據
	@echo "💾 備份數據..."
	@scripts/deploy.sh backup

restore: ## 還原數據 (使用方法: make restore BACKUP_DIR=path/to/backup)
	@echo "🔄 還原數據..."
	@scripts/deploy.sh restore $(BACKUP_DIR)

# 清理操作
clean: ## 清理未使用的資源
	@echo "🧹 清理未使用的資源..."
	@docker system prune -f
	@docker volume prune -f

clean-all: ## 清理所有資源 (危險操作)
	@echo "⚠️  清理所有資源..."
	@read -p "確定要刪除所有容器、鏡像和數據嗎？(y/N): " confirm && [ "$$confirm" = "y" ]
	@scripts/deploy.sh cleanup

# 開發工具
shell: ## 進入 GameHub 容器 shell
	@docker exec -it gamehub-server bash

db-shell: ## 進入資料庫 shell
	@docker exec -it gamehub-postgres psql -U gamehub gamehub

redis-shell: ## 進入 Redis shell
	@docker exec -it gamehub-redis redis-cli

# 測試和調試
test: ## 運行測試
	@echo "🧪 運行測試..."
	@docker-compose exec gamehub go test ./...

lint: ## 運行代碼檢查
	@echo "🔍 運行代碼檢查..."
	@docker run --rm -v $(PWD)/GameHub:/app -w /app golangci/golangci-lint:latest golangci-lint run

# 更新操作
update: ## 更新服務
	@echo "⬆️  更新服務..."
	@docker-compose pull
	@docker-compose up -d

# 配置管理
config: ## 檢查配置
	@docker-compose config

env-setup: ## 設置環境變量
	@if [ ! -f .env ]; then \
		echo "📝 創建環境配置文件..."; \
		cp .env.example .env; \
		echo "✅ 請編輯 .env 文件並配置相應的值"; \
	else \
		echo "✅ 環境配置文件已存在"; \
	fi

# 監控相關
monitoring: ## 啟動監控服務
	@echo "📊 啟動監控服務..."
	@docker-compose --profile monitoring up -d prometheus grafana
	@echo "✅ 監控服務已啟動"
	@echo "🌐 Grafana: http://localhost:3000 (admin/admin123)"
	@echo "🌐 Prometheus: http://localhost:9090"

# 工具服務
tools: ## 啟動工具服務
	@echo "🛠️  啟動工具服務..."
	@docker-compose --profile tools up -d
	@echo "✅ 工具服務已啟動"

# 快速操作
quick-start: env-setup build start health ## 快速啟動 (設置環境 + 構建 + 啟動 + 健康檢查)

quick-dev: env-setup dev-build ## 快速開始開發 (設置環境 + 開發環境構建)

# 文檔生成
docs: ## 生成 API 文檔
	@echo "📚 生成 API 文檔..."
	@docker run --rm -v $(PWD):/workspace -w /workspace swaggerapi/swagger-codegen-cli generate \
		-i docs/api/swagger.yaml \
		-l html2 \
		-o docs/api/html

# 安全掃描
security-scan: ## 運行安全掃描
	@echo "🔒 運行安全掃描..."
	@docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
		-v $(PWD):/src aquasec/trivy fs /src

# 性能測試
benchmark: ## 運行性能測試
	@echo "⚡ 運行性能測試..."
	@docker run --rm --network=host -v $(PWD)/tests:/tests loadimpact/k6 run /tests/load-test.js