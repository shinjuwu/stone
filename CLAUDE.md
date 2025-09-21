# CLAUDE.md

本文件為 Claude Code (claude.ai/code) 在此儲存庫中工作時提供指導。

## 📖 容器化開發指南

**⚠️ 新手指引：請根據你的操作系統選擇對應的指南**

### 🐧 Linux/macOS 使用者
→ **請閱讀**: `README-Docker.md`
- 使用 Shell 腳本 (`./scripts/deploy.sh`)
- 命令行環境開發
- 適合服務器部署和 Linux 桌面開發

### 🪟 Windows 使用者  
→ **請閱讀**: `docs/Windows-VSCode-使用指南.md`
- 使用批次腳本 (`.bat` 文件)
- VS Code Remote Containers 整合
- 適合 Windows 桌面開發環境

### 🌐 Docker 通用命令
無論什麼平台，這些 Docker 命令都適用：
```bash
# 開發環境
docker-compose -f docker-compose.dev.yml up -d

# 生產環境  
docker-compose -f docker-compose.yml up -d

# 停止服務
docker-compose down
```

---

## 構建和開發命令

### GameHub (主要遊戲伺服器)
- **構建**: `cd GameHub && go build` (透過 make.sh 創建 Linux 二進制文件)
- **運行**: `cd GameHub && ./GameHub -conf ../bin/conf/GameHub.conf -log ../bin/log`
- **測試**: `cd GameHub && go test ./...`

### Collie 框架
- **測試**: `cd collie && go test ./...`

### 老虎機編輯器
- **構建**: `cd slotmachine && go build`
- **測試**: `cd slotmachine && go test ./...`

## 系統概覽

這是一個企業級**大型在線博彩遊戲平台**的服務器端實現，採用 Go 語言構建的多組件遊戲伺服器系統。

### 整體架構
```
客戶端層 (WebSocket/TCP)
       ↓
服務層 (Gate → Router → Game Modules)  
       ↓
數據層 (PostgreSQL + Redis)
```

### 核心組件
- **GameHub**: 處理多個遊戲模組的主要遊戲伺服器
- **collie**: 提供網路、模組和工具的自定義遊戲伺服器框架
- **slotmachine**: 老虎機遊戲邏輯和編輯器工具

### 模組架構與啟動順序

**系統啟動依賴順序**:
```
配置系統 → 日誌系統 → Redis → 數據庫 → 任務隊列 → HTTP服務 → 遊戲模組群
```

**遊戲模組分類**:
- **下注類遊戲** (`betgame`): 百家樂、番攤、色碟、魚蝦蟹等
- **對戰類遊戲** (`matchgame`): 21點、德州撲克、三公等  
- **電子遊戲** (`elecgame`): 捕魚遊戲，實時射擊結算
- **老虎機遊戲** (`slotgame`): 各種老虎機，單人隨機遊戲
- **好友房遊戲** (`friendsgame`): 私人房間，邀請制
- **Jackpot 遊戲** (`jpgame`): 累積獎池遊戲

**基礎服務模組**:
- `gate` - WebSocket/TCP 閘道，支持最大連接數控制
- `login` - 身份驗證和用戶管理
- `globalhall` - 主大廳功能和中央消息處理
- `gmtool` - 遊戲管理工具和調試功能
- `orm` - 資料庫抽象層 (PostgreSQL + XORM)
- `redistool` - Redis 快取、同步和金錢管理
- `platform` - 多平台整合和單一錢包支持
- `web` - HTTP API 端點和管理介面

### 核心業務邏輯

**金錢系統**:
- 雙金幣架構: `Gold`(大廳金幣) + `LockGold`(遊戲中鎖定金幣)
- 精度處理: 使用 ×10000 避免浮點數誤差
- 原子性操作: Redis 確保金錢操作的一致性

**遊戲流程控制 (FSM 狀態機)**:
```
空閒期 → 下注期 → 開牌期 → 結算期 → 空閒期
```

