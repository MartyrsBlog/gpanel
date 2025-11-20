# GPanel - 服务器管理面板

<div align="center">

![GPanel Logo](docs/images/logo.png)

**基于 Go + Gin + Vue3 的现代化服务器管理面板**

[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Vue](https://img.shields.io/badge/Vue-3.3+-green.svg)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Build](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/your-repo/gpanel/actions)

[功能特性](#功能特性) • [快速开始](#快速开始) • [安装指南](#安装指南) • [使用文档](#使用文档) • [开发指南](#开发指南)

</div>

## 📖 项目介绍

GPanel 是一个现代化的服务器管理面板，采用 Go + Gin 后端和 Vue3 前端技术栈构建。它提供了类似宝塔面板的功能，包括系统监控、网站管理、文件管理、数据库管理、Docker 管理等，同时支持插件扩展。

### 🎯 设计目标

- **简洁易用** - 现代化的用户界面，操作简单直观
- **高性能** - 基于 Go 语言开发，内存占用低，响应速度快
- **安全可靠** - 内置安全机制，支持 SSL/TLS 加密
- **插件化** - 支持插件系统，功能可灵活扩展
- **跨平台** - 支持 Linux 主流发行版

## ✨ 功能特性

### 🖥️ 系统监控
- **实时监控** - CPU、内存、磁盘、网络使用情况
- **进程管理** - 查看和管理系统进程
- **系统信息** - 详细的系统和硬件信息
- **性能图表** - 可视化性能数据展示

### 🌐 网站管理
- **Nginx 管理** - 虚拟主机创建、配置、删除
- **SSL 证书** - 自动申请 Let's Encrypt 证书
- **域名管理** - 多域名支持和管理
- **PHP 版本** - 支持多 PHP 版本切换

### 📁 文件管理
- **在线浏览** - 文件和目录浏览
- **文件操作** - 上传、下载、复制、移动、删除
- **压缩解压** - 支持 ZIP、TAR.GZ 等格式
- **权限管理** - 文件和目录权限设置

### 🗄️ 数据库管理
- **MySQL/MariaDB** - 数据库创建、删除、备份
- **用户管理** - 数据库用户权限管理
- **备份恢复** - 自动备份和手动恢复
- **性能监控** - 数据库性能状态监控

### 🐳 Docker 管理
- **容器管理** - 容器的创建、启动、停止、删除
- **镜像管理** - 镜像的拉取、删除、管理
- **日志查看** - 容器日志实时查看
- **资源监控** - 容器资源使用情况

### 🔌 插件系统
- **插件加载** - 动态加载和卸载插件
- **API 扩展** - 插件可扩展 API 接口
- **配置管理** - 插件配置统一管理
- **示例插件** - 提供 Redis 管理插件示例

## 🚀 快速开始

### 环境要求

- **操作系统**: Linux (Ubuntu 18.04+, Debian 9+, CentOS 7+, Arch Linux)
- **Go**: 1.21 或更高版本
- **Node.js**: 18.0 或更高版本
- **内存**: 最低 512MB，推荐 1GB+
- **磁盘**: 最低 1GB 可用空间

### 一键安装

```bash
# 下载安装脚本
curl -fsSL https://raw.githubusercontent.com/your-repo/gpanel/main/scripts/install.sh -o install.sh

# 执行安装
chmod +x install.sh
sudo ./install.sh
```

### Docker 部署

```bash
# 克隆项目
git clone https://github.com/your-repo/gpanel.git
cd gpanel

# 启动服务
docker-compose up -d
```

### 手动安装

```bash
# 1. 克隆项目
git clone https://github.com/your-repo/gpanel.git
cd gpanel

# 2. 安装依赖
make install-deps

# 3. 构建应用
make build

# 4. 启动服务
make run
```

## 📖 安装指南

### 系统要求

详细的环境要求和兼容性说明：

| 组件 | 最低版本 | 推荐版本 |
|------|---------|---------|
| 操作系统 | Ubuntu 18.04+ | Ubuntu 20.04+ |
| Go | 1.21 | 1.21+ |
| Node.js | 18.0 | 18.0+ |
| 内存 | 512MB | 1GB+ |
| 磁盘 | 1GB | 5GB+ |

### 安装步骤

#### 方法一：一键安装脚本（推荐）

```bash
# 基础安装
curl -fsSL https://raw.githubusercontent.com/your-repo/gpanel/main/scripts/install.sh | sudo bash

# 指定端口安装
curl -fsSL https://raw.githubusercontent.com/your-repo/gpanel/main/scripts/install.sh | sudo bash -s -- --port 8888

# 不安装 Docker 和 Nginx
curl -fsSL https://raw.githubusercontent.com/your-repo/gpanel/main/scripts/install.sh | sudo bash -s -- --no-docker --no-nginx
```

#### 方法二：源码编译安装

```bash
# 1. 安装 Go 和 Node.js
# Ubuntu/Debian
sudo apt update
sudo apt install -y golang nodejs npm

# CentOS/RHEL
sudo yum install -y golang nodejs npm

# 2. 克隆代码
git clone https://github.com/your-repo/gpanel.git
cd gpanel

# 3. 安装依赖
go mod download
cd web && npm install && cd ..

# 4. 构建前端
cd web && npm run build && cd ..

# 5. 构建后端
go build -o gpanel cmd/server/main.go

# 6. 启动服务
./gpanel
```

#### 方法三：Docker 安装

```bash
# 使用 Docker Compose
git clone https://github.com/your-repo/gpanel.git
cd gpanel
docker-compose up -d

# 或单独使用 Docker
docker run -d \
  --name gpanel \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v /opt/gpanel:/app/data \
  gpanel:latest
```

### 配置说明

主要配置文件位于 `/etc/gpanel/config.yaml`：

```yaml
server:
  port: "8080"
  host: "0.0.0.0"

database:
  type: "sqlite"  # sqlite or mysql
  sqlite:
    path: "./data/gpanel.db"
  mysql:
    host: "localhost"
    port: "3306"
    database: "gpanel"
    username: "root"
    password: ""

auth:
  jwt_secret: "your-secret-key"
  token_expire: 86400  # 24 hours
  default_user: "admin"
  default_passwd: "admin123"

ssl:
  email: "admin@example.com"
  staging: false
```

## 📚 使用文档

### 基础使用

#### 登录系统

1. 打开浏览器访问 `http://your-server-ip:8080`
2. 使用默认账号登录：
   - 用户名：`admin`
   - 密码：`admin123`

#### 仪表盘

仪表盘提供了系统概览信息：

- **系统状态** - CPU、内存、磁盘使用率
- **实时监控** - 网络流量、系统负载
- **快捷操作** - 常用功能快速入口
- **通知消息** - 系统通知和警告

#### 网站管理

1. **创建网站**
   - 点击"网站管理" → "创建网站"
   - 填写域名、网站目录等信息
   - 选择 PHP 版本（可选）
   - 点击创建完成

2. **SSL 证书**
   - 选择网站 → "SSL 管理"
   - 填写邮箱地址
   - 点击"申请证书"自动申请 Let's Encrypt 证书

#### 文件管理

文件管理器功能：

- **浏览文件** - 点击目录进入，支持路径导航
- **上传文件** - 拖拽或点击上传按钮
- **文件操作** - 右键菜单或工具栏操作
- **压缩解压** - 支持常见压缩格式

#### 数据库管理

1. **创建数据库**
   - 进入"数据库管理"页面
   - 点击"创建数据库"
   - 填写数据库名、用户名、密码
   - 选择数据库类型（MySQL/MariaDB）

2. **备份数据库**
   - 选择数据库 → "备份"
   - 下载备份文件或保存到服务器

### 高级功能

#### Docker 管理

Docker 功能需要服务器已安装 Docker：

- **容器列表** - 查看所有容器状态
- **容器操作** - 启动、停止、重启、删除
- **镜像管理** - 拉取、删除、查看镜像
- **日志查看** - 实时查看容器日志

#### 插件管理

1. **安装插件**
   - 将插件文件放置到 `plugins/` 目录
   - 在插件管理页面扫描插件
   - 点击启用插件

2. **开发插件**
   - 参考插件开发文档
   - 实现 Plugin 接口
   - 编译为 .so 文件

## 🔧 开发指南

### 开发环境搭建

```bash
# 1. 克隆项目
git clone https://github.com/your-repo/gpanel.git
cd gpanel

# 2. 安装依赖
make install-deps

# 3. 启动开发环境
make dev
```

### 项目结构

```
gpanel/
├── cmd/                    # 应用入口
│   └── server/
│       └── main.go        # 主程序
├── internal/              # 内部包
│   ├── api/              # API 处理
│   ├── config/           # 配置管理
│   ├── middleware/       # 中间件
│   ├── plugin/           # 插件系统
│   └── ...
├── web/                  # 前端代码
│   ├── src/
│   ├── public/
│   └── ...
├── templates/            # 模板文件
├── scripts/             # 脚本文件
├── docker/              # Docker 配置
└── docs/               # 文档
```

### API 文档

启动服务后访问 `http://localhost:8080/swagger/index.html` 查看完整 API 文档。

#### 主要 API 接口

```bash
# 认证相关
POST /api/auth/login      # 用户登录
POST /api/auth/logout     # 用户登出
GET  /api/auth/info       # 获取用户信息

# 系统监控
GET  /api/system/monitor  # 系统监控数据
GET  /api/system/processes # 进程列表
GET  /api/system/disk     # 磁盘信息

# 网站管理
GET  /api/site/list       # 网站列表
POST /api/site/create     # 创建网站
POST /api/site/delete     # 删除网站

# 文件管理
GET  /api/file/list       # 文件列表
POST /api/file/upload     # 上传文件
GET  /api/file/download   # 下载文件

# 数据库管理
GET  /api/database/list   # 数据库列表
POST /api/database/create # 创建数据库
GET  /api/database/backup # 备份数据库

# Docker 管理
GET  /api/docker/containers # 容器列表
POST /api/docker/container/start # 启动容器
GET  /api/docker/images   # 镜像列表
```

### 插件开发

#### 插件接口

```go
type Plugin interface {
    Name() string
    Version() string
    Description() string
    Author() string
    Init() error
    Routes(rg *gin.RouterGroup)
    Cleanup() error
}
```

#### 插件示例

```go
package main

import (
    "github.com/gin-gonic/gin"
    "gpanel/internal/plugin"
)

type RedisPlugin struct {
    plugin.BasePlugin
}

func NewPlugin() plugin.Plugin {
    return &RedisPlugin{
        BasePlugin: plugin.NewBasePlugin(
            "redis",
            "1.0.0",
            "Redis management plugin",
            "GPanel Team",
        ),
    }
}

func (p *RedisPlugin) Routes(rg *gin.RouterGroup) {
    rg.GET("/status", p.getRedisStatus)
    rg.POST("/restart", p.restartRedis)
}

func (p *RedisPlugin) getRedisStatus(c *gin.Context) {
    // 实现 Redis 状态查询
}

func (p *RedisPlugin) restartRedis(c *gin.Context) {
    // 实现 Redis 重启
}
```

#### 插件编译

```bash
# 编译插件
go build -buildmode=plugin -o redis.so redis.go

# 创建插件清单
cat > manifest.json << EOF
{
  "name": "redis",
  "version": "1.0.0",
  "description": "Redis management plugin",
  "author": "GPanel Team",
  "entry": "redis.so",
  "enabled": true
}
EOF
```

### 贡献指南

我们欢迎所有形式的贡献！

#### 提交 Issue

- 使用 Issue 模板报告 Bug
- 提供详细的重现步骤
- 包含系统环境信息

#### 提交 PR

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

#### 代码规范

- Go 代码遵循 `gofmt` 和 `golint` 规范
- 前端代码使用 ESLint 和 Prettier
- 提交信息遵循 Conventional Commits

## 📋 待办事项

- [ ] 添加更多数据库支持 (PostgreSQL, MongoDB)
- [ ] 实现集群管理功能
- [ ] 添加监控告警系统
- [ ] 支持容器化部署
- [ ] 添加更多插件示例
- [ ] 实现备份恢复功能
- [ ] 添加多语言支持
- [ ] 实现主题切换功能

## 🤝 贡献者

感谢所有为 GPanel 做出贡献的开发者！

<a href="https://github.com/your-repo/gpanel/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=your-repo/gpanel" />
</a>

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架
- [Vue.js](https://github.com/vuejs/vue) - 前端框架
- [Element Plus](https://github.com/element-plus/element-plus) - UI 组件库
- [ECharts](https://github.com/apache/echarts) - 图表库
- [gopsutil](https://github.com/shirou/gopsutil) - 系统信息库

## 📞 联系我们

- **官方网站**: https://gpanel.dev
- **文档地址**: https://docs.gpanel.dev
- **问题反馈**: https://github.com/your-repo/gpanel/issues
- **交流群组**: [加入讨论](https://github.com/your-repo/gpanel/discussions)

---

⭐ 如果这个项目对你有帮助，请给我们一个 Star！