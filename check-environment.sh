#!/bin/bash

# GameHub 環境檢測腳本
# 檢測當前環境並提供解決方案

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

print_success() {
    echo -e "${CYAN}[SUCCESS]${NC} $1"
}

# 檢測操作系統
detect_os() {
    if [[ -f /proc/version ]] && grep -q Microsoft /proc/version; then
        echo "WSL"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "Linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macOS"
    elif [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "msys" ]]; then
        echo "Windows"
    else
        echo "Unknown"
    fi
}

# 檢查 Docker 狀態
check_docker() {
    local docker_available=false
    local docker_running=false
    
    if command -v docker &> /dev/null; then
        docker_available=true
        print_info "✓ Docker 命令可用"
        
        if docker info &> /dev/null; then
            docker_running=true
            print_success "✓ Docker 服務運行正常"
        else
            print_warn "⚠ Docker 已安裝但服務未運行"
        fi
    else
        print_error "✗ Docker 未安裝或不可用"
    fi
    
    echo "$docker_available:$docker_running"
}

# 檢查 Docker Compose
check_docker_compose() {
    if command -v docker-compose &> /dev/null; then
        print_info "✓ Docker Compose 可用"
        return 0
    else
        print_error "✗ Docker Compose 未安裝"
        return 1
    fi
}

# 檢查端口占用
check_key_ports() {
    local ports=(3563 3564 8080 5432 6379 6380 6381)
    local occupied_ports=()
    
    print_step "檢查關鍵端口..."
    
    for port in "${ports[@]}"; do
        if ss -tuln 2>/dev/null | grep -q ":$port " || netstat -tuln 2>/dev/null | grep -q ":$port "; then
            occupied_ports+=($port)
        fi
    done
    
    if [ ${#occupied_ports[@]} -gt 0 ]; then
        print_warn "以下端口已被占用: ${occupied_ports[*]}"
        print_info "可以使用 ./configure-ports.sh auto 自動解決"
    else
        print_success "✓ 所有關鍵端口都可用"
    fi
}

# 提供解決方案
provide_solution() {
    local os=$1
    local docker_status=$2
    
    local docker_available=$(echo $docker_status | cut -d: -f1)
    local docker_running=$(echo $docker_status | cut -d: -f2)
    
    echo ""
    echo "=========================================="
    echo "🔧 解決方案建議"
    echo "=========================================="
    
    case $os in
        "WSL")
            print_step "WSL 環境解決方案："
            if [ "$docker_available" = "false" ]; then
                echo "1. 安裝 Docker Desktop for Windows"
                echo "   下載地址: https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe"
                echo ""
                echo "2. 配置 WSL 整合："
                echo "   - 打開 Docker Desktop"
                echo "   - Settings → Resources → WSL Integration" 
                echo "   - 啟用 WSL 整合"
                echo ""
                echo "3. 或者直接在 WSL 中安裝 Docker："
                echo "   curl -fsSL https://get.docker.com -o get-docker.sh"
                echo "   sudo sh get-docker.sh"
            elif [ "$docker_running" = "false" ]; then
                echo "Docker 已安裝但未運行，請："
                echo "1. 如果使用 Docker Desktop，請啟動 Docker Desktop"
                echo "2. 如果是 WSL 原生安裝，運行: sudo service docker start"
            fi
            ;;
        "Linux")
            print_step "Linux 環境解決方案："
            if [ "$docker_available" = "false" ]; then
                echo "安裝 Docker："
                echo "curl -fsSL https://get.docker.com -o get-docker.sh"
                echo "sudo sh get-docker.sh"
                echo "sudo usermod -aG docker \$USER"
                echo ""
                echo "安裝 Docker Compose："
                echo "sudo curl -L \"https://github.com/docker/compose/releases/latest/download/docker-compose-\$(uname -s)-\$(uname -m)\" -o /usr/local/bin/docker-compose"
                echo "sudo chmod +x /usr/local/bin/docker-compose"
            elif [ "$docker_running" = "false" ]; then
                echo "啟動 Docker 服務："
                echo "sudo systemctl start docker"
                echo "sudo systemctl enable docker"
            fi
            ;;
        "macOS")
            print_step "macOS 環境解決方案："
            if [ "$docker_available" = "false" ]; then
                echo "安裝 Docker Desktop for Mac："
                echo "https://desktop.docker.com/mac/main/amd64/Docker.dmg"
            elif [ "$docker_running" = "false" ]; then
                echo "請啟動 Docker Desktop 應用程序"
            fi
            ;;
        "Windows")
            print_step "Windows 環境解決方案："
            echo "建議使用 Docker Desktop for Windows"
            echo "下載地址: https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe"
            ;;
        *)
            print_warn "未識別的操作系統，請手動安裝 Docker"
            ;;
    esac
}

# 顯示快速啟動指南
show_quick_start() {
    echo ""
    echo "=========================================="
    echo "🚀 Docker 配置完成後的快速啟動"
    echo "=========================================="
    echo ""
    echo "1. 檢查環境："
    echo "   bash check-environment.sh"
    echo ""
    echo "2. 檢查端口衝突："
    echo "   bash configure-ports.sh check"
    echo ""
    echo "3. 自動解決端口衝突（如有）："
    echo "   bash configure-ports.sh auto"
    echo ""
    echo "4. 部署客戶端開發服務器："
    echo "   bash deploy-client-dev.sh"
    echo ""
    echo "5. 測試連接："
    echo "   bash test-client-connection.sh"
    echo ""
}

# 主函數
main() {
    echo "=========================================="
    echo "🔍 GameHub 環境檢測"
    echo "=========================================="
    
    # 檢測操作系統
    local os=$(detect_os)
    print_step "檢測到操作系統: $os"
    
    # 檢查 Docker
    print_step "檢查 Docker 環境..."
    local docker_status=$(check_docker)
    
    # 檢查 Docker Compose
    print_step "檢查 Docker Compose..."
    check_docker_compose
    
    # 檢查端口
    check_key_ports
    
    # 檢查項目文件
    print_step "檢查項目文件..."
    local required_files=("docker-compose.client-dev.yml" "GameHub/Dockerfile.fixed" "docker/config/GameHub.client-dev.conf")
    local missing_files=()
    
    for file in "${required_files[@]}"; do
        if [[ ! -f "$file" ]]; then
            missing_files+=("$file")
        fi
    done
    
    if [ ${#missing_files[@]} -gt 0 ]; then
        print_error "缺少必要文件: ${missing_files[*]}"
    else
        print_success "✓ 所有必要文件都存在"
    fi
    
    # 提供解決方案
    provide_solution "$os" "$docker_status"
    
    # 快速啟動指南
    show_quick_start
    
    # 總結
    echo "=========================================="
    echo "📋 環境檢測總結"
    echo "=========================================="
    
    local docker_available=$(echo $docker_status | cut -d: -f1)
    local docker_running=$(echo $docker_status | cut -d: -f2)
    
    if [ "$docker_available" = "true" ] && [ "$docker_running" = "true" ]; then
        print_success "🎉 環境已就緒，可以開始部署！"
        echo ""
        echo "運行以下命令開始："
        echo "bash deploy-client-dev.sh"
    else
        print_warn "❌ 環境未就緒，請按照上述指南配置 Docker"
    fi
}

# 執行主函數
main "$@"