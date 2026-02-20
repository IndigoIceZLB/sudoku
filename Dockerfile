# ----------------------------
# 阶段 1: 构建阶段 (Builder)
# ----------------------------
FROM golang:alpine AS builder

# 设置工作目录
WORKDIR /app

# 为了加快构建，先复制依赖描述文件并下载依赖
COPY go.mod go.sum ./
# 使用国内代理（可选，Render构建时不需要，但写上也无妨），下载依赖
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go mod download

# 复制所有源代码
COPY . .

# 编译 Go 程序
# CGO_ENABLED=0: 关闭 CGO，确保生成静态链接的二进制文件
# GOOS=linux: 目标系统为 Linux
# -o server: 输出文件名为 server
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# ----------------------------
# 阶段 2: 运行阶段 (Runner)
# ----------------------------
# 使用极简的 Alpine 镜像
FROM alpine:latest

# 安装基础证书（如果后续要请求 HTTPS 外部 API，这一步很重要）
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 从构建阶段只复制编译好的二进制文件
COPY --from=builder /app/server .

# 声明端口（仅文档作用，实际由 Render 环境变量控制）
EXPOSE 8080

# 启动命令
CMD ["./server"]