# GameHub Windows + VS Code 容器化開發指南

本指南將幫助您在 Windows 上使用 VS Code 和 Docker 容器進行 GameHub 開發。

## 前置要求

### 必需軟件
- Windows 10/11 (支援 WSL2)
- Docker Desktop for Windows
- VS Code
- Git for Windows

### VS Code 擴展
- Remote - Containers
- Go (golang.go)
- Docker (ms-azuretools.vscode-docker)
- GitLens

## 快速開始

### 1. 初始設置

```batch
# 1. 克隆專案（如果尚未完成）
git clone <your-repo-url>
cd stone

# 2. 運行初始設置腳本
scripts\windows-setup.bat
```

### 2. 啟動開發環境

```batch
# 方法1: 使用便捷腳本
dev-start.bat

# 方法2: 手動命令
docker-compose -f docker-compose.dev.yml up -d
```

### 3. 連接到容器

**選項 A: 使用 Dev Container（推薦）**
1. 在 VS Code 中按 `Ctrl+Shift+P`
2. 輸入 `Remote-Containers: Reopen in Container`
3. 選擇 GameHub Development Environment

**選項 B: 附加到運行中的容器**
1. 在 VS Code 中按 `Ctrl+Shift+P`
2. 輸入 `Remote-Containers: Attach to Running Container`
3. 選擇 `stone-gamehub-dev-1`

## 開發工作流程

### 日常開發

1. **啟動環境**
   ```batch
   dev-start.bat
   ```

2. **打開 VS Code**
   - 腳本會詢問是否自動打開 VS Code
   - 或手動運行: `code .`

3. **連接到容器**
   - VS Code 會自動檢測到 `.devcontainer` 配置
   - 選擇 "Reopen in Container"

4. **開始編碼**
   - 所有 Go 工具已預配置
   - 代碼自動格式化和 lint
   - 支援調試和測試

### 常用 VS Code 任務

在 VS Code 中按 `Ctrl+Shift+P` 然後輸入 `Tasks: Run Task`:

- **start-dependencies**: 啟動數據庫和相關服務
- **stop-dependencies**: 停止所有服務
- **build-gamehub**: 構建遊戲服務器
- **test-all**: 運行所有測試
- **test-coverage**: 運行測試並顯示覆蓋率
- **lint**: 運行代碼檢查
- **docker-build-dev**: 重新構建開發容器

### 調試

#### 本地調試
1. 設置斷點
2. 按 `F5` 或選擇 "Launch GameHub Server"
3. 服務器將在調試模式下啟動

#### 容器調試
1. 確保容器正在運行
2. 選擇 "Attach to GameHub Container"
3. 容器內的 Delve 調試器會處理斷點

#### 複合調試
- 選擇 "Launch Full Stack" 配置
- 自動啟動依賴服務並啟動調試

### 測試

```bash
# 在容器內運行
go test ./...                    # 所有測試
go test -cover ./...             # 包含覆蓋率
go test -v ./internal/game/...   # 特定包的詳細測試
```

或使用 VS Code 任務:
- `Ctrl+Shift+P` → `Tasks: Run Task` → `test-all`

## 服務訪問

開發環境啟動後，您可以訪問：

| 服務 | 地址 | 用途 |
|------|------|------|
| GameHub 服務器 | http://localhost | 主遊戲服務 |
| pgAdmin | http://localhost:5050 | PostgreSQL 管理 |
| Redis Commander | http://localhost:8081 | Redis 監控 |
| API 文檔 | http://localhost:8083 | 服務 API 文檔 |

### 數據庫連接信息

**PostgreSQL (via pgAdmin)**
- 服務器: postgres
- 端口: 5432
- 數據庫: gamehub_dev
- 用戶名: gamehub_user
- 密碼: dev_password_123

**Redis (via Redis Commander)**
- 主機: redis
- 端口: 6379
- 密碼: dev_redis_pass

## 故障排除

### 常見問題