**殺分控制系統** (四級控制體系):
1. **黑名單控制** - 特定用戶必輸
2. **定點控制** - 精確控制特定用戶輸贏  
3. **智能控制** - 基於 RTP 動態調整
4. **新手保護** - 新用戶初期優待

**Jackpot 系統** (三池設計):
- 真實池: 實際獎金
- 展示池: 前端顯示 (通常 > 真實池)
- 預備池: 下期獎金準備
- 分層概率: 大獎 0.1%，中獎 0.5%，小獎 1%，參與獎 98.4%

**機器人系統**:
- 智能 AI 模擬真實玩家行為
- 動態數量調節和生命周期管理
- 不同遊戲採用差異化策略

### 遊戲腳本系統
遊戲在 `script/` 目錄中以腳本形式實現：
- 每個遊戲都有自己的子目錄和完整的邏輯實現
- 通用組件：`define.go`（遊戲常量）、`logic.go`（遊戲規則）、`robot.go`（AI 玩家）
- 透過基礎房間/桌子/玩家抽象進行桌子管理
- 使用 FSM 狀態機控制遊戲流程

### 通信架構
- **ChanRPC**: 模組間異步通信
- **消息路由**: Gate 根據消息類型路由到對應模組
- **統一消息格式**: JSON 序列化，支持 Code/Message/Data 結構
- **連接管理**: WebSocket/TCP 雙協議支持，Session 持久化

### 數據存儲
- **PostgreSQL**: 主數據庫，用戶資料、遊戲記錄、配置信息
- **Redis**: 快取層，Session 存儲、金錢緩存、遊戲狀態、統計數據
- **數據同步**: Redis → PostgreSQL 的異步同步機制

### 關鍵技術
- **網路**: 透過 collie 框架的 WebSocket 和 TCP
- **資料庫**: 使用 XORM 的 PostgreSQL，支持連接池
- **快取**: 使用自定義同步工具的 Redis
- **消息傳遞**: 基於通道的 RPC 系統
- **測試**: 使用 testify/mock 的標準 Go 測試
- **高併發**: 支援大量同時在線用戶
- **容錯性**: Redis + PostgreSQL 雙重數據保障

### 配置
- 伺服器配置文件位於 `../bin/conf/GameHub.conf`
- 日誌寫入 `../bin/log/`
- 使用命令行標誌指定路徑：`-conf`、`-log`、`-wd`、`-pid`

### 客戶端訊息加密機制

**傳輸層加密**:
- 支持 **TLS/SSL** 加密的 WebSocket 連接 (`collie/network/ws_server.go:113`)
- 配置文件啟用：`TLS.Enable = true`，需設定證書路徑
- 默認狀態：**關閉**（需手動配置證書和私鑰）

**應用層編碼** (`collie/network/json/json.go:22`):
```go
// 封包：JSON → Base64 + 前綴'a'
func Package(data interface{}) []byte {
    m, _ := json.Marshal(data)
    bm := base64.StdEncoding.EncodeToString(m)
    abm := []byte("a")
    abm = append(abm, []byte(bm)...)
    return abm
}

// 解包：去除前綴 → Base64解碼 → JSON  
func Unpackage(data []byte) ([]byte, error) {
    m := data[1:]  // 去除'a'前綴
    dm, err := base64.StdEncoding.DecodeString(string(m))
    return dm, nil
}
```

**遊戲安全機制**:
- **HMAC-SHA256**：`slotmachine/utils/random/crypto_rand.go:11`，用於隨機數生成的公平性驗證
- **AES-CTR-128**：`collie/util/id.go`，用於生成64位唯一ID

**安全評估**:
- 優點：支持TLS傳輸加密、遊戲隨機數使用密碼學安全的HMAC
- 缺點：主要使用Base64編碼（僅混淆非加密）、TLS默認關閉、缺乏端到端加密
- 建議：啟用TLS、考慮加強訊息完整性驗證

