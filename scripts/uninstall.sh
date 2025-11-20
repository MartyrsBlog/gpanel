#!/bin/bash

# GPanel 一键卸载脚本
# 支持 Ubuntu/Debian/CentOS/Arch Linux

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置变量
INSTALL_DIR="/opt/gpanel"
SERVICE_NAME="gpanel"
USER="gpanel"
REMOVE_DOCKER="false"
REMOVE_NGINX="false"
REMOVE_DATA="false"
BACKUP_DATA="true"

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

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# 显示banner
show_banner() {
    echo -e "${BLUE}"
    cat << 'EOF'
 ____             _                    
|  _ \           | |                   
| |_) | __ _  ___| | _____ _ __   __ _ 
|  _ < / _` |/ __| |/ / _ \ '_ \ / _` |
| |_) | (_| | (__|   <  __/ | | | (_| |
|____/ \__,_|\___|_|\_\___|_| |_|\__,_|
                                          
    服务器管理面板 一键卸载脚本
EOF
    echo -e "${NC}"
}

# 检查是否为root用户
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "此脚本需要以root权限运行"
        exit 1
    fi
}

# 检测操作系统
detect_os() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        OS=$NAME
        VER=$VERSION_ID
    else
        log_error "无法检测操作系统"
        exit 1
    fi
    
    log_info "检测到操作系统: $OS $VER"
}

# 检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 确认卸载
confirm_uninstall() {
    echo ""
    log_warn "警告：此操作将完全卸载 GPanel 及其相关组件！"
    echo "将要执行的操作："
    echo "  - 停止并禁用 GPanel 服务"
    echo "  - 删除 GPanel 系统服务"
    echo "  - 删除 GPanel 安装目录 ($INSTALL_DIR)"
    echo "  - 删除 GPanel 用户 ($USER)"
    
    if [[ "$REMOVE_NGINX" == "true" ]]; then
        echo "  - 删除 Nginx 配置"
    fi
    
    if [[ "$REMOVE_DATA" == "true" ]]; then
        echo "  - 删除所有数据文件"
    else
        echo "  - 保留数据文件"
    fi
    
    if [[ "$BACKUP_DATA" == "true" ]]; then
        echo "  - 备份数据到 /tmp/gpanel-backup-$(date +%Y%m%d-%H%M%S)"
    fi
    
    echo ""
    read -p "确定要继续吗？(输入 'yes' 确认): " confirm
    
    if [[ "$confirm" != "yes" ]]; then
        log_info "取消卸载操作"
        exit 0
    fi
}

# 停止服务
stop_services() {
    log_info "停止 GPanel 服务..."
    
    # 停止GPanel服务
    if systemctl is-active --quiet $SERVICE_NAME 2>/dev/null; then
        systemctl stop $SERVICE_NAME
        log_info "GPanel 服务已停止"
    fi
    
    # 禁用GPanel服务
    if systemctl is-enabled --quiet $SERVICE_NAME 2>/dev/null; then
        systemctl disable $SERVICE_NAME
        log_info "GPanel 服务已禁用"
    fi
}

# 备份数据
backup_data() {
    if [[ "$BACKUP_DATA" == "true" ]] && [[ -d "$INSTALL_DIR" ]]; then
        BACKUP_DIR="/tmp/gpanel-backup-$(date +%Y%m%d-%H%M%S)"
        log_info "备份数据到 $BACKUP_DIR..."
        
        mkdir -p $BACKUP_DIR
        
        # 备份配置文件
        if [[ -f "$INSTALL_DIR/config/config.yaml" ]]; then
            cp -r $INSTALL_DIR/config $BACKUP_DIR/
        fi
        
        # 备份数据库
        if [[ -f "$INSTALL_DIR/data/gpanel.db" ]]; then
            cp -r $INSTALL_DIR/data $BACKUP_DIR/
        fi
        
        # 备份插件
        if [[ -d "$INSTALL_DIR/plugins" ]]; then
            cp -r $INSTALL_DIR/plugins $BACKUP_DIR/
        fi
        
        # 备份日志
        if [[ -d "$INSTALL_DIR/logs" ]]; then
            cp -r $INSTALL_DIR/logs $BACKUP_DIR/
        fi
        
        # 备份Nginx配置
        if [[ -f "/etc/nginx/sites-available/gpanel" ]]; then
            mkdir -p $BACKUP_DIR/nginx
            cp /etc/nginx/sites-available/gpanel $BACKUP_DIR/nginx/
        fi
        
        # 创建备份信息文件
        cat > $BACKUP_DIR/backup_info.txt << EOF
备份时间: $(date)
备份原因: GPanel 卸载
GPanel 版本: $(cat $INSTALL_DIR/VERSION 2>/dev/null || echo "未知")
操作系统: $OS $VER
安装目录: $INSTALL_DIR
服务名称: $SERVICE_NAME

目录说明:
- config/: 配置文件
- data/: 数据库文件
- plugins/: 插件文件
- logs/: 日志文件
- nginx/: Nginx配置文件

恢复方法:
1. 停止相关服务
2. 将备份文件复制回原位置
3. 重启服务
EOF
        
        log_info "数据备份完成: $BACKUP_DIR"
    fi
}

# 删除系统服务
remove_service() {
    log_info "删除系统服务..."
    
    # 删除systemd服务文件
    if [[ -f "/etc/systemd/system/${SERVICE_NAME}.service" ]]; then
        rm -f /etc/systemd/system/${SERVICE_NAME}.service
        systemctl daemon-reload
        log_info "系统服务已删除"
    fi
    
    # 重载systemd
    systemctl daemon-reload
}

