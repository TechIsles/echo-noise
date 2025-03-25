# 使用 Node.js Alpine 镜像作为前端构建阶段
FROM node:18-alpine AS frontend-build

# 设置前端工作目录
WORKDIR /app/web

# 复制前端依赖文件并安装依赖
COPY ./web/package.json ./web/package-lock.json* ./
RUN npm install --production

# 复制前端源代码并构建项目
COPY ./web/ .
RUN npm run generate  # 生成静态文件

# 复制构建后的文件到后端 public 目录
RUN cp -r .output/public /app/public/

# 使用 Golang Alpine 镜像作为后端构建阶段
FROM golang:1.22-alpine AS backend-build

# 设置后端工作目录
WORKDIR /app

# 添加阿里云的镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装构建时所需的工具（仅限构建阶段）
RUN apk add --no-cache --update gcc musl-dev

# 设置环境变量启用 cgo
ENV CGO_ENABLED=1
ENV GOTOOLCHAIN=auto

# 复制后端代码和必要文件
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./config ./config

# 创建并设置权限
RUN mkdir -p /app/data && chmod -R 777 /app/data

# 构建 Go 后端应用
RUN go build -o /app/noise ./cmd/server/main.go

# 使用更轻量的 Alpine 镜像作为运行时阶段
FROM alpine:3.18 AS final

WORKDIR /app

# 复制构建阶段的文件
COPY --from=backend-build /app/config /app/config
COPY --from=backend-build /app/noise /app/noise
COPY --from=frontend-build /app/public /app/public

# 暴露端口
EXPOSE 1314

# 运行后端服务
CMD ["/app/noise"]
