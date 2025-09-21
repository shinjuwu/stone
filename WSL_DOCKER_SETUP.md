# 🐳 WSL + Docker 環境設置指南

## 🔍 問題分析
你當前在 WSL 環境中，需要正確配置 Docker 才能運行客戶端開發服務器。

## 🚀 解決方案

### 方法一：使用 Docker Desktop (推薦)

#### 1. 安裝 Docker Desktop
- 下載並安裝 [Docker Desktop for Windows](https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe)
- 確保安裝時勾選 "Use WSL 2 instead of Hyper-V"

#### 2. 配置 WSL 整合
1. 打開 Docker Desktop
2. 進入 Settings → Resources → WSL Integration
3. 勾選 "Enable integration with my default WSL distro"
4. 在 "Enable integration with additional distros" 中選擇你的 WSL 發行版
5. 點擊 "Apply & Restart"

#### 3. 驗證安裝
```bash
# 在 WSL 中執行
docker --version
docker-compose --version
```

### 方法二：直接在 WSL 中安裝 Docker

#### 1. 更新套件列表
```bash
sudo apt update
```

#### 2. 安裝 Docker
```bash
# 安裝依賴
sudo apt install -y apt-transport-https ca-certificates curl gnupg lsb-release

# 添加 Docker 官方 GPG 鑰匙
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# 添加 Docker 存儲庫
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 安裝 Docker
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io
```

#### 3. 安裝 Docker Compose
```bash
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 4. 配置用戶權限
```bash
sudo usermod -aG docker $USER
newgrp docker
```

#### 5. 啟動 Docker 服務
```bash
sudo service docker start
```

### 方法三：在 Windows 主機上運行

如果 WSL 配置困難，可以直接在 Windows 上運行：

#### 1. 複製項目到 Windows
```bash
# 將項目複製到 Windows 目錄
cp -r /mnt/c/work/stone /mnt/c/GameHub-Windows/
```

#### 2. 在 Windows PowerShell 中運行
```powershell
# 打開 PowerShell 並切換到項目目錄
cd C:\GameHub-Windows

# 運行部署腳本 (使用 Git Bash 或 PowerShell)
.\deploy-client-dev.sh  # 在 Git Bash 中
# 或者
bash deploy-client-dev.sh  # 如果有 bash 環境
```

## 🔧 快速驗證

運行以下命令檢查 Docker 是否可用：

```bash
# 檢查 Docker
docker --version

# 檢查 Docker Compose  
docker-compose --version

# 測試 Docker 運行
docker run hello-world

# 檢查 Docker 服務狀態
docker info
```

## 🚀 Docker 配置完成後

一旦 Docker 可用，就可以運行部署腳本：

```bash
# 回到項目目錄
cd /mnt/c/work/stone

# 運行部署腳本
bash deploy-client-dev.sh

# 或者使用端口配置工具
bash configure-ports.sh check
bash deploy-client-dev.sh
```

## 🎯 常見問題解決

### Docker Desktop 相關
- **問題**: Docker Desktop 啟動失敗
- **解決**: 確保 Windows 功能中的 "Windows Subsystem for Linux" 和 "Virtual Machine Platform" 已啟用

### WSL 整合問題
- **問題**: WSL 中看不到 Docker 命令
- **解決**: 檢查 Docker Desktop 的 WSL 整合設置

### 權限問題
- **問題**: 提示需要 sudo 權限
- **解決**: 確保用戶已添加到 docker 組：`sudo usermod -aG docker $USER`

### 服務啟動問題
- **問題**: Docker 服務未運行
- **解決**: 
  - WSL: `sudo service docker start`
  - Windows: 啟動 Docker Desktop

## 📝 推薦配置

我推薦使用 **Docker Desktop + WSL 整合** 的方案，因為：

1. ✅ 圖形化管理界面
2. ✅ 自動更新
3. ✅ 完整的 WSL 支持
4. ✅ 穩定性較好
5. ✅ 官方支持

## 🔄 配置完成後的下一步

1. 驗證 Docker 可用性
2. 運行 `bash configure-ports.sh check` 檢查端口
3. 運行 `bash deploy-client-dev.sh` 部署服務器
4. 運行 `bash test-client-connection.sh` 測試連接

---

**注意**: 配置完成後，記得重新啟動 WSL 或重新登錄以確保環境變量生效。