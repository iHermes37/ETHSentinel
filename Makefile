# ETH Sentinel SDK — Makefile

MODULE := github.com/ETHSentinel
PROTO_DIR := proto
GEN_DIR   := gen

# ─────────────────────────────────────────────
#  Proto 代码生成
# ─────────────────────────────────────────────

.PHONY: proto
proto:
	@echo ">>> Generating protobuf code..."
	@mkdir -p $(GEN_DIR)
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/sentinel/v1/events.proto \
		$(PROTO_DIR)/sentinel/v1/sentinel.proto
	@echo ">>> Done: $(GEN_DIR)/sentinel/v1/"

# ─────────────────────────────────────────────
#  安装 protoc 插件（首次使用）
# ─────────────────────────────────────────────

.PHONY: install-tools
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# ─────────────────────────────────────────────
#  构建
# ─────────────────────────────────────────────

.PHONY: build
build: proto
	go build ./...

.PHONY: build-server
build-server:
	go build -o bin/sentinel-server ./cmd/server

# ─────────────────────────────────────────────
#  测试
# ─────────────────────────────────────────────

.PHONY: test
test:
	go test ./... -race -count=1

.PHONY: test-short
test-short:
	go test ./... -short -count=1

# ─────────────────────────────────────────────
#  代码质量
# ─────────────────────────────────────────────

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: tidy
tidy:
	go mod tidy

# ─────────────────────────────────────────────
#  运行
# ─────────────────────────────────────────────

.PHONY: run-server
run-server:
	ETH_RPC_URL=$(ETH_RPC_URL) \
	ETH_WS_URL=$(ETH_WS_URL)   \
	ETH_PROXY=$(ETH_PROXY)     \
	GRPC_ADDR=:50051           \
	go run ./cmd/server

.PHONY: run-example
run-example:
	go run ./cmd/example

# ─────────────────────────────────────────────
#  Docker
# ─────────────────────────────────────────────

.PHONY: docker-build
docker-build:
	docker build -t eth-sentinel:latest .

.PHONY: clean
clean:
	rm -rf bin/ $(GEN_DIR)/
