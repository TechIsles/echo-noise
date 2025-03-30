# 前端构建阶段
FROM node:22.14.0-alpine AS frontend-build

# 设置工作目录
WORKDIR /app/web

# 复制前端依赖文件并安装依赖
COPY ./web/package.json ./web/package-lock.json* ./
RUN npm install --production

# 复制前端源代码并构建
COPY ./web/ .
RUN npm run generate

# 将构建结果复制到公共目录
RUN cp -r .output/public /app/public/

# 后端构建阶段
FROM golang:1.24.1-alpine AS backend-build

# 设置工作目录
WORKDIR /app

# 安装构建时所需的工具
RUN apk add --no-cache gcc musl-dev

# 禁用 CGO（如果不依赖 CGO，可以禁用以避免安装 gcc）
# ENV CGO_ENABLED=0

# 复制 Go 模块文件并下载依赖
COPY ./go.mod ./go.sum ./
RUN go mod download

# 复制项目文件
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./config ./config

# 创建必要的目录并设置权限
RUN mkdir -p /app/data /app/public && chmod -R 755 /app/data

# 编译 Go 应用
RUN go build -o /app/noise ./cmd/server/main.go

# 运行时阶段
FROM alpine:latest AS final

# 设置工作目录
WORKDIR /app

# 从后端构建阶段复制配置文件和二进制文件
COPY --from=backend-build /app/config /app/config
COPY --from=backend-build /app/noise /app/noise

# 从前端构建阶段复制静态文件
COPY --from=frontend-build /app/public /app/public

# 安装运行时所需的工具（如 CA 证书）
RUN apk add --no-cache ca-certificates

# 复制数据库文件到运行时目录
COPY ./data/noise.db /app/data/

# 暴露应用端口
EXPOSE 1314

# 启动应用
CMD ["/app/noise"]
