# 多阶段构建 Dockerfile
FROM node:18-alpine AS frontend-builder

WORKDIR /app/web

# 复制前端依赖文件
COPY web/package*.json ./

# 安装前端依赖
RUN npm ci --only=production

# 复制前端源码
COPY web/ ./

# 构建前端
RUN npm run build

# Go 构建阶段
FROM golang:1.21-alpine AS backend-builder

# 安装必要的系统包
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# 复制 Go 依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gpanel cmd/server/main.go

# 最终运行镜像
FROM alpine:latest

# 安装必要的系统包
RUN apk --no-cache add ca-certificates tzdata nginx

# 创建非root用户
RUN addgroup -g 1001 -S gpanel && \
    adduser -u 1001 -S gpanel -G gpanel

WORKDIR /app

# 从构建阶段复制文件
COPY --from=backend-builder /app/gpanel .
COPY --from=frontend-builder /app/web/dist ./web/dist
COPY --from=backend-builder /app/internal ./internal
COPY --from=backend-builder /app/templates ./templates
COPY --from=backend-builder /app/internal/config/config.yaml ./config.yaml

# 创建必要的目录
RUN mkdir -p /app/data /app/logs /app/plugins && \
    chown -R gpanel:gpanel /app

# 复制 nginx 配置
COPY docker/nginx.conf /etc/nginx/nginx.conf
COPY docker/default.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80 443

# 切换到非root用户
USER gpanel

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost/health || exit 1

# 启动脚本
COPY docker/entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["server"]