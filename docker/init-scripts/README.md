# GameHub 數據庫初始化腳本

此目錄包含 GameHub 數據庫的初始化腳本，當 PostgreSQL 容器首次啟動時會自動執行。

## 📁 腳本執行順序

PostgreSQL 會按照檔案名的字母順序執行這些腳本：

### 1. `00-setup-database.sql`
- 數據庫基本設置
- 創建必要的擴展 (uuid-ossp, pgcrypto)
- 設置時區
- 創建通用函數

### 2. `01-init-database.sql` 
- 創建核心系統表結構
- members (用戶表)
- game_records (遊戲記錄表)  
- agents (代理商表)
- game_settings (遊戲設置表)
- 創建基本索引和觸發器

### 3. `02-gamelist.sql`
- 創建 gamelist 表
- 插入所有可用遊戲的基本信息
- 包含遊戲 ID、遊戲代碼、狀態等

### 4. `03-gameinfo.sql`
- 創建 gameinfo 表
- 插入代理商與遊戲的對應關係
- 配置哪些代理商可以提供哪些遊戲

### 5. `04-lobbyinfo.sql`
- 創建 lobbyinfo 表  
- 插入大廳桌台配置信息
- 配置代理商、遊戲和桌台的三方關係

### 6. `99-finalize-initialization.sql`
- 完成初始化的最後步驟
- 更新數據庫統計資訊
- 創建摘要視圖
- 輸出初始化完成報告

## 🎯 主要數據表說明

### gamelist
所有可用遊戲的主列表
- `game_id`: 遊戲唯一 ID
- `game_code`: 遊戲代碼 (如: 'baccarat', 'blackjack')
- `status`: 遊戲狀態 (1=啟用, 0=禁用)

### gameinfo  
代理商遊戲配置表
- `agent_id`: 代理商 ID
- `game_id`: 遊戲 ID (關聯到 gamelist)
- `game_code`: 遊戲代碼
- `status`: 此配置的狀態

### lobbyinfo
大廳桌台配置表
- `agent_id`: 代理商 ID
- `game_id`: 遊戲 ID
- `table_id`: 桌台 ID
- `status`: 桌台狀態

## 🔧 使用方式

### 自動初始化 (推薦)
當你使用 docker-compose 啟動開發環境時，這些腳本會自動執行：

```bash
# 啟動開發環境 (會自動執行初始化)
docker-compose -f docker-compose.dev.yml up -d postgres

# 查看初始化日誌
docker-compose -f docker-compose.dev.yml logs postgres
```

### 手動重新初始化
如果需要重新初始化數據庫：

```bash
# 停止並移除數據
docker-compose -f docker-compose.dev.yml down -v

# 重新啟動 (會重新執行初始化腳本)
docker-compose -f docker-compose.dev.yml up -d postgres
```

### 驗證初始化結果
連接到數據庫並檢查：

```sql
-- 檢查初始化摘要
SELECT * FROM v_database_summary;

-- 檢查遊戲列表
SELECT game_id, game_code FROM gamelist WHERE status = 1 ORDER BY game_id;

-- 檢查代理商配置
SELECT DISTINCT agent_id FROM gameinfo ORDER BY agent_id;
```

## ⚠️ 注意事項

1. **數據持久化**: 初始化後的數據會持久化在 Docker volume 中
2. **僅首次執行**: 這些腳本只在數據庫首次創建時執行
3. **修改腳本**: 如果修改了腳本，需要重新創建容器和數據卷
4. **備份重要**: 生產環境部署前請先備份現有數據

## 🐛 故障排除

### 查看初始化日誌
```bash
docker-compose -f docker-compose.dev.yml logs postgres | grep -E "(ERROR|initialization)"
```

### 常見問題
- **權限錯誤**: 確保腳本檔案有正確的讀取權限
- **語法錯誤**: 檢查 SQL 語法，特別是字符編碼問題
- **依賴順序**: 確保腳本按正確順序執行

### 重置數據庫
```bash
# 完全重置 (警告：會刪除所有數據)
docker-compose -f docker-compose.dev.yml down -v
docker volume rm stone_postgres_dev_data
docker-compose -f docker-compose.dev.yml up -d postgres
```