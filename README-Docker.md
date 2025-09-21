# GameHub 容器化部署指南

## 📋 目錄
- [系統需求](#系統需求)
- [快速開始](#快速開始)
- [開發環境](#開發環境)
- [生產環境](#生產環境)
- [配置說明](#配置說明)
- [常用命令](#常用命令)
- [故障排除](#故障排除)
- [性能優化](#性能優化)

## 🔧 系統需求

### 最低要求
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **內存**: 4GB RAM
- **存儲**: 20GB 可用空間
- **操作系統**: Linux/macOS/Windows

### 推薦配置
- **CPU**: 4 核心以上
- **內存**: 8GB RAM
- **存儲**: SSD 50GB 以上
- **網路**: 1Gbps

## 🚀 快速開始

### 1. 克隆專案並設置環境

```bash
# 克隆專案
git clone <your-repo-url>
cd stone

# 複製環境變量配置
cp .env.example .env

# 編輯環境變量（根據需要修改）
vim .env
```

### 2. 一鍵部署

```bash
# 開發環境部署
./scripts/deploy.sh deploy dev

# 生產環境部署
./scripts/deploy.sh deploy prod
```

### 3. 驗證部署

```bash
# 檢查服務狀態
./scripts/deploy.sh health

# 查看服務日誌
./scripts/deploy.sh logs
```

訪問服務：
- **遊戲服務**: http://localhost
- **API 文檔**: http://localhost:8083 (開發環境)
- **數據庫管理**: http://localhost:5050 (開發環境)
- **Redis 管理**: http://localhost:8081 (開發環境)

## 🛠️ 開發環境

### 啟動開發環境

```bash
# 使用開發配置啟動
docker-compose -f docker-compose.dev.yml up -d

# 或使用部署腳本
./scripts/deploy.sh deploy dev
```

### 開發特性

- **熱重載**: 代碼更改後自動重新編譯
- **調試端口**: 可連接調試器
- **開發工具**: 包含 pgAdmin、Redis Commander
- **詳細日誌**: Debug 級別日誌輸出

### 開發工具訪問

| 服務 | URL | 用戶名 | 密碼 |
|------|-----|--------|------|
| pgAdmin | http://localhost:5050 | admin@gamehub.dev | admin123 |
| Redis Commander | http://localhost:8081 | - | - |
| API 文檔 | http://localhost:8083 | - | - |
| 老虎機編輯器 | http://localhost:8082 | - | - |

### 調試代碼

```bash
# 進入容器調試
docker exec -it gamehub-server-dev bash

# 查看實時日誌
docker logs -f gamehub-server-dev

# 重啟特定服務
docker-compose -f docker-compose.dev.yml restart gamehub-dev
```

## 🏭 生產環境

### 準備生產環境

1. **設置環境變量**
```bash
# 編輯生產環境配置
cp .env.example .env.prod
vim .env.prod
```

2. **SSL 證書**
```bash
# 放置 SSL 證書
mkdir -p docker/nginx/ssl
cp your-cert.pem docker/nginx/ssl/cert.pem
cp your-key.pem docker/nginx/ssl/key.pem
```

3. **安全配置**
```bash
# 設置安全密碼
export POSTGRES_PASSWORD=$(openssl rand -base64 32)
export REDIS_PASSWORD=$(openssl rand -base64 32)
```

### 生產環境部署

```bash
# 使用生產配置
export COMPOSE_FILE=docker-compose.yml
export ENV_FILE=.env.prod

# 部署
./scripts/deploy.sh deploy prod
```

### 監控和維護

```bash
# 查看系統狀態
docker stats

# 查看服務日誌
docker-compose logs -f --tail=100

# 備份數據
./scripts/deploy.sh backup

# 健康檢查
./scripts/deploy.sh health
```

## ⚙️ 配置說明

### 環境變量配置

主要配置項說明：

```bash
# 數據庫配置
POSTGRES_DB=gamehub           # 數據庫名稱
POSTGRES_USER=gamehub         # 數據庫用戶
POSTGRES_PASSWORD=xxx         # 數據庫密碼

# Redis 配置  
REDIS_PASSWORD=xxx            # Redis 密碼

# 服務配置
SERVER_ID=1                   # 服務器 ID
PLATFORM=DEV                  # 平台環境 (DEV/QA/PROD)
LOG_LEVEL=3                   # 日誌級別 (1-5)

# 端口配置
GAMEHUB_WS_PORT=3563         # WebSocket 端口
GAMEHUB_TCP_PORT=3564        # TCP 端口
GAMEHUB_HTTP_PORT=8080       # HTTP 端口
```

### 服務配置文件

- **GameHub**: `docker/config/GameHub.conf`
- **Nginx**: `docker/nginx/nginx.conf`
- **Redis**: `docker/redis/redis.conf`

## 📝 常用命令

### 基本操作

```bash
# 啟動所有服務
docker-compose up -d

# 停止所有服務
docker-compose down

# 重啟服務
docker-compose restart [service_name]

# 查看服務狀態
docker-compose ps

# 查看服務日誌
docker-compose logs [service_name]
```

### 部署腳本操作

```bash
# 完整部署
./scripts/deploy.sh deploy [dev|staging|prod]

# 僅啟動服務
./scripts/deploy.sh start [env]

# 停止服務
./scripts/deploy.sh stop

# 重新構建映像
./scripts/deploy.sh build

# 備份數據
./scripts/deploy.sh backup

# 還原數據
./scripts/deploy.sh restore backup_directory

# 清理所有資源
./scripts/deploy.sh cleanup
```

### 維護操作

```bash
# 更新服務
docker-compose pull
docker-compose up -d

# 清理舊映像
docker image prune -f

# 查看資源使用
docker system df

# 數據庫操作
docker exec -it gamehub-postgres psql -U gamehub -d gamehub

# Redis 操作  
docker exec -it gamehub-redis redis-cli
```

## 🔍 故障排除

### 常見問題

#### 1. 容器啟動失敗

```bash
# 查看詳細錯誤
docker-compose logs [service_name]

# 檢查配置文件
docker-compose config

# 重新構建映像
docker-compose build --no-cache [service_name]
```

#### 2. 數據庫連接失敗

```bash
# 檢查數據庫狀態
docker exec gamehub-postgres pg_isready -U gamehub

# 查看數據庫日誌
docker logs gamehub-postgres

# 重置數據庫密碼
docker exec -it gamehub-postgres psql -U postgres
```

#### 3. Redis 連接問題

```bash
# 測試 Redis 連接
docker exec gamehub-redis redis-cli ping

# 檢查 Redis 配置
docker exec gamehub-redis redis-cli CONFIG GET "*"
```

#### 4. 端口衝突

```bash
# 檢查端口使用
netstat -tulpn | grep :3563

# 修改 .env 文件中的端口配置
vim .env
```

### 日誌分析

```bash
# 查看所有服務日誌
docker-compose logs --follow

# 查看特定時間範圍的日誌
docker-compose logs --since="2024-01-01T00:00:00Z" --until="2024-01-01T23:59:59Z"

# 搜索錯誤日誌
docker-compose logs | grep -i error

# 保存日誌到文件
docker-compose logs > logs/debug_$(date +%Y%m%d_%H%M%S).log
```

## 🚀 性能優化

### 資源限制

在 `docker-compose.yml` 中添加資源限制：

```yaml
services:
  gamehub:
    deploy:
      resources:
        limits:
          memory: 2G
          cpus: '1.0'
        reservations:
          memory: 1G
          cpus: '0.5'
```

### 數據庫優化

```sql
-- PostgreSQL 性能調優
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
SELECT pg_reload_conf();
```

### Redis 優化

```bash
# Redis 內存優化
docker exec gamehub-redis redis-cli CONFIG SET maxmemory 512mb
docker exec gamehub-redis redis-cli CONFIG SET maxmemory-policy allkeys-lru
```

### 監控指標

```bash
# 監控容器資源使用
docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"

# 監控數據庫性能
docker exec gamehub-postgres psql -U gamehub -c "SELECT * FROM pg_stat_activity;"

# 監控 Redis 性能
docker exec gamehub-redis redis-cli INFO stats
```

## 📚 更多資源

- [Docker 官方文檔](https://docs.docker.com/)
- [Docker Compose 參考](https://docs.docker.com/compose/)
- [PostgreSQL Docker 鏡像](https://hub.docker.com/_/postgres)
- [Redis Docker 鏡像](https://hub.docker.com/_/redis)
- [Nginx Docker 鏡像](https://hub.docker.com/_/nginx)

## 🤝 貢獻

如果你發現問題或有改進建議，請提交 Issue 或 Pull Request。

## 📄 許可證

請參考專案根目錄的 LICENSE 文件。