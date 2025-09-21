# 🎮 GameHub 客戶端開發服務器

專為客戶端開發者設計的一鍵部署方案，無需 TLS 配置，直接使用 IP 連線。

## 🚀 快速開始

### 一鍵部署
```bash
# 啟動服務器
./deploy-client-dev.sh

# 或者
./deploy-client-dev.sh start
```

### 系統需求
- Docker 20.10+
- Docker Compose 2.0+
- 可用端口：3563, 3564, 8080, 5432, 6379

## 📡 連接信息

部署成功後，客戶端可以使用以下信息連接：

```
WebSocket: ws://YOUR_IP:3563
TCP:       YOUR_IP:3564
HTTP API:  http://YOUR_IP:8080
```

## 🎯 客戶端配置示例

### dytRobot 測試工具配置
```json
{
    "serverIP": "YOUR_SERVER_IP",
    "wsPort": 3563,
    "tcpPort": 3564,
    "httpPort": 8080,
    "ssl": false
}
```

### WebSocket 連接示例
```javascript
// JavaScript 客戶端
const ws = new WebSocket('ws://YOUR_SERVER_IP:3563');

ws.onopen = function() {
    console.log('連接成功');
    // 發送登錄消息
    ws.send(JSON.stringify({
        "Login": {
            "Account": "test_user",
            "Passwd": "123456"
        }
    }));
};

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('收到消息:', data);
};
```

### TCP 連接示例
```python
# Python 客戶端
import socket
import json

def connect_tcp():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect(('YOUR_SERVER_IP', 3564))
    
    # 發送登錄消息
    login_msg = json.dumps({
        "Login": {
            "Account": "test_user", 
            "Passwd": "123456"
        }
    })
    
    sock.send(login_msg.encode())
    response = sock.recv(1024)
    print('收到回應:', response.decode())
    
    return sock
```

## 🛠️ 管理命令

```bash
# 查看服務狀態
./deploy-client-dev.sh status

# 查看實時日誌
./deploy-client-dev.sh logs

# 重啟服務
./deploy-client-dev.sh restart

# 停止服務
./deploy-client-dev.sh stop

# 重新構建（代碼更新後）
./deploy-client-dev.sh rebuild

# 清理所有數據
./deploy-client-dev.sh clean
```

## 🎲 支持的遊戲

服務器支持以下遊戲類型：
- **電子遊戲**: Plinko, 捕魚, 老虎機
- **對戰遊戲**: 21點, 德州撲克, 三公
- **好友房**: 私人房間遊戲
- **Jackpot**: 累積獎池遊戲

## 🧪 測試流程

### 1. 基本連接測試
```bash
# 測試 HTTP API
curl http://YOUR_SERVER_IP:8080/health

# 測試 WebSocket（使用 wscat）
npm install -g wscat
wscat -c ws://YOUR_SERVER_IP:3563
```

### 2. 遊戲功能測試
```bash
# 使用 dytRobot 測試工具
cd dytRobot
./dytRobot.exe

# 或使用命令行版本
./plinko_optimized
```

### 3. 性能測試
```bash
# 查看服務器資源使用
docker stats

# 查看日誌中的性能指標
./deploy-client-dev.sh logs | grep -E "(latency|memory|cpu)"
```

## 🔧 配置說明

### 服務器配置特點
- **TLS**: 關閉，直接明文連接
- **端口**: 標準端口，無需額外映射
- **日誌級別**: Debug (4)，詳細輸出
- **最大連接數**: 1000（適合開發測試）
- **數據庫**: 獨立的開發數據庫
- **Redis**: 無密碼，簡化配置

### 環境變量
```bash
# 數據庫
DB_NAME=gamehub_client_dev
DB_USER=gamehub_dev
DB_PASSWORD=dev123

# Redis
REDIS_PASSWORD=""  # 無密碼

# 服務
PLATFORM=DEV
LOG_LEVEL=4
SERVER_ID=1
```

## 🚨 故障排除

### 常見問題

#### 1. 端口被占用
```bash
# 檢查端口使用
netstat -tulpn | grep -E "(3563|3564|8080|5432|6379)"

# 關閉占用端口的進程
sudo lsof -t -i:3563 | xargs sudo kill -9
```

#### 2. Docker 服務未啟動
```bash
# 啟動 Docker
sudo systemctl start docker

# 檢查 Docker 狀態
sudo systemctl status docker
```

#### 3. 容器啟動失敗
```bash
# 查看詳細錯誤
./deploy-client-dev.sh logs

# 重新構建
./deploy-client-dev.sh rebuild
```

#### 4. 連接被拒絕
```bash
# 檢查防火牆
sudo ufw status
sudo ufw allow 3563
sudo ufw allow 3564
sudo ufw allow 8080

# 檢查服務器狀態
./deploy-client-dev.sh status
curl http://localhost:8080/health
```

## 📝 開發注意事項

### 消息格式
所有客戶端消息都使用 JSON 格式，並經過 Base64 編碼：
```javascript
// 原始消息
const message = {"Login": {"Account": "test", "Passwd": "123"}};

// 編碼後發送
const encoded = 'a' + btoa(JSON.stringify(message));
websocket.send(encoded);
```

### 認證流程
1. 發送 Login 消息
2. 接收 LoginResponse
3. 發送 JoinRoom 消息
4. 開始遊戲交互

### 遊戲流程
1. 選擇遊戲房間
2. 發送 PlayerAction 消息
3. 接收遊戲結果
4. 處理獎勵和積分

## 📚 相關文檔

- [Plinko 遊戲測試指南](USAGE_GUIDE.md)
- [消息協議說明](GameHub/msg/)
- [完整部署指南](LINUX_DEPLOYMENT_GUIDE.md)

---

**技術支持**: 如有問題請檢查日誌或聯繫開發團隊