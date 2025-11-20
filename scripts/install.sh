#!/bin/bash

# GPanel 一键安装脚本
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
VERSION="latest"
PORT="8080"
WITH_DOCKER="true"
WITH_NGINX="true"

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
                                          
    服务器管理面板 一键安装脚本
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

# 更新系统
update_system() {
    log_info "更新系统包..."
    
    case $OS in
        "Ubuntu"|"Debian"*)
            apt update && apt upgrade -y
            ;;
        "CentOS"*|"Red Hat"*)
            yum update -y
            ;;
        "Arch Linux")
            pacman -Syu --noconfirm
            ;;
        *)
            log_error "不支持的操作系统: $OS"
            exit 1
            ;;
    esac
}

# 安装基础依赖
install_dependencies() {
    log_info "安装基础依赖..."
    
    case $OS in
        "Ubuntu"|"Debian"*)
            apt install -y curl wget git unzip tar sudo
            ;;
        "CentOS"*|"Red Hat"*)
            yum install -y curl wget git unzip tar sudo
            ;;
        "Arch Linux")
            pacman -S --noconfirm curl wget git unzip tar sudo
            ;;
    esac
}

# 安装Go
install_go() {
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}')
        log_info "Go 已安装: $GO_VERSION"
        return
    fi
    
    log_info "安装 Go..."
    GO_VERSION="1.21.5"
    cd /tmp
    
    if [[ $(uname -m) == "x86_64" ]]; then
        ARCH="amd64"
    else
        ARCH="386"
    fi
    
    wget -q https://golang.org/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz
    tar -C /usr/local -xzf go${GO_VERSION}.linux-${ARCH}.tar.gz
    
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    export PATH=$PATH:/usr/local/go/bin
    
    rm -f go${GO_VERSION}.linux-${ARCH}.tar.gz
    log_info "Go 安装完成"
}

# 安装Node.js
install_nodejs() {
    if command_exists node && command_exists npm; then
        NODE_VERSION=$(node -v)
        NPM_VERSION=$(npm -v)
        log_info "Node.js 已安装: $NODE_VERSION"
        log_info "npm 已安装: $NPM_VERSION"
        return
    fi
    
    log_info "安装 Node.js..."
    cd /tmp
    
    # 使用NodeSource仓库安装
    if [[ $OS == *"Ubuntu"* ]] || [[ $OS == *"Debian"* ]]; then
        curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
        apt install -y nodejs
    elif [[ $OS == *"CentOS"* ]] || [[ $OS == *"Red Hat"* ]]; then
        curl -fsSL https://rpm.nodesource.com/setup_18.x | bash -
        yum install -y nodejs
    elif [[ $OS == *"Arch Linux"* ]]; then
        pacman -S --noconfirm nodejs npm
    fi
    
    log_info "Node.js 安装完成"
}

# 安装Docker
install_docker() {
    if [[ "$WITH_DOCKER" != "true" ]]; then
        return
    fi
    
    if command_exists docker; then
        DOCKER_VERSION=$(docker --version)
        log_info "Docker 已安装: $DOCKER_VERSION"
        return
    fi
    
    log_info "安装 Docker..."
    
    case $OS in
        "Ubuntu"|"Debian"*)
            apt install -y apt-transport-https ca-certificates curl gnupg lsb-release
            curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
            echo "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list
            apt update && apt install -y docker-ce docker-ce-cli containerd.io
            ;;
        "CentOS"*|"Red Hat"*)
            yum install -y yum-utils
            yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
            yum install -y docker-ce docker-ce-cli containerd.io
            ;;
        "Arch Linux")
            pacman -S --noconfirm docker
            ;;
    esac
    
    # 启动Docker服务
    systemctl enable docker
    systemctl start docker
    
    log_info "Docker 安装完成"
}

# 安装Nginx
install_nginx() {
    if [[ "$WITH_NGINX" != "true" ]]; then
        return
    fi
    
    if command_exists nginx; then
        NGINX_VERSION=$(nginx -v 2>&1)
        log_info "Nginx 已安装: $NGINX_VERSION"
        return
    fi
    
    log_info "安装 Nginx..."
    
    case $OS in
        "Ubuntu"|"Debian"*)
            apt install -y nginx
            ;;
        "CentOS"*|"Red Hat"*)
            yum install -y nginx
            ;;
        "Arch Linux")
            pacman -S --noconfirm nginx
            ;;
    esac
    
    # 启动Nginx服务
    systemctl enable nginx
    systemctl start nginx
    
    log_info "Nginx 安装完成"
}

# 创建用户
create_user() {
    if id "$USER" &>/dev/null; then
        log_info "用户 $USER 已存在"
    else
        log_info "创建用户: $USER"
        useradd -r -s /bin/false -d $INSTALL_DIR $USER
    fi
    
    # 如果安装了Docker，将用户添加到docker组
    if [[ "$WITH_DOCKER" == "true" ]] && command_exists docker; then
        log_info "将用户 $USER 添加到docker组"
        usermod -aG docker $USER
    fi
}

