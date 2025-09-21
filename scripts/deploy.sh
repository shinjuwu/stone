#!/bin/bash

# GameHub 容器化部署腳本
# 使用方法: ./scripts/deploy.sh [dev|staging|prod]

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日誌函數
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 檢查依賴
check_dependencies() {
    log_step "檢查依賴工具..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安裝，請先安裝 Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安裝，請先安裝 Docker Compose"
        exit 1
    fi
    
    log_info "依賴檢查完成"
}

# 環境設置
setup_environment() {
    local env=${1:-dev}
    log_step "設置 $env 環境..."
    
    # 複製環境變量文件
    if [[ ! -f ".env" ]]; then
        if [[ -f ".env.example" ]]; then
            cp .env.example .env
            log_info "已創建 .env 文件，請根據需要修改配置"
        else
            log_error ".env.example 文件不存在"
            exit 1
        fi
    fi
    
    # 創建必要的目錄
    mkdir -p docker/nginx/ssl
    mkdir -p data/postgres
    mkdir -p data/redis
    mkdir -p logs/nginx
    mkdir -p logs/gamehub
    
    log_info "環境設置完成"
}

# 構建映像
build_images() {
    log_step "構建 Docker 映像..."
    
    # 構建 GameHub 主服務
    log_info "構建 GameHub 主服務映像..."
    docker-compose build gamehub
    
    # 如果有老虎機編輯器，也構建它
    if [[ -f "slotmachine/Dockerfile.editor" ]]; then
        log_info "構建老虎機編輯器映像..."
        docker-compose build slot-editor
    fi
    
    log_info "映像構建完成"
}

# 啟動服務
start_services() {
    local env=${1:-dev}
    log_step "啟動服務..."
    
    case $env in
        "dev")
            # 開發環境：啟動所有服務包括工具
            docker-compose --profile tools up -d
            ;;
        "staging"|"prod")
            # 生產環境：只啟動核心服務
            docker-compose up -d postgres redis gamehub nginx
            ;;
        *)
            log_error "未知環境: $env"
            exit 1
            ;;
    esac
    
    log_info "服務啟動完成"
}

# 健康檢查
health_check() {
    log_step "執行健康檢查..."
    
    local max_attempts=30
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        log_info "檢查服務狀態 (第 $attempt/$max_attempts 次)..."
        
        # 檢查容器狀態
        if docker-compose ps | grep -q "Up"; then
            # 檢查服務是否響應
            if curl -f http://localhost/health > /dev/null 2>&1; then
                log_info "所有服務運行正常！"
                return 0
            fi
        fi
        
        sleep 10
        ((attempt++))
    done
    
    log_error "健康檢查失敗"
    docker-compose logs
    return 1
}

# 停止服務
stop_services() {
    log_step "停止服務..."
    docker-compose down
    log_info "服務已停止"
}

# 清理資源
cleanup() {
    log_step "清理資源..."
    docker-compose down -v --remove-orphans
    docker system prune -f
    log_info "清理完成"
}

# 備份數據
backup_data() {
    log_step "備份數據..."
    
    local backup_dir="backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$backup_dir"
    
    # 備份數據庫
    docker-compose exec -T postgres pg_dump -U gamehub gamehub > "$backup_dir/database.sql"
    
    # 備份 Redis 數據
    docker-compose exec -T redis redis-cli --rdb - > "$backup_dir/redis.rdb"
    
    # 備份配置文件
    cp -r docker/config "$backup_dir/"
    
    log_info "備份完成: $backup_dir"
}

# 還原數據
restore_data() {
    local backup_dir=$1
    
    if [[ -z "$backup_dir" ]]; then
        log_error "請指定備份目錄"
        exit 1
    fi
    
    if [[ ! -d "$backup_dir" ]]; then
        log_error "備份目錄不存在: $backup_dir"
        exit 1
    fi
    
    log_step "從 $backup_dir 還原數據..."
    
    # 還原數據庫
    if [[ -f "$backup_dir/database.sql" ]]; then
        docker-compose exec -T postgres psql -U gamehub gamehub < "$backup_dir/database.sql"
    fi
    
    # 還原 Redis 數據
    if [[ -f "$backup_dir/redis.rdb" ]]; then
        docker-compose stop redis
        docker cp "$backup_dir/redis.rdb" $(docker-compose ps -q redis):/data/dump.rdb
        docker-compose start redis
    fi
    
    log_info "數據還原完成"
}

# 查看日誌
view_logs() {
    local service=${1:-}
    
    if [[ -n "$service" ]]; then
        docker-compose logs -f "$service"
    else
        docker-compose logs -f
    fi
}

# 主函數
main() {
    local command=${1:-deploy}
    local env=${2:-dev}
    
    echo "======================================"
    echo "GameHub 容器化部署腳本"
    echo "======================================"
    
    case $command in
        "deploy")
            check_dependencies
            setup_environment "$env"
            build_images
            start_services "$env"
            health_check
            ;;
        "start")
            start_services "$env"
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            stop_services
            start_services "$env"
            ;;
        "build")
            build_images
            ;;
        "logs")
            view_logs "$env"
            ;;
        "backup")
            backup_data
            ;;
        "restore")
            restore_data "$env"
            ;;
        "cleanup")
            cleanup
            ;;
        "health")
            health_check
            ;;
        *)
            echo "使用方法: $0 {deploy|start|stop|restart|build|logs|backup|restore|cleanup|health} [dev|staging|prod]"
            echo ""
            echo "命令說明:"
            echo "  deploy  - 完整部署（默認）"
            echo "  start   - 啟動服務"
            echo "  stop    - 停止服務"
            echo "  restart - 重啟服務"
            echo "  build   - 構建映像"
            echo "  logs    - 查看日誌"
            echo "  backup  - 備份數據"
            echo "  restore - 還原數據"
            echo "  cleanup - 清理所有資源"
            echo "  health  - 健康檢查"
            exit 1
            ;;
    esac
}

# 執行主函數
main "$@"