#### 1. Docker 未運行
```
錯誤: Cannot connect to the Docker daemon
解決: 啟動 Docker Desktop for Windows
```

#### 2. 端口衝突
```
錯誤: Port already in use
解決: 
docker-compose -f docker-compose.dev.yml down
netstat -ano | findstr :5432
taskkill /PID <PID> /F
```

#### 3. 容器構建失敗
```bash
# 清理並重新構建
docker-compose -f docker-compose.dev.yml down
docker system prune -f
docker-compose -f docker-compose.dev.yml build --no-cache
```

#### 4. VS Code 無法連接到容器
```
解決步驟:
1. 確認容器正在運行: docker ps
2. 重新安裝 Remote-Containers 擴展
3. 重啟 VS Code
4. 檢查 .devcontainer/devcontainer.json 配置
```

#### 5. Go 模組問題
```bash
# 在容器內執行
cd /app/src/GameHub
go mod tidy
go mod download
```

### 調試命令

```bash
# 查看容器狀態
docker-compose -f docker-compose.dev.yml ps

# 查看服務日誌
docker-compose -f docker-compose.dev.yml logs gamehub-dev

# 進入容器 shell
docker-compose -f docker-compose.dev.yml exec gamehub-dev bash

# 重啟特定服務
docker-compose -f docker-compose.dev.yml restart gamehub-dev
```

## 高級配置

### 自定義環境變量

編輯 `.env` 文件來自定義配置:

```env
# 數據庫設置
POSTGRES_DB=gamehub_dev
POSTGRES_USER=gamehub_user
POSTGRES_PASSWORD=dev_password_123

# Redis 設置
REDIS_PASSWORD=dev_redis_pass

# 應用設置
LOG_LEVEL=debug
DEBUG_MODE=true
```

### VS Code 設置自定義

編輯 `.vscode/settings.json` 來調整開發環境:

```json
{
    "go.toolsManagement.autoUpdate": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "editor.formatOnSave": true
}
```

### 添加新的調試配置

在 `.vscode/launch.json` 中添加:

```json
{
    "name": "Debug Specific Module",
    "type": "go",
    "request": "launch",
    "mode": "debug",
    "program": "${workspaceFolder}/GameHub/internal/yourmodule",
    "args": ["--config", "dev.conf"]
}
```

## 生產部署

當開發完成後，使用生產配置:

```bash
# 構建生產鏡像
docker-compose build

# 啟動生產環境
docker-compose up -d

# 查看生產日誌
docker-compose logs -f gamehub
```

## 備份和恢復

### 數據庫備份
```bash
# 備份
docker-compose exec postgres pg_dump -U gamehub_user gamehub_dev > backup.sql

# 恢復
docker-compose exec -T postgres psql -U gamehub_user gamehub_dev < backup.sql
```

### Redis 備份
```bash
# 備份
docker-compose exec redis redis-cli --rdb dump.rdb

# 恢復會在容器重啟時自動進行
```

## 效能優化建議

1. **資源分配**: 在 Docker Desktop 中為容器分配足夠的 CPU 和記憶體
2. **卷映射**: 使用命名卷而非綁定掛載以提高 I/O 效能
3. **網路**: 使用 Docker 內部網路減少延遲
4. **構建緩存**: 利用多階段構建和層緩存

## 團隊協作

### Git 工作流程
1. 每個功能使用獨立分支
2. 提交前運行完整測試套件
3. 使用 pre-commit hooks 進行代碼檢查
4. 定期同步主分支

### 代碼審查檢查表
- [ ] 代碼格式化 (goimports)
- [ ] 靜態分析通過 (golangci-lint)
- [ ] 所有測試通過
- [ ] 測試覆蓋率 > 80%
- [ ] 文檔更新
- [ ] 安全性檢查

---

如需更多幫助，請查看:
- [Docker 官方文檔](https://docs.docker.com/)
- [VS Code Remote-Containers 文檔](https://code.visualstudio.com/docs/remote/containers)
- [Go 開發最佳實踐](https://golang.org/doc/effective_go)