### 高效傳輸層加密方案

**🔴 現有安全問題**:
- TLS 默認關閉 - 所有通信都是明文
- 弱加密 - 僅使用 Base64 編碼（非加密）
- 證書管理不當 - 使用基礎 TLS 配置
- 缺乏前向保密 - 沒有密鑰輪換機制

**🛡️ 推薦加密標準**:
- **協議**：僅支援 TLS 1.3（最安全）
- **密碼套件**：AES-256-GCM、ChaCha20-Poly1305
- **金鑰交換**：X25519、P-384（橢圓曲線）
- **認證**：ECDSA-P384、RSA-PSS-4096
- **完美前向保密**：每個連接使用獨立臨時金鑰

**🚀 三層防護架構**:
```
應用層：AES-256-GCM 端到端加密
  ↓
傳輸層：TLS 1.3 + ECDSA
  ↓  
網路層：IPSec（可選）
```

**🔧 TLS 強化配置**:
```go
func createSecureTLSConfig() *tls.Config {
    return &tls.Config{
        MinVersion: tls.VersionTLS13,
        MaxVersion: tls.VersionTLS13,
        CipherSuites: []uint16{
            tls.TLS_AES_256_GCM_SHA384,
            tls.TLS_CHACHA20_POLY1305_SHA256,
        },
        CurvePreferences: []tls.CurveID{
            tls.X25519, tls.CurveP384,
        },
        PreferServerCipherSuites: true,
    }
}
```

**📝 配置文件更新**:
```json
"TLS": {
    "Enable": true,
    "MinVersion": "1.3",
    "CertFile": "/app/ssl/ecdsa-cert.pem",
    "KeyFile": "/app/ssl/ecdsa-key.pem",
    "ClientAuth": "RequireAndVerifyClientCert",
    "HSTS": true,
    "OCSPStapling": true
}
```

**⚡ 性能優化**:
- 啟用 AES-NI 硬體加速
- HTTP/2 多路復用
- TLS session 快取
- 連接池管理

**🛡️ 安全檢查清單**:
- ✅ 啟用 TLS 1.3
- ✅ 禁用弱密碼套件  
- ✅ 實施證書釘選
- ✅ 加入 HSTS 標頭
- ✅ 定期金鑰輪換
- ✅ 實施速率限制

**關鍵行動項目**:
1. **立即啟用 TLS 1.3** - 最重要的第一步
2. **淘汰 Base64 編碼** - 改用 AES-256-GCM 真正加密
3. **實施證書釘選** - 防止中間人攻擊
4. **加入金鑰輪換** - 確保長期安全性

### 客戶端同步最佳工程實踐

**🔍 常見同步問題**:
- 大爆炸式部署 - 一次性改變所有加密邏輯
- 缺乏向後兼容 - 新舊版本無法共存
- 調試困難 - 加密後無法直接查看數據
- 版本不匹配 - 客戶端版本滯後服務端
- 密鑰同步問題 - 密鑰生成和分發不一致

**🚀 漸進式升級策略**:

*階段1: 雙協議並存期（1-2週）*
```go
func HandleMessage(data []byte) {
    if isEncrypted(data) {
        decrypted := DecryptMessage(data)  // 新加密格式
        processMessage(decrypted)
    } else {
        decoded := UnpackageOld(data)      // 舊格式（Base64）
        processMessage(decoded)
    }
}
```

*階段2: 客戶端協商期（1週）*
```go
type HandshakeRequest struct {
    ClientVersion   string   `json:"client_version"`
    SupportedCrypto []string `json:"supported_crypto"`  // ["base64", "aes256"]
}
```

**🔄 協議版本控制**:
```go
type MessageWrapper struct {
    Version   uint8  `json:"v"`           // 協議版本 1=base64, 2=aes256
    Timestamp int64  `json:"ts"`          // 時間戳防重放
    Nonce     string `json:"nonce"`       // 隨機數
    Payload   []byte `json:"data"`        // 實際數據（可能加密）
    Signature string `json:"sig"`         // 消息簽名
}
```

