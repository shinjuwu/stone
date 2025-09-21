# 🐧 GameHub Linux 部署完整指南

## 📋 目錄
- [系統需求與準備](#系統需求與準備)
- [部署方法比較](#部署方法比較)
- [容器化部署 (推薦)](#容器化部署-推薦)
- [原生部署](#原生部署)
- [生產環境優化](#生產環境優化)
- [監控與運維](#監控與運維)
- [故障排除](#故障排除)

## 🔧 系統需求與準備

### 最低系統需求
```bash
# 操作系統
Ubuntu 20.04+ / CentOS 8+ / RHEL 8+ / Debian 11+

# 硬件需求
CPU: 4 核心
內存: 8GB RAM
存儲: 50GB SSD
網絡: 1Gbps

# 必要軟件
Docker 20.10+
Docker Compose 2.0+
Git
```

### 系統準備
```bash
# Ubuntu/Debian 系統準備
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl wget git vim net-tools htop

# CentOS/RHEL 系統準備
sudo yum update -y
sudo yum install -y curl wget git vim net-tools htop

# 安裝 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# 安裝 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 驗證安裝
docker --version
docker-compose --version
```

## ⚖️ 部署方法比較

| 部署方式 | 優點 | 缺點 | 適用場景 |
|---------|------|------|---------|
| **容器化部署** | • 環境一致性<br>• 快速部署<br>• 易於擴展<br>• 隔離性好 | • 學習成本<br>• 資源開銷 | **推薦**<br>生產環境 |
| **原生部署** | • 性能最佳<br>• 資源利用率高<br>• 直接控制 | • 環境複雜<br>• 依賴管理困難 | 高性能需求 |
| **混合部署** | • 平衡性能和管理 | • 複雜度高 | 大型集群 |

## 🐳 容器化部署 (推薦)

### 快速部署
```bash
# 1. 克隆項目
git clone <your-repo-url>
cd stone

# 2. 一鍵部署開發環境
make quick-dev

# 3. 一鍵部署生產環境
make quick-start
```

### 詳細部署步驟

#### 1. 環境配置
```bash
# 複製並編輯環境變量
cp .env.example .env
vim .env

# 關鍵配置項
POSTGRES_PASSWORD=your_secure_password
REDIS_PASSWORD=your_redis_password
SERVER_ID=1
PLATFORM=PROD
LOG_LEVEL=2
```

#### 2. 生產環境配置
```bash
# 創建生產環境配置
cat > .env.prod << EOF
# 數據庫配置
POSTGRES_DB=gamehub_prod
POSTGRES_USER=gamehub_prod
POSTGRES_PASSWORD=$(openssl rand -base64 32)

# Redis 配置
REDIS_PASSWORD=$(openssl rand -base64 32)

# 服務配置
SERVER_ID=1
PLATFORM=PROD
LOG_LEVEL=2

# 端口配置 (生產環境)
GAMEHUB_WS_PORT=3563
GAMEHUB_TCP_PORT=3564
GAMEHUB_HTTP_PORT=8080
HTTP_PORT=80
HTTPS_PORT=443

# 監控配置
GRAFANA_PASSWORD=$(openssl rand -base64 16)
EOF
```

#### 3. SSL 證書配置
```bash
# 創建 SSL 目錄
mkdir -p docker/nginx/ssl

# 使用 Let's Encrypt (推薦)
sudo apt install certbot
sudo certbot certonly --standalone -d your-domain.com
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem docker/nginx/ssl/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem docker/nginx/ssl/key.pem

# 或使用自簽名證書 (測試用)
openssl req -x509 -newkey rsa:4096 -keyout docker/nginx/ssl/key.pem \
    -out docker/nginx/ssl/cert.pem -days 365 -nodes \
    -subj "/C=TW/ST=Taiwan/L=Taipei/O=GameHub/CN=localhost"
```

#### 4. 部署執行
```bash
# 方法一：使用 Makefile (推薦)
make prod

# 方法二：使用部署腳本
./scripts/deploy.sh deploy prod

# 方法三：直接使用 Docker Compose
export COMPOSE_FILE=docker-compose.yml
export ENV_FILE=.env.prod
docker-compose up -d postgres redis gamehub nginx
```

#### 5. 部署驗證
```bash
# 檢查服務狀態
make ps
# 或
docker-compose ps

# 健康檢查
make health
# 或
./scripts/deploy.sh health

# 查看日誌
make logs
# 或
docker-compose logs -f --tail=100
```

### 開發環境部署
```bash
# 啟動開發環境（包含管理工具）
make dev

# 開發工具訪問地址
echo "開發工具列表："
echo "• 主服務: http://localhost"
echo "• pgAdmin: http://localhost:5050 (admin@gamehub.dev/admin123)"
echo "• Redis Commander: http://localhost:8081"
echo "• API 文檔: http://localhost:8083"
echo "• 老虎機編輯器: http://localhost:8082"
```

## 🖥️ 原生部署

### 系統依賴安裝
```bash
# Ubuntu/Debian
sudo apt install -y golang-go postgresql-14 redis-server nginx

# CentOS/RHEL
sudo yum install -y golang postgresql14-server redis nginx

# 環境變量設置
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc
```

### 數據庫設置
```bash
# PostgreSQL 設置
sudo -u postgres createuser -P gamehub
sudo -u postgres createdb -O gamehub gamehub

# 導入數據庫結構
sudo -u postgres psql -d gamehub -f docker/init-scripts/01-init-database.sql
sudo -u postgres psql -d gamehub -f gameinfo.sql
sudo -u postgres psql -d gamehub -f gamelist.sql
sudo -u postgres psql -d gamehub -f lobbyinfo.sql

# Redis 設置
sudo systemctl enable redis
sudo systemctl start redis
```

### 應用編譯與部署
```bash
# 編譯 GameHub
cd GameHub
CGO_ENABLED=0 GOOS=linux go build -o gamehub .

# 創建運行目錄
sudo mkdir -p /opt/gamehub/{bin,conf,log}
sudo cp gamehub /opt/gamehub/bin/
sudo cp GameHub.conf.example /opt/gamehub/conf/GameHub.conf

# 編輯配置文件
sudo vim /opt/gamehub/conf/GameHub.conf
# 修改數據庫和 Redis 連接信息

# 創建 systemd 服務
sudo tee /etc/systemd/system/gamehub.service << EOF
[Unit]
Description=GameHub Server
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=gamehub
Group=gamehub
WorkingDirectory=/opt/gamehub
ExecStart=/opt/gamehub/bin/gamehub -conf conf/GameHub.conf -log log
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 創建用戶和設置權限
sudo useradd -r -s /bin/false gamehub
sudo chown -R gamehub:gamehub /opt/gamehub

# 啟動服務
sudo systemctl daemon-reload
sudo systemctl enable gamehub
sudo systemctl start gamehub
```

### Nginx 配置
```bash
# 創建 Nginx 配置
sudo tee /etc/nginx/sites-available/gamehub << EOF
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
    }
    
    location /ws {
        proxy_pass http://localhost:3563;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

# 啟用配置
sudo ln -s /etc/nginx/sites-available/gamehub /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 🏭 生產環境優化

### 資源限制配置
```bash
# 在 docker-compose.yml 中添加資源限制
services:
  gamehub:
    deploy:
      resources:
        limits:
          memory: 4G
          cpus: '2.0'
        reservations:
          memory: 2G
          cpus: '1.0'
    ulimits:
      nofile:
        soft: 65536
        hard: 65536
```

### 數據庫優化
```sql
-- PostgreSQL 性能調優
ALTER SYSTEM SET shared_buffers = '1GB';
ALTER SYSTEM SET effective_cache_size = '3GB';
ALTER SYSTEM SET maintenance_work_mem = '256MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;
ALTER SYSTEM SET random_page_cost = 1.1;
SELECT pg_reload_conf();

-- 創建索引
CREATE INDEX CONCURRENTLY idx_player_id ON game_records(player_id);
CREATE INDEX CONCURRENTLY idx_game_time ON game_records(created_at);
```

### Redis 優化
```bash
# Redis 配置優化
docker exec gamehub-redis redis-cli CONFIG SET maxmemory 2gb
docker exec gamehub-redis redis-cli CONFIG SET maxmemory-policy allkeys-lru
docker exec gamehub-redis redis-cli CONFIG SET save "900 1 300 10 60 10000"
docker exec gamehub-redis redis-cli CONFIG SET tcp-keepalive 300
```

### 系統優化
```bash
# 內核參數調優
cat >> /etc/sysctl.conf << EOF
# 網絡優化
net.core.somaxconn = 32768
net.core.netdev_max_backlog = 32768
net.ipv4.tcp_max_syn_backlog = 32768
net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_keepalive_time = 1200
net.ipv4.tcp_max_tw_buckets = 32768

# 文件描述符限制
fs.file-max = 1000000
EOF

sysctl -p

# 用戶限制
cat >> /etc/security/limits.conf << EOF
* soft nofile 65536
* hard nofile 65536
* soft nproc 32768
* hard nproc 32768
EOF
```

## 📊 監控與運維

### 監控系統部署
```bash
# 啟動監控服務
make monitoring

# 訪問監控界面
echo "監控服務："
echo "• Grafana: http://localhost:3000 (admin/admin123)"
echo "• Prometheus: http://localhost:9090"
```

### 基本監控指標
```bash
# 系統資源監控
make stats

# 應用日誌監控
tail -f logs/gamehub/*.log | grep -E "(ERROR|WARN|FATAL)"

# 數據庫性能監控
docker exec gamehub-postgres psql -U gamehub -c "
SELECT 
    datname,
    numbackends,
    xact_commit,
    xact_rollback,
    blks_read,
    blks_hit
FROM pg_stat_database 
WHERE datname = 'gamehub';"

# Redis 性能監控
docker exec gamehub-redis redis-cli INFO stats | grep -E "(keyspace|memory|stats)"
```

### 自動化備份
```bash
# 創建備份腳本
cat > /opt/backup_gamehub.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/opt/backups/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

# 數據庫備份
docker exec gamehub-postgres pg_dump -U gamehub gamehub > "$BACKUP_DIR/database.sql"

# Redis 備份
docker exec gamehub-redis redis-cli --rdb - > "$BACKUP_DIR/redis.rdb"

# 配置文件備份
cp -r docker/config "$BACKUP_DIR/"

# 清理舊備份 (保留30天)
find /opt/backups -type d -mtime +30 -exec rm -rf {} +

echo "備份完成: $BACKUP_DIR"
EOF

chmod +x /opt/backup_gamehub.sh

# 設置定時備份
echo "0 2 * * * /opt/backup_gamehub.sh" | crontab -
```

### 日誌輪轉
```bash
# 創建 logrotate 配置
cat > /etc/logrotate.d/gamehub << EOF
/opt/gamehub/log/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 gamehub gamehub
    postrotate
        systemctl reload gamehub
    endscript
}
EOF
```

## 🔧 故障排除

### 常見問題診斷

#### 1. 服務無法啟動
```bash
# 檢查服務狀態
make ps
docker-compose logs gamehub

# 檢查端口占用
netstat -tulpn | grep -E "(3563|3564|8080)"

# 檢查配置文件
docker-compose config
```

#### 2. 數據庫連接失敗
```bash
# 測試數據庫連接
docker exec gamehub-postgres pg_isready -U gamehub

# 檢查數據庫日誌
docker logs gamehub-postgres

# 手動連接測試
make db-shell
```

#### 3. Redis 連接問題
```bash
# 測試 Redis 連接
docker exec gamehub-redis redis-cli ping

# 檢查 Redis 配置
docker exec gamehub-redis redis-cli CONFIG GET "*"

# 進入 Redis 命令行
make redis-shell
```

#### 4. 性能問題診斷
```bash
# 容器資源使用
docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"

# 系統資源監控
htop
iotop
netstat -i

# 應用性能分析
docker exec gamehub-server top
docker exec gamehub-server ps aux
```

### 緊急恢復程序

#### 數據恢復
```bash
# 從備份恢復數據庫
make restore BACKUP_DIR=/path/to/backup

# 手動恢復數據庫
docker exec -i gamehub-postgres psql -U gamehub gamehub < backup/database.sql

# 恢復 Redis 數據
docker-compose stop redis
docker cp backup/redis.rdb $(docker-compose ps -q redis):/data/dump.rdb
docker-compose start redis
```

#### 服務恢復
```bash
# 快速重啟
make restart

# 完全重新部署
make stop
make clean
make quick-start

# 回滾到上一個版本
git checkout <previous-tag>
make prod-build
```

## 📚 部署檢查清單

### 部署前檢查
- [ ] 系統需求滿足
- [ ] Docker 和 Docker Compose 已安裝
- [ ] 環境變量配置完成
- [ ] SSL 證書準備就緒
- [ ] 端口規劃確認
- [ ] 數據庫初始化腳本準備

### 部署後驗證
- [ ] 所有容器正常運行
- [ ] 健康檢查通過
- [ ] Web 服務可訪問
- [ ] WebSocket 連接正常
- [ ] 數據庫連接正常
- [ ] Redis 緩存正常
- [ ] 日誌輸出正常
- [ ] 監控系統運行

### 生產環境檢查
- [ ] 資源限制配置
- [ ] 備份策略實施
- [ ] 監控告警設置
- [ ] 日誌輪轉配置
- [ ] 安全加固完成
- [ ] 災難恢復計劃

---

**總結**: 本指南提供了 GameHub 在 Linux 環境下的完整部署方案，推薦使用容器化部署以獲得最佳的一致性和可維護性。生產環境部署時請特別注意安全配置和監控設置。