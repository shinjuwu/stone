#!/bin/bash

# GameHub 客戶端開發服務器端口配置腳本
# 用於手動配置端口以避免衝突

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

# 檢查端口是否可用
is_port_available() {
    local port=$1
    if ss -tuln 2>/dev/null | grep -q ":$port " || netstat -tuln 2>/dev/null | grep -q ":$port "; then
        return 1  # 端口被占用
    else
        return 0  # 端口可用
    fi
}

# 建議可用端口
suggest_port() {
    local base_port=$1
    local max_attempts=20
    
    for ((i=0; i<max_attempts; i++)); do
        local test_port=$((base_port + i))
        if is_port_available $test_port; then
            echo $test_port
            return 0
        fi
    done
    
    return 1
}

# 顯示當前配置
show_current_config() {
    print_step "當前端口配置："
    
    if [ -f "docker-compose.client-dev.yml" ]; then
        local ws_port=$(grep -A 10 "gamehub:" docker-compose.client-dev.yml | grep "3563:" | awk -F: '{print $2}' | tr -d '"')
        local tcp_port=$(grep -A 10 "gamehub:" docker-compose.client-dev.yml | grep "3564:" | awk -F: '{print $2}' | tr -d '"')
        local http_port=$(grep -A 10 "gamehub:" docker-compose.client-dev.yml | grep "8080:" | awk -F: '{print $2}' | tr -d '"')
        local db_port=$(grep -A 10 "postgres:" docker-compose.client-dev.yml | grep "5432:" | awk -F: '{print $2}' | tr -d '"')
        local redis_port=$(grep -A 10 "redis:" docker-compose.client-dev.yml | grep "6379\"" | awk -F: '{print $2}' | tr -d '"')
        
        echo "  WebSocket: ${ws_port:-3563}"
        echo "  TCP:       ${tcp_port:-3564}"
        echo "  HTTP:      ${http_port:-8080}"
        echo "  PostgreSQL: ${db_port:-5432}"
        echo "  Redis:     ${redis_port:-6381}"
    else
        print_error "docker-compose.client-dev.yml 文件不存在"
        exit 1
    fi
}

# 檢查端口衝突
check_port_conflicts() {
    print_step "檢查端口衝突..."
    
    local ports=(3563 3564 8080 5432)
    local redis_port=$(grep -A 10 "redis:" docker-compose.client-dev.yml | grep "6379\"" | awk -F: '{print $2}' | tr -d '"')
    if [ -n "$redis_port" ]; then
        ports+=($redis_port)
    fi
    
    local conflicts=()
    
    for port in "${ports[@]}"; do
        if ! is_port_available $port; then
            conflicts+=($port)
        fi
    done
    
    if [ ${#conflicts[@]} -gt 0 ]; then
        print_warn "以下端口存在衝突: ${conflicts[*]}"
        return 1
    else
        print_info "所有端口都可用 ✓"
        return 0
    fi
}

# 自動解決端口衝突
auto_resolve_conflicts() {
    print_step "自動解決端口衝突..."
    
    local current_redis_port=$(grep -A 10 "redis:" docker-compose.client-dev.yml | grep "6379\"" | awk -F: '{print $2}' | tr -d '"')
    
    if [ -n "$current_redis_port" ] && ! is_port_available $current_redis_port; then
        print_info "Redis 端口 $current_redis_port 被占用，尋找替代端口..."
        
        local new_redis_port=$(suggest_port 6381)
        if [ -n "$new_redis_port" ]; then
            update_redis_port $new_redis_port
            print_info "Redis 端口已更新為: $new_redis_port"
        else
            print_error "無法找到可用的 Redis 端口"
            return 1
        fi
    fi
    
    # 檢查其他關鍵端口
    local ports_services=("5432:PostgreSQL" "8080:HTTP" "3563:WebSocket" "3564:TCP")
    
    for port_service in "${ports_services[@]}"; do
        local port=$(echo $port_service | cut -d: -f1)
        local service=$(echo $port_service | cut -d: -f2)
        
        if ! is_port_available $port; then
            print_warn "$service 端口 $port 被占用"
            local suggested_port=$(suggest_port $((port + 10)))
            if [ -n "$suggested_port" ]; then
                print_info "建議 $service 使用端口: $suggested_port"
                read -p "是否更新 $service 端口到 $suggested_port？(y/N): " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    update_service_port $service $port $suggested_port
                fi
            fi
        fi
    done
}

# 更新 Redis 端口
update_redis_port() {
    local new_port=$1
    sed -i "s/\"[0-9]*:6379\"/\"$new_port:6379\"/" docker-compose.client-dev.yml
}

# 更新服務端口
update_service_port() {
    local service=$1
    local old_port=$2
    local new_port=$3
    
    case $service in
        "PostgreSQL")
            sed -i "s/\"$old_port:5432\"/\"$new_port:5432\"/" docker-compose.client-dev.yml
            ;;
        "HTTP")
            sed -i "s/\"$old_port:8080\"/\"$new_port:8080\"/" docker-compose.client-dev.yml
            ;;
        "WebSocket")
            sed -i "s/\"$old_port:3563\"/\"$new_port:3563\"/" docker-compose.client-dev.yml
            ;;
        "TCP")
            sed -i "s/\"$old_port:3564\"/\"$new_port:3564\"/" docker-compose.client-dev.yml
            ;;
    esac
    print_info "$service 端口已更新為: $new_port"
}