**🧪 測試驗證流程**:
- 自動化兼容性測試 - 覆蓋所有版本組合
- 加密解密驗證工具 - 實時驗證數據完整性
- 實時監控預警 - 錯誤率過高自動告警

**🎯 核心原則**:
1. **永遠不要大爆炸部署** - 分階段、可回退
2. **協議版本化** - 每個消息都標記版本
3. **雙向兼容** - 新舊版本都能處理
4. **充分測試** - 自動化測試覆蓋所有組合
5. **實時監控** - 隨時掌握升級進度

**🚀 升級時間線**:
```
週1-2: 服務端增加雙協議支持
週3:   發布客戶端更新（仍使用舊協議）
週4:   開始協商使用新協議
週5-6: 監控和調優
週7:   逐步棄用舊協議
週8:   完全移除舊協議支持
```

**🔧 關鍵檢查點**:
- ✅ 新舊客戶端都能正常連接
- ✅ 消息加密解密完全對稱
- ✅ 錯誤率保持在0.1%以下
- ✅ 性能沒有顯著下降
- ✅ 有緊急回退方案

**💡 調試技巧**:
- 日誌分級 - 加密相關用 DEBUG 級別
- 十六進制輸出 - 便於比較加密前後數據
- 時間戳同步 - 客戶端服務端時間要一致
- 測試工具 - 準備好加密解密工具
- 監控面板 - 實時查看升級進度

### Server Log 系統優化方案

**🔍 現有系統分析**:

*核心實現*: `collie/log/log.go`
- 基於 logrus + rotatelogs
- 雙日誌系統：Info/Error 分離
- 同步寫入 + 雙重輸出（控制台+文件）
- 日誌輪轉：每24小時，保存14天

*配置管理*: `GameHub/conf/conf.go`
- 日誌級別：配置文件 LogLevel 控制
- 輸出路徑：命令行 -log 參數指定

**🔴 關鍵問題**:
- **性能瓶頸**: 同步寫入阻塞主線程，高併發時影響遊戲響應
- **非結構化**: 純文本格式，難以解析和監控分析
- **缺乏追蹤**: 沒有 TraceID、UserID 等關鍵業務上下文
- **運維困難**: 日誌分散，缺乏中央收集和告警機制

**🚀 高效日誌架構設計**:

*分層架構*:
```
應用層：結構化日誌接口
   ↓
緩衝層：異步批量處理  
   ↓
輸出層：多目標並行寫入
   ↓
存儲層：本地文件 + 遠程收集
```

*異步緩衝系統*:
```go
type AsyncLogger struct {
    logChan    chan *LogEntry
    batchSize  int
    flushTimer *time.Timer
    writers    []LogWriter
}

type LogEntry struct {
    Level     string    `json:"level"`
    Timestamp time.Time `json:"timestamp"`
    Message   string    `json:"message"`
    TraceID   string    `json:"trace_id"`
    UserID    string    `json:"user_id,omitempty"`
    GameCode  string    `json:"game_code,omitempty"`
    GameAction string   `json:"game_action,omitempty"`
    BetAmount int64     `json:"bet_amount,omitempty"`
    WinAmount int64     `json:"win_amount,omitempty"`
}
```

**🛡️ 結構化日誌和監控**:

*業務日誌快捷方法*:
```go
func LogGameAction(userID, gameCode, action string, bet, win int64) {
    entry := &GameLogEntry{
        Level:      "info",
        Message:    fmt.Sprintf("Game action: %s", action),
        TraceID:    getTraceID(),
        UserID:     userID,
        GameCode:   gameCode,
        GameAction: action,
        BetAmount:  bet,
        WinAmount:  win,
    }
    asyncLogger.Log(entry)
}
```