# 创建目录
create_directories() {
    log_info "创建安装目录..."
    mkdir -p $INSTALL_DIR
    mkdir -p $INSTALL_DIR/{data,logs,plugins,config,templates,web}
    mkdir -p /var/www
    mkdir -p /etc/letsencrypt
}

# 下载GPanel
download_gpanel() {
    log_info "下载 GPanel..."
    cd /tmp
    
    # 清理可能存在的旧目录
    rm -rf gpanel
    
    # 克隆源码仓库
    if ! git clone https://github.com/MartyrsBlog/gpanel.git; then
        log_error "克隆仓库失败，请检查网络连接"
        exit 1
    fi
    
    # 复制源码到安装目录
    cp -r gpanel/* $INSTALL_DIR/
    
    # 构建前端
    log_info "构建前端..."
    cd $INSTALL_DIR/web
    npm install
    npm run build
    
    # 构建后端
    log_info "构建后端..."
    cd $INSTALL_DIR
    export PATH=$PATH:/usr/local/go/bin
    go mod download
    go mod tidy
    go build -o gpanel cmd/server/main.go
    
    # 清理临时文件
    rm -rf /tmp/gpanel
    
    log_info "GPanel 构建完成"
}

# 设置权限
set_permissions() {
    log_info "设置文件权限..."
    chown -R $USER:$USER $INSTALL_DIR
    chmod +x $INSTALL_DIR/gpanel
    chmod -R 755 $INSTALL_DIR
}

# 创建systemd服务
create_service() {
    log_info "创建系统服务..."
    
    cat > /etc/systemd/system/${SERVICE_NAME}.service << EOF
[Unit]
Description=GPanel Server Management Panel
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/gpanel
Restart=always
RestartSec=5
Environment=GIN_MODE=release
Environment=PORT=$PORT

[Install]
WantedBy=multi-user.target
EOF
    
    systemctl daemon-reload
    systemctl enable $SERVICE_NAME
}

# 配置Nginx
configure_nginx() {
    if [[ "$WITH_NGINX" != "true" ]]; then
        return
    fi
    
    log_info "配置 Nginx..."
    
    cat > /etc/nginx/sites-available/gpanel << EOF
server {
    listen 80;
    server_name _;
    
    location / {
        proxy_pass http://127.0.0.1:$PORT;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF
    
    # 启用站点
    ln -sf /etc/nginx/sites-available/gpanel /etc/nginx/sites-enabled/
    rm -f /etc/nginx/sites-enabled/default
    
    # 测试配置
    nginx -t
    
    # 重启Nginx
    systemctl restart nginx
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    # 启动GPanel
    systemctl start $SERVICE_NAME
    
    # 检查服务状态
    sleep 3
    if systemctl is-active --quiet $SERVICE_NAME; then
        log_info "GPanel 服务启动成功"
    else
        log_error "GPanel 服务启动失败"
        systemctl status $SERVICE_NAME
        exit 1
    fi
}

# 显示安装结果
show_result() {
    log_info "安装完成！"
    echo ""
    echo -e "${GREEN}访问地址: http://localhost:$PORT${NC}"
    echo -e "${GREEN}默认账号: admin${NC}"
    echo -e "${GREEN}默认密码: admin123${NC}"
    echo ""
    echo "常用命令:"
    echo "  启动服务: systemctl start $SERVICE_NAME"
    echo "  停止服务: systemctl stop $SERVICE_NAME"
    echo "  重启服务: systemctl restart $SERVICE_NAME"
    echo "  查看状态: systemctl status $SERVICE_NAME"
    echo "  查看日志: journalctl -u $SERVICE_NAME -f"
    echo ""
    echo -e "${YELLOW}请及时修改默认密码以确保安全！${NC}"
}

# 清理函数
cleanup() {
    log_info "清理临时文件..."
    rm -f /tmp/gpanel-*
}

# 错误处理
error_handler() {
    log_error "安装过程中发生错误，正在清理..."
    cleanup
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
            --version)
                VERSION="$2"
                shift 2
                ;;
            --port)
                PORT="$2"
                shift 2
                ;;
            --no-docker)
                WITH_DOCKER="false"
                shift
                ;;
            --no-nginx)
                WITH_NGINX="false"
                shift
                ;;
            --help)
                echo "用法: $0 [选项]"
                echo "选项:"
                echo "  --version VERSION    指定版本 (默认: latest)"
                echo "  --port PORT          指定端口 (默认: 8080)"
                echo "  --no-docker          不安装Docker"
                echo "  --no-nginx           不安装Nginx"
                echo "  --help               显示帮助信息"
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                exit 1
                ;;
        esac
    done
    
    log_info "开始安装 GPanel..."
    log_info "安装目录: $INSTALL_DIR"
    log_info "服务端口: $PORT"
    
    # 执行安装步骤
    check_root
    detect_os
    update_system
    install_dependencies
    install_go
    install_nodejs
    install_docker
    install_nginx
    create_user
    create_directories
    download_gpanel
    set_permissions
    create_service
    configure_nginx
    start_services
    show_result
    cleanup
    
    log_info "GPanel 安装完成！"
}

# 执行主函数
main "$@"