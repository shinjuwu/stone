# 🚀 GameHub 客戶端開發服務器 - 快速參考

## 一鍵命令

```bash
# 🚀 啟動服務器
./deploy-client-dev.sh

# 🔧 檢查端口衝突
./configure-ports.sh check

# 🔧 自動解決端口衝突
./configure-ports.sh auto

# 🧪 測試連接
./test-client-connection.sh

# 📊 查看狀態
./deploy-client-dev.sh status

# 📋 查看日誌
./deploy-client-dev.sh logs

# 🔄 重啟服務
./deploy-client-dev.sh restart

# ⛔ 停止服務
./deploy-client-dev.sh stop
```

## 連接信息

| 服務 | 端口 | 協議 | 用途 |
|------|------|------|------|
| WebSocket | 3563 | ws:// | 遊戲客戶端主要連接 |
| TCP | 3564 | tcp:// | 備用TCP連接 |
| HTTP API | 8080 | http:// | REST API和健康檢查 |

## 客戶端配置模板

```json
{
    "server": "YOUR_SERVER_IP",
    "websocket_port": 3563,
    "tcp_port": 3564,
    "http_port": 8080,
    "use_ssl": false,
    "timeout": 30000
}
```

## 快速測試

```bash
# 測試 HTTP API
curl http://YOUR_IP:8080/health

# 測試端口連通性
nc -zv YOUR_IP 3563
nc -zv YOUR_IP 3564

# 使用 dytRobot 測試
./plinko_optimized
```

## 故障排除

| 問題 | 解決方案 |
|------|----------|
| 端口被占用 | `./configure-ports.sh auto` |
| 容器未啟動 | `./deploy-client-dev.sh rebuild` |
| 連接被拒絕 | 檢查防火牆: `sudo ufw allow PORT` |
| 服務異常 | `./deploy-client-dev.sh logs` |
| Redis 端口衝突 | `./configure-ports.sh check` |

## 端口配置

```bash
# 檢查當前端口配置
./configure-ports.sh show

# 檢查端口衝突
./configure-ports.sh check

# 自動解決衝突
./configure-ports.sh auto

# 手動配置端口
./configure-ports.sh manual

# 重置為默認端口
./configure-ports.sh reset
```

## 支持的遊戲

- ✅ Plinko (彈球遊戲)
- ✅ 老虎機系列
- ✅ 捕魚遊戲
- ✅ 21點
- ✅ 德州撲克
- ✅ 三公
- ✅ Jackpot 遊戲

---
📝 詳細文檔: [CLIENT_DEV_SETUP.md](CLIENT_DEV_SETUP.md)