# 删除Nginx配置
remove_nginx_config() {
    if [[ "$REMOVE_NGINX" == "true" ]]; then
        log_info "删除 Nginx 配置..."
        
        # 删除Nginx站点配置
        if [[ -f "/etc/nginx/sites-available/gpanel" ]]; then
            rm -f /etc/nginx/sites-available/gpanel
        fi
        
        if [[ -f "/etc/nginx/sites-enabled/gpanel" ]]; then
            rm -f /etc/nginx/sites-enabled/gpanel
        fi
        
        # 测试Nginx配置
        if command_exists nginx; then
            nginx -t && systemctl reload nginx || log_warn "Nginx 配置重载失败"
        fi
        
        log_info "Nginx 配置已删除"
    fi
}

# 删除用户和组
remove_user() {
    log_info "删除用户和组..."
    
    if id "$USER" &>/dev/null; then
        # 强制删除用户及其主目录
        userdel -rf $USER 2>/dev/null || true
        log_info "用户 $USER 已删除"
    fi
}

# 删除安装目录
remove_installation() {
    if [[ -d "$INSTALL_DIR" ]]; then
        log_info "删除安装目录..."
        
        if [[ "$REMOVE_DATA" == "true" ]]; then
            # 完全删除
            rm -rf $INSTALL_DIR
            log_info "安装目录已完全删除"
        else
            # 保留数据，只删除程序文件
            if [[ -f "$INSTALL_DIR/gpanel" ]]; then
                rm -f $INSTALL_DIR/gpanel
            fi
            if [[ -d "$INSTALL_DIR/web" ]]; then
                rm -rf $INSTALL_DIR/web
            fi
            if [[ -d "$INSTALL_DIR/templates" ]]; then
                rm -rf $INSTALL_DIR/templates
            fi
            log_info "程序文件已删除，数据文件保留"
        fi
    fi
}

# 删除Docker相关
remove_docker() {
    if [[ "$REMOVE_DOCKER" == "true" ]]; then
        log_info "删除 Docker 相关..."
        
        # 停止并删除GPanel Docker容器
        if command_exists docker; then
            # 停止容器
            docker stop gpanel 2>/dev/null || true
            docker rm gpanel 2>/dev/null || true
            
            # 删除镜像
            docker rmi gpanel:latest 2>/dev/null || true
            
            # 删除网络
            docker network rm gpanel-network 2>/dev/null || true
            
            log_info "Docker 相关已删除"
        fi
    fi
}

# 清理系统
cleanup_system() {
    log_info "清理系统..."
    
    # 清理systemd失败的服务
    systemctl reset-failed 2>/dev/null || true
    
    # 清理临时文件
    rm -f /tmp/gpanel-* 2>/dev/null || true
    
    log_info "系统清理完成"
}

# 显示卸载结果
show_result() {
    log_info "卸载完成！"
    echo ""
    
    if [[ "$BACKUP_DATA" == "true" ]]; then
        echo -e "${GREEN}数据已备份到 /tmp/gpanel-backup-*${NC}"
    fi
    
    if [[ "$REMOVE_DATA" != "true" ]]; then
        echo -e "${YELLOW}数据文件保留在 $INSTALL_DIR${NC}"
    fi
    
    echo ""
    echo "手动清理项目（可选）："
    
    if [[ "$REMOVE_NGINX" != "true" ]] && command_exists nginx; then
        echo "  - 手动删除 Nginx 配置: rm -f /etc/nginx/sites-{available,enabled}/gpanel"
    fi
    
    if [[ "$REMOVE_DOCKER" != "true" ]] && command_exists docker; then
        echo "  - 手动删除 Docker 容器: docker stop gpanel && docker rm gpanel"
    fi
    
    if [[ -d "/var/www" ]]; then
        echo "  - 检查并清理网站文件: /var/www"
    fi
    
    if [[ -d "/etc/letsencrypt" ]]; then
        echo "  - 检查并清理SSL证书: /etc/letsencrypt"
    fi
    
    echo ""
    echo -e "${GREEN}感谢使用 GPanel！${NC}"
}

# 错误处理
error_handler() {
    log_error "卸载过程中发生错误"
    exit 1
}

# 主函数
main() {
    # 设置错误处理
    trap error_handler ERR
    
    show_banner
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            --remove-docker)
                REMOVE_DOCKER="true"
                shift
                ;;
            --remove-nginx)
                REMOVE_NGINX="true"
                shift
                ;;
            --remove-data)
                REMOVE_DATA="true"
                shift
                ;;
            --no-backup)
                BACKUP_DATA="false"
                shift
                ;;
            --help)
                echo "用法: $0 [选项]"
                echo "选项:"
                echo "  --remove-docker    同时删除Docker相关"
                echo "  --remove-nginx     同时删除Nginx配置"
                echo "  --remove-data      同时删除所有数据文件"
                echo "  --no-backup        不备份数据文件"
                echo "  --help             显示帮助信息"
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                exit 1
                ;;
        esac
    done
    
    log_info "开始卸载 GPanel..."
    
    # 执行卸载步骤
    check_root
    detect_os
    confirm_uninstall
    stop_services
    backup_data
    remove_service
    remove_nginx_config
    remove_docker
    remove_user
    remove_installation
    cleanup_system
    show_result
    
    log_info "GPanel 卸载完成！"
}

# 执行主函数
main "$@"