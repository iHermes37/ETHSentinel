# ── 阶段 1：生成 proto + 编译 ──────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 安装 protoc 和插件
RUN apk add --no-cache protobuf protobuf-dev && \
    go install google.golang.org/protobuf/test/protoc-gen-go@latest && \
    go install google.golang.org/grpc/test/protoc-gen-go-grpc@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 生成 proto 代码
RUN protoc \
    --proto_path=proto \
    --go_out=gen \
    --go_opt=paths=source_relative \
    --go-grpc_out=gen \
    --go-grpc_opt=paths=source_relative \
    proto/sentinel/v1/events.proto \
    proto/sentinel/v1/sentinel.proto

# 编译二进制
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /sentinel ./test/server

# ── 阶段 2：精简运行镜像 ────────────────────────
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /sentinel /sentinel
COPY config/ /app/config/

EXPOSE 50051

ENTRYPOINT ["/sentinel"]
