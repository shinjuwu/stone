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

# 檢查端口占用
check_ports() {
    print_step "檢查端口占用..."
    
    local ports=(3563 3564 8080 5432 6379)
    local occupied_ports=()
    
    for port in "${ports[@]}"; do
        if netstat -tuln 2>/dev/null | grep -q ":$port " || ss -tuln 2>/dev/null | grep -q ":$port "; then
            occupied_ports+=($port)
        fi
    done
    
    if [ ${#occupied_ports[@]} -gt 0 ]; then
        print_warn "以下端口已被占用: ${occupied_ports[*]}"
        print_warn "這可能會導致服務啟動失敗"
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

# 顯示幫助信息
show_help() {
    echo "GameHub 客戶端開發服務器部署腳本"
    echo ""
    echo "用法: $0 [選項]"
    echo ""
    echo "選項:"
    echo "  start          啟動服務（默認）"
    echo "  stop           停止服務"
    echo "  restart        重啟服務"
    echo "  logs           查看日誌"
    echo "  status         查看服務狀態"
    echo "  clean          清理服務和數據"
    echo "  rebuild        重新構建並啟動"
    echo "  help           顯示此幫助信息"
    echo ""
    echo "示例:"
    echo "  $0              # 啟動服務"
    echo "  $0 start        # 啟動服務"
    echo "  $0 stop         # 停止服務"
    echo "  $0 logs         # 查看日誌"
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

# 清理服務
clean_services() {
    print_warn "這將刪除所有容器和數據，無法恢復！"
    read -p "確定要清理所有數據嗎？(y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_step "清理服務和數據..."
        docker-compose -f docker-compose.client-dev.yml down -v --remove-orphans
        docker system prune -f
        print_success "清理完成"
    else
        print_info "清理已取消"
    fi
}

# 重新構建
rebuild_services() {
    print_step "重新構建並啟動服務..."
    docker-compose -f docker-compose.client-dev.yml down
    docker-compose -f docker-compose.client-dev.yml build --no-cache
    start_services
    wait_for_services
    show_connection_info
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
        "clean")
            clean_services
            ;;
        "rebuild")
            rebuild_services
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "未知命令: $command"
            show_help
            exit 1
            ;;
    esac
}

# 執行主函數
main "$@"