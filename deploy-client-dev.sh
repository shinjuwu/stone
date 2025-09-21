#!/bin/bash

# GameHub 客戶端開發服務器一鍵部署腳本
# 用途：為客戶端開發者提供簡單的測試服務器
# 特點：無需 TLS，直接 IP 連線，快速啟動

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 打印帶顏色的消息
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

# 獲取本機 IP 地址
get_local_ip() {
    # 嘗試不同方法獲取本機 IP
    local ip=""
    
    # 方法1：使用 hostname 命令
    if command -v hostname &> /dev/null; then
        ip=$(hostname -I 2>/dev/null | awk '{print $1}')
    fi
    
    # 方法2：使用 ip 命令
    if [ -z "$ip" ] && command -v ip &> /dev/null; then
        ip=$(ip route get 8.8.8.8 2>/dev/null | grep -oP 'src \K\S+')
    fi
    
    # 方法3：使用 ifconfig 命令
    if [ -z "$ip" ] && command -v ifconfig &> /dev/null; then
        ip=$(ifconfig 2>/dev/null | grep -E "inet [0-9]" | grep -v "127.0.0.1" | head -1 | awk '{print $2}')
    fi
    
    # 默認回退
    if [ -z "$ip" ]; then
        ip="localhost"
    fi
    
    echo "$ip"
}

# 檢查依賴
check_dependencies() {
    print_step "檢查系統依賴..."
    
    # 檢查 Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安裝，請先安裝 Docker"
        echo "安裝命令：curl -fsSL https://get.docker.com -o get-docker.sh && sudo sh get-docker.sh"
        exit 1
    fi
    
    # 檢查 Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose 未安裝，請先安裝 Docker Compose"
        echo "安裝命令：sudo curl -L \"https://github.com/docker/compose/releases/latest/download/docker-compose-\$(uname -s)-\$(uname -m)\" -o /usr/local/bin/docker-compose && sudo chmod +x /usr/local/bin/docker-compose"
        exit 1
    fi
    
    # 檢查 Docker 服務
    if ! docker info &> /dev/null; then
        print_error "Docker 服務未運行，請啟動 Docker 服務"
        echo "啟動命令：sudo systemctl start docker"
        exit 1
    fi
    
    print_info "依賴檢查完成 ✓"
}

# 建議替代的 Redis 端口
suggest_alternative_redis_port() {
    local alternative_ports=(6382 6383 6384 6385 6386)
    local available_port=""
    
    print_step "尋找可用的 Redis 端口..."
    
    for port in "${alternative_ports[@]}"; do
        if ! ss -tuln 2>/dev/null | grep -q ":$port " && ! netstat -tuln 2>/dev/null | grep -q ":$port "; then
            available_port=$port
            break
        fi
    done
    
    if [ -n "$available_port" ]; then
        print_info "找到可用端口: $available_port"
        print_step "更新 Redis 配置到端口 $available_port..."
        
        # 更新 docker-compose 文件中的端口
        sed -i "s/\"6381:6379\"/\"$available_port:6379\"/" docker-compose.client-dev.yml
        
        print_success "Redis 端口已更新為: $available_port"
    else
        print_error "無法找到可用的 Redis 端口"
        exit 1
    fi
}