# 手動配置端口
manual_configure() {
    print_step "手動配置端口..."
    
    echo "請輸入要使用的端口 (直接回車使用建議值)："
    
    # Redis 端口配置
    local suggested_redis=$(suggest_port 6381)
    read -p "Redis 端口 (建議: $suggested_redis): " redis_port
    redis_port=${redis_port:-$suggested_redis}
    
    if is_port_available $redis_port; then
        update_redis_port $redis_port
        print_info "Redis 端口設置為: $redis_port"
    else
        print_error "端口 $redis_port 不可用"
        return 1
    fi
    
    # 其他端口可以類似配置...
    print_info "端口配置完成"
}

# 重置為默認端口
reset_to_defaults() {
    print_step "重置為默認端口配置..."
    
    # 備份當前配置
    cp docker-compose.client-dev.yml docker-compose.client-dev.yml.backup
    
    # 重置端口
    sed -i 's/"[0-9]*:5432"/"5432:5432"/' docker-compose.client-dev.yml
    sed -i 's/"[0-9]*:6379"/"6381:6379"/' docker-compose.client-dev.yml
    sed -i 's/"[0-9]*:3563"/"3563:3563"/' docker-compose.client-dev.yml
    sed -i 's/"[0-9]*:3564"/"3564:3564"/' docker-compose.client-dev.yml
    sed -i 's/"[0-9]*:8080"/"8080:8080"/' docker-compose.client-dev.yml
    
    print_info "已重置為默認端口配置"
    print_info "備份文件: docker-compose.client-dev.yml.backup"
}

# 顯示幫助
show_help() {
    echo "GameHub 端口配置工具"
    echo ""
    echo "用法: $0 [選項]"
    echo ""
    echo "選項:"
    echo "  check      檢查端口衝突"
    echo "  auto       自動解決衝突"
    echo "  manual     手動配置端口"
    echo "  reset      重置為默認端口"
    echo "  show       顯示當前配置"
    echo "  help       顯示此幫助"
    echo ""
}

# 主函數
main() {
    local command=${1:-check}
    
    echo "=========================================="
    echo "🔧 GameHub 端口配置工具"
    echo "=========================================="
    
    case $command in
        "check")
            show_current_config
            echo ""
            if check_port_conflicts; then
                print_info "🎉 沒有端口衝突，可以正常部署"
            else
                print_warn "❌ 存在端口衝突，建議運行 'auto' 自動解決"
            fi
            ;;
        "auto")
            show_current_config
            echo ""
            auto_resolve_conflicts
            echo ""
            print_info "✅ 端口衝突已自動解決"
            show_current_config
            ;;
        "manual")
            manual_configure
            ;;
        "reset")
            reset_to_defaults
            show_current_config
            ;;
        "show")
            show_current_config
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