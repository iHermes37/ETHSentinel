# ── 阶段 1：生成 proto + 编译 ──────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 安装 protoc 和插件
RUN apk add --no-cache protobuf protobuf-dev && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 改动一：proto_path 改为 server/proto
# 改动二：新增 mempool.proto 和 wallet.proto
RUN protoc \
    --proto_path=server/proto \
    --go_out=gen \
    --go_opt=paths=source_relative \
    --go-grpc_out=gen \
    --go-grpc_opt=paths=source_relative \
    server/proto/sentinel/v1/events.proto \
    server/proto/sentinel/v1/sentinel.proto \
    server/proto/sentinel/v1/mempool.proto \
    server/proto/sentinel/v1/wallet.proto

# 改动三：入口改为 ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /sentinel ./cmd/server

# ── 阶段 2：精简运行镜像 ────────────────────────
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /sentinel /sentinel

EXPOSE 50051

ENTRYPOINT ["/sentinel"]