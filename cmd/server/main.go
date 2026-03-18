// cmd/server — ETH Sentinel 完整 gRPC Server 启动入口
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ETHSentinel/internal/chain"
	"github.com/ETHSentinel/internal/conn"
	"github.com/ETHSentinel/internal/parser"
	"github.com/ETHSentinel/internal/scanner"
	grpcserver "github.com/ETHSentinel/server/grpc"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	rpcURL := envOrDefault("ETH_RPC_URL", "https://mainnet.infura.io/v3/YOUR_KEY")
	wsURL := envOrDefault("ETH_WS_URL", "wss://mainnet.infura.io/ws/v3/YOUR_KEY")
	proxy := envOrDefault("ETH_PROXY", "")
	grpcAddr := envOrDefault("GRPC_ADDR", ":50051")

	pool := conn.NewMultiChainPool(logger)
	pool.RegisterChain(chain.MustGet(chain.ChainETH), &conn.NodeConfig{
		Name: "ethereum", RPCURL: rpcURL, WSURL: wsURL, ProxyURL: proxy,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ethClient, err := pool.GetRPC(ctx, chain.ChainETH)
	if err != nil {
		logger.Fatal("connect eth rpc failed", zap.Error(err))
	}

	engine, err := parser.NewEngine(ethClient)
	if err != nil {
		logger.Fatal("init parser engine failed", zap.Error(err))
	}

	sc := scanner.New(ethClient, engine, logger)
	sentinelH := grpcserver.NewSentinelHandler(sc, ethClient, logger)
	mempoolH := grpcserver.NewMempoolHandler(pool, logger)
	walletH := grpcserver.NewWalletHandler(pool, logger)

	srv := grpcserver.NewServer(
		grpcserver.ServerConfig{Addr: grpcAddr},
		sentinelH, mempoolH, walletH, pool, logger,
	)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("gRPC server error", zap.Error(err))
		}
	}()
	logger.Info("ETH Sentinel started", zap.String("grpc", grpcAddr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutCancel()
	srv.Stop(shutCtx)
	pool.Close()
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
