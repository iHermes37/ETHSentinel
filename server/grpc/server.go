// Package grpc 提供 gRPC 服务启动和管理
package grpc

import (
	"context"
	"fmt"
	"net"

	sentinelv1 "github.com/ETHSentinel/gen/sentinel/v1"
	"github.com/ETHSentinel/internal/conn"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// ServerConfig gRPC 服务配置
type ServerConfig struct {
	Addr string
}

// Server gRPC 服务包装器
type Server struct {
	cfg        ServerConfig
	grpcServer *grpc.Server
	logger     *zap.Logger
}

// NewServer 创建并注册所有 gRPC 服务
func NewServer(
	cfg ServerConfig,
	sentinel *SentinelHandler,
	mempoolH *MempoolHandler,
	walletH *WalletHandler,
	_ *conn.MultiChainPool,
	logger *zap.Logger,
) *Server {
	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggingUnaryInterceptor(logger),
			RecoveryUnaryInterceptor(logger),
		),
		grpc.ChainStreamInterceptor(
			LoggingStreamInterceptor(logger),
			RecoveryStreamInterceptor(logger),
		),
	)

	sentinelv1.RegisterSentinelServiceServer(grpcSrv, sentinel)
	sentinelv1.RegisterMempoolServiceServer(grpcSrv, mempoolH)
	sentinelv1.RegisterWalletServiceServer(grpcSrv, walletH)

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcSrv, healthSrv)
	healthSrv.SetServingStatus("sentinel.v1.SentinelService", grpc_health_v1.HealthCheckResponse_SERVING)
	healthSrv.SetServingStatus("sentinel.v1.MempoolService", grpc_health_v1.HealthCheckResponse_SERVING)
	healthSrv.SetServingStatus("sentinel.v1.WalletService", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(grpcSrv)

	return &Server{cfg: cfg, grpcServer: grpcSrv, logger: logger}
}

// Start 启动 gRPC 服务（阻塞）
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return fmt.Errorf("grpc: listen on %s: %w", s.cfg.Addr, err)
	}
	s.logger.Info("gRPC server starting", zap.String("addr", s.cfg.Addr))
	return s.grpcServer.Serve(lis)
}

// Stop 优雅停止
func (s *Server) Stop(ctx context.Context) {
	s.logger.Info("gRPC server stopping")
	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()
	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
	case <-stopped:
	}
	s.logger.Info("gRPC server stopped")
}
