# Stage 1: Build Frontend
FROM --platform=$BUILDPLATFORM node:20-alpine AS node-builder
WORKDIR /webapp
# 复制 package.json 和 lock 文件以利用 layer 缓存
COPY webapp/package*.json ./
# 使用 npm ci 以获得更快且可重现的构建（如果存在 package-lock.json）
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi
# 复制其余源码并构建
COPY webapp/ .
RUN npm run build

# Stage 2: Build Backend
FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine AS builder
# 安装 Git 以支持某些 Go 模块依赖
RUN apk add --no-cache git
WORKDIR /app
# 复制 go.mod 和 go.sum 以利用 layer 缓存
COPY go.mod go.sum ./
RUN go mod download
# 复制源码
COPY . .
# 从 node-builder 复制构建好的 web 目录，用于 go embed
COPY --from=node-builder /web ./web
# 编译二进制，针对目标平台
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w" -o BingPaper .

# Stage 3: Final Image
FROM alpine:3.21
# 安装运行时必需的证书和时区数据
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
# 创建必要目录
RUN mkdir -p data static
# 从构建阶段复制二进制文件
COPY --from=builder /app/BingPaper .
# 复制静态资源（如果有些资源没有被 embed）
COPY --from=builder /app/static ./static
# 复制默认配置
COPY config.example.yaml ./data/config.yaml

EXPOSE 8080
VOLUME ["/app/data"]
ENTRYPOINT ["./BingPaper"]
