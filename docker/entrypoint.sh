#!/bin/sh

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查环境变量
check_env() {
    log_info "检查环境变量..."
    
    if [ -z "$GIN_MODE" ]; then
        export GIN_MODE=release
        log_warn "GIN_MODE 未设置，使用默认值: release"
    fi
    
    if [ -z "$TZ" ]; then
        export TZ=Asia/Shanghai
        log_warn "TZ 未设置，使用默认值: Asia/Shanghai"
    fi
}

# 初始化目录
init_dirs() {
    log_info "初始化目录..."
    
    mkdir -p /app/data
    mkdir -p /app/logs
    mkdir -p /app/plugins
    mkdir -p /var/www
    mkdir -p /var/log/nginx
    mkdir -p /etc/letsencrypt
    
    # 设置权限
    chown -R gpanel:gpanel /app
    chmod -R 755 /app
}

# 检查配置文件
check_config() {
    log_info "检查配置文件..."
    
    if [ ! -f "/app/config.yaml" ]; then
        log_warn "配置文件不存在，使用默认配置"
        cp /app/internal/config/config.yaml /app/config.yaml
    fi
    
    # 检查数据库连接
    if ! grep -q "type: sqlite" /app/config.yaml; then
        log_info "检测到非SQLite数据库，等待数据库连接..."
        # 这里可以添加数据库连接检查逻辑
    fi
}

# 启动Nginx
start_nginx() {
    log_info "启动Nginx..."
    
    # 检查Nginx配置
    nginx -t
    
    if [ $? -eq 0 ]; then
        nginx
        log_info "Nginx 启动成功"
    else
        log_error "Nginx 配置错误"
        exit 1
    fi
}

# 启动应用
start_app() {
    log_info "启动 GPanel 应用..."
    
    case "$1" in
        "server")
            ./gpanel
            ;;
        "migrate")
            ./gpanel migrate
            ;;
        "backup")
            ./gpanel backup
            ;;
        *)
            log_error "未知命令: $1"
            log_info "可用命令: server, migrate, backup"
            exit 1
            ;;
    esac
}

# 健康检查
health_check() {
    log_info "执行健康检查..."
    
    # 检查应用是否启动
    for i in 1 2 3 4 5; do
        if wget --no-verbose --tries=1 --spider http://localhost/health; then
            log_info "健康检查通过"
            return 0
        fi
        log_warn "健康检查失败，重试 $i/5"
        sleep 2
    done
    
    log_error "健康检查失败"
    return 1
}

# 信号处理
cleanup() {
    log_info "收到停止信号，正在关闭服务..."
    
    # 停止应用
    if [ -n "$GOLANG_PID" ]; then
        kill -TERM $GOLANG_PID
        wait $GOLANG_PID
    fi
    
    # 停止Nginx
    nginx -s quit
    
    log_info "服务已停止"
    exit 0
}

# 设置信号处理
trap cleanup SIGTERM SIGINT

# 主函数
main() {
    log_info "GPanel 容器启动中..."
    log_info "版本: $(cat /app/VERSION 2>/dev/null || echo 'unknown')"
    log_info "构建时间: $(cat /app/BUILD_TIME 2>/dev/null || echo 'unknown')"
    
    # 检查环境
    check_env
    
    # 初始化目录
    init_dirs
    
    # 检查配置
    check_config
    
    # 启动Nginx
    start_nginx
    
    # 启动应用
    log_info "启动应用服务: $1"
    start_app "$1" &
    GOLANG_PID=$!
    
    # 等待应用启动
    sleep 3
    
    # 健康检查
    if [ "$1" = "server" ]; then
        health_check
    fi
    
    log_info "GPanel 启动完成!"
    
    # 等待进程结束
    wait $GOLANG_PID
}

# 执行主函数
main "$@"