MODULE := github.com/ETHSentinel
PROTO_DIR := server/proto
GEN_DIR   := gen

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
		$(PROTO_DIR)/sentinel/v1/sentinel.proto \
		$(PROTO_DIR)/sentinel/v1/mempool.proto \
		$(PROTO_DIR)/sentinel/v1/wallet.proto
	@echo ">>> Done: $(GEN_DIR)/sentinel/v1/"

.PHONY: install-tools
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: build
build:
	go build ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test ./... -race -count=1

.PHONY: run-server
run-server:
	ETH_RPC_URL=$(ETH_RPC_URL) \
	ETH_WS_URL=$(ETH_WS_URL)   \
	ETH_PROXY=$(ETH_PROXY)     \
	GRPC_ADDR=:50051           \
	go run ./test/server

.PHONY: run-example
run-example:
	go run ./test/example

.PHONY: clean
clean:
	rm -rf bin/ $(GEN_DIR)/sentinel/v1/*.pb.go