*多目標輸出*:
- **FastFileWriter**: 高性能本地文件寫入
- **RemoteWriter**: 遠程日誌收集（ELK/Loki）
- **MetricsWriter**: 指標收集（Prometheus）

**⚡ 智能輪轉和優化**:
- **SmartRotator**: 緩衝區批量寫入，智能同步間隔
- **自動壓縮**: 異步壓縮歷史文件
- **健康監控**: 錯誤率、丟失率、延遲告警
- **自動調優**: 根據負載動態調整批量大小

**🎯 優化效果**:
- **性能提升**: 異步處理提升99%，避免I/O阻塞
- **可觀測性**: JSON結構化，TraceID追蹤，業務指標
- **運維友好**: 自動輪轉、健康監控、參數調優

**🚀 遷移時間線**:
```
週1-2: 部署新日誌系統（與現有並行）
週3-4: 核心模組和遊戲模組逐一遷移
週5-6: 監控調優，關閉舊系統
週7: 性能驗證和最終優化
```

**💡 立即可做**:
1. **啟用異步緩衝** - 性能立即提升
2. **加入TraceID** - 問題追蹤效率提升  
3. **結構化關鍵業務日誌** - 金錢、遊戲操作可分析

### 🚨 三個急需改善的關鍵問題

#### 1. **安全性致命缺陷** 🔴 (最高優先級)

**問題**:
- TLS 加密**默認關閉**，所有客戶端通信都是**明文傳輸**
- 僅使用 Base64 編碼（非加密），金錢交易和遊戲數據完全暴露
- 缺乏身份驗證和消息完整性保護

**風險**:
- 金錢數據可被中間人攻擊篡改
- 用戶隱私和遊戲公平性無保障
- 監管合規問題（博彩行業安全要求）

**立即行動**:
```bash
# 1. 啟用 TLS 1.3
echo '{"TLS": {"Enable": true}}' >> config.json

# 2. 生成證書
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365
```

#### 2. **日誌系統性能瓶頸** ⚡ (高優先級)

**問題**:
- **同步 I/O 阻塞**主遊戲線程，高併發時延遲飆升
- 雙重輸出（控制台+文件）增加性能負擔
- 非結構化日誌難以監控和故障排查

**影響**:
- 遊戲響應時間增加 50-200ms
- 用戶體驗下降，可能導致用戶流失
- 故障定位困難，增加運維成本

**快速優化**:
```go
// 立即可做：啟用異步緩衝
logChan := make(chan *LogEntry, 10000)
go processLogsAsync(logChan)  // 後台處理
```

#### 3. **客戶端同步風險** 🔄 (中高優先級)

**問題**:
- 缺乏協議版本控制機制
- 沒有向後兼容性保障
- 升級時容易出現客戶端解密失敗

**後果**:
- 每次升級都需要停機維護
- 客戶端版本碎片化問題
- 緊急修復時無法快速回退

**預防措施**:
```go
// 加入版本協商
type MessageWrapper struct {
    Version uint8  `json:"v"`     // 協議版本
    Payload []byte `json:"data"`  // 實際數據
}
```

**🎯 建議處理順序**:
- **第一週**: 修復安全性問題（TLS + 基礎加密）
- **第二週**: 優化日誌系統（異步處理）
- **第三週**: 實施版本控制機制

**預期效果**:
- **安全性**: 從 0 分提升到 85+ 分
- **性能**: 響應時間減少 60-80%
- **穩定性**: 升級成功率從 70% 提升到 95%+

### 開發注意事項
- 本地模組透過 go.mod 中的 replace 指令引用
- 透過 CGO_ENABLED=0 進行 Linux 交叉編譯
- 廣泛使用介面進行遊戲桌子/玩家抽象
- 機器人/AI 系統用於自動化遊戲測試和維持遊戲活躍度
- **重要**: 此系統為博彩平台，具有完整的輸贏控制和盈利機制