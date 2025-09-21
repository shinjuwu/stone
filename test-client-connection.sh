#!/bin/bash

# GameHub 客戶端連接測試腳本
# 用於驗證部署的服務器是否可以正常連接

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

print_test() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

# 獲取本機 IP
get_local_ip() {
    local ip=""
    if command -v hostname &> /dev/null; then
        ip=$(hostname -I 2>/dev/null | awk '{print $1}')
    fi
    if [ -z "$ip" ]; then
        ip="localhost"
    fi
    echo "$ip"
}

# 測試端口連通性
test_port() {
    local host=$1
    local port=$2
    local service=$3
    
    print_test "測試 $service ($host:$port)"
    
    if timeout 5 bash -c "cat < /dev/null > /dev/tcp/$host/$port" 2>/dev/null; then
        print_info "✓ $service 端口連通"
        return 0
    else
        print_error "✗ $service 端口無法連接"
        return 1
    fi
}

# 測試 HTTP API
test_http_api() {
    local host=$1
    local port=$2
    
    print_test "測試 HTTP API ($host:$port)"
    
    local response=$(curl -s -w "%{http_code}" -o /dev/null http://$host:$port/health 2>/dev/null || echo "000")
    
    if [ "$response" = "200" ]; then
        print_info "✓ HTTP API 正常響應"
        return 0
    else
        print_error "✗ HTTP API 響應異常 (HTTP $response)"
        return 1
    fi
}

# 測試 WebSocket 連接
test_websocket() {
    local host=$1
    local port=$2
    
    print_test "測試 WebSocket 連接 ($host:$port)"
    
    # 檢查是否安裝了 wscat
    if ! command -v wscat &> /dev/null; then
        print_warn "wscat 未安裝，跳過 WebSocket 測試"
        print_info "安裝命令: npm install -g wscat"
        return 0
    fi
    
    # 使用 timeout 限制測試時間
    if timeout 10 wscat -c ws://$host:$port --no-check 2>/dev/null <<< '{"test":"connection"}' | grep -q "Connected"; then
        print_info "✓ WebSocket 連接成功"
        return 0
    else
        print_error "✗ WebSocket 連接失敗"
        return 1
    fi
}

# 測試數據庫連接
test_database() {
    print_test "測試數據庫連接"
    
    if docker exec gamehub-postgres-client-dev pg_isready -U gamehub_dev &>/dev/null; then
        print_info "✓ 數據庫連接正常"
        return 0
    else
        print_error "✗ 數據庫連接失敗"
        return 1
    fi
}

# 測試 Redis 連接
test_redis() {
    print_test "測試 Redis 連接"
    
    if docker exec gamehub-redis-client-dev redis-cli ping 2>/dev/null | grep -q "PONG"; then
        print_info "✓ Redis 連接正常"
        return 0
    else
        print_error "✗ Redis 連接失敗"
        return 1
    fi
}

# 檢查容器狀態
check_containers() {
    print_test "檢查容器狀態"
    
    local containers=("gamehub-postgres-client-dev" "gamehub-redis-client-dev" "gamehub-server-client-dev")
    local all_running=true
    
    for container in "${containers[@]}"; do
        if docker ps --filter "name=$container" --filter "status=running" | grep -q "$container"; then
            print_info "✓ $container 運行中"
        else
            print_error "✗ $container 未運行"
            all_running=false
        fi
    done
    
    if $all_running; then
        return 0
    else
        return 1
    fi
}

# 生成測試報告
generate_report() {
    local host=$1
    local total_tests=$2
    local passed_tests=$3
    
    echo ""
    echo "=========================================="
    echo "📊 連接測試報告"
    echo "=========================================="
    echo "測試時間: $(date)"
    echo "服務器地址: $host"
    echo "通過測試: $passed_tests/$total_tests"
    echo ""
    
    if [ $passed_tests -eq $total_tests ]; then
        print_info "🎉 所有測試通過！客戶端可以正常連接服務器。"
        echo ""
        echo "客戶端連接信息："
        echo "  WebSocket: ws://$host:3563"
        echo "  TCP:       $host:3564"
        echo "  HTTP API:  http://$host:8080"
    else
        print_warn "⚠️  部分測試失敗，請檢查服務器配置。"
        echo ""
        echo "故障排除建議："
        echo "1. 檢查服務器狀態: ./deploy-client-dev.sh status"
        echo "2. 查看服務器日誌: ./deploy-client-dev.sh logs"
        echo "3. 重啟服務器: ./deploy-client-dev.sh restart"
    fi
    echo ""
}

# 主測試函數
main() {
    local host=$(get_local_ip)
    local total_tests=0
    local passed_tests=0
    
    echo "=========================================="
    echo "🧪 GameHub 客戶端連接測試"
    echo "=========================================="
    echo "測試服務器: $host"
    echo ""
    
    # 檢查容器狀態
    ((total_tests++))
    if check_containers; then
        ((passed_tests++))
    fi
    
    echo ""
    
    # 測試各個端口
    ((total_tests++))
    if test_port $host 8080 "HTTP API"; then
        ((passed_tests++))
    fi
    
    ((total_tests++))
    if test_port $host 3563 "WebSocket"; then
        ((passed_tests++))
    fi
    
    ((total_tests++))
    if test_port $host 3564 "TCP"; then
        ((passed_tests++))
    fi
    
    echo ""
    
    # 測試 HTTP API 響應
    ((total_tests++))
    if test_http_api $host 8080; then
        ((passed_tests++))
    fi
    
    echo ""
    
    # 測試數據庫和 Redis
    ((total_tests++))
    if test_database; then
        ((passed_tests++))
    fi
    
    ((total_tests++))
    if test_redis; then
        ((passed_tests++))
    fi
    
    # 生成報告
    generate_report $host $total_tests $passed_tests
    
    # 返回測試結果
    if [ $passed_tests -eq $total_tests ]; then
        exit 0
    else
        exit 1
    fi
}

# 執行測試
main "$@"