# 檢查端口占用
check_ports() {
    print_step "檢查端口占用..."
    
    local ports=(3563 3564 8080 5432 6381)  # 改為檢查 6381
    local occupied_ports=()
    
    for port in "${ports[@]}"; do
        if ss -tuln 2>/dev/null | grep -q ":$port " || netstat -tuln 2>/dev/null | grep -q ":$port "; then
            occupied_ports+=($port)
        fi
    done
    
    if [ ${#occupied_ports[@]} -gt 0 ]; then
        print_warn "以下端口已被占用: ${occupied_ports[*]}"
        print_warn "這可能會導致服務啟動失敗"
        
        # 如果是 Redis 端口被占用，提供解決方案
        for port in "${occupied_ports[@]}"; do
            if [ "$port" = "6381" ]; then
                print_info "Redis 端口 6381 被占用，嘗試使用其他端口..."
                suggest_alternative_redis_port
                return 0
            fi
        done
        
        read -p "是否繼續部署？(y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "部署已取消"
            exit 0
        fi
    else
        print_info "端口檢查完成 ✓"
    fi
}

# 構建服務
build_services() {
    print_step "構建 GameHub 服務鏡像..."
    
    if ! docker-compose -f docker-compose.client-dev.yml build; then
        print_error "服務構建失敗"
        exit 1
    fi
    
    print_info "服務構建完成 ✓"
}

# 啟動服務
start_services() {
    print_step "啟動客戶端開發服務器..."
    
    # 停止可能存在的舊容器
    docker-compose -f docker-compose.client-dev.yml down 2>/dev/null || true
    
    # 啟動服務
    if ! docker-compose -f docker-compose.client-dev.yml up -d; then
        print_error "服務啟動失敗"
        exit 1
    fi
    
    print_info "服務啟動完成 ✓"
}

# 等待服務就緒
wait_for_services() {
    print_step "等待服務啟動完成..."
    
    local max_attempts=60
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        print_info "檢查服務狀態 ($attempt/$max_attempts)..."
        
        # 檢查容器狀態
        local running_containers=$(docker-compose -f docker-compose.client-dev.yml ps --services --filter "status=running" | wc -l)
        local total_containers=$(docker-compose -f docker-compose.client-dev.yml ps --services | wc -l)
        
        if [[ $running_containers -eq $total_containers ]] && [[ $total_containers -gt 0 ]]; then
            # 檢查 GameHub 服務是否響應
            if curl -f http://localhost:8080/health &>/dev/null; then
                print_success "所有服務啟動完成！"
                return 0
            fi
        fi
        
        sleep 5
        ((attempt++))
    done
    
    print_error "服務啟動超時，請檢查日誌"
    docker-compose -f docker-compose.client-dev.yml logs
    return 1
}

# 顯示連接信息
show_connection_info() {
    local local_ip=$(get_local_ip)
    local redis_port=$(docker port gamehub-redis-client-dev 6379/tcp 2>/dev/null | cut -d: -f2)
    
    echo ""
    echo "=========================================="
    echo "🎮 GameHub 客戶端開發服務器部署完成！"
    echo "=========================================="
    echo ""
    echo "📍 服務器連接信息："
    echo "   WebSocket: ws://$local_ip:3563"
    echo "   TCP:       $local_ip:3564"
    echo "   HTTP API:  http://$local_ip:8080"
    echo ""
    echo "🔧 管理界面："
    echo "   服務器狀態: http://$local_ip:8080/health"
    if [ -n "$redis_port" ]; then
        echo "   Redis 端口: $redis_port (避免衝突自動選擇)"
    fi
    echo ""
    echo "📊 服務狀態："
    docker-compose -f docker-compose.client-dev.yml ps
    echo ""
    echo "💡 常用命令："
    echo "   查看日誌: docker-compose -f docker-compose.client-dev.yml logs -f"
    echo "   停止服務: docker-compose -f docker-compose.client-dev.yml down"
    echo "   重啟服務: docker-compose -f docker-compose.client-dev.yml restart"
    echo ""
    echo "🎯 客戶端配置示例："
    echo "   服務器地址: $local_ip"
    echo "   WebSocket端口: 3563"
    echo "   TCP端口: 3564"
    echo ""
    print_success "部署完成！客戶端現在可以連接到服務器進行開發測試。"
}

# 停止服務
stop_services() {
    print_step "停止客戶端開發服務器..."
    docker-compose -f docker-compose.client-dev.yml down
    print_success "服務已停止"
}

# 查看日誌
show_logs() {
    docker-compose -f docker-compose.client-dev.yml logs -f --tail=100
}

# 查看狀態
show_status() {
    local local_ip=$(get_local_ip)
    
    echo "服務狀態："
    docker-compose -f docker-compose.client-dev.yml ps
    echo ""
    echo "連接信息："
    echo "WebSocket: ws://$local_ip:3563"
    echo "TCP: $local_ip:3564"
    echo "HTTP API: http://$local_ip:8080"
}

# 主函數
main() {
    local command=${1:-start}
    
    case $command in
        "start")
            echo "=========================================="
            echo "🚀 GameHub 客戶端開發服務器部署"
            echo "=========================================="
            check_dependencies
            check_ports
            build_services
            start_services
            wait_for_services
            show_connection_info
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            stop_services
            sleep 2
            start_services
            wait_for_services
            show_connection_info
            ;;
        "logs")
            show_logs
            ;;
        "status")
            show_status
            ;;
        *)
            echo "使用方法: $0 {start|stop|restart|logs|status}"
            exit 1
            ;;
    esac
}

# 執行主函數
main "$@"