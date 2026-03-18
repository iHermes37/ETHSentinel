// Package grpc 提供 gRPC 服务启动和管理
package grpc

import (
	"context"
	"fmt"
	"net"

	sentinelv1 "github.com/ETHSentinel/gen/sentinel/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// ServerConfig gRPC 服务配置
type ServerConfig struct {
	Addr string // 监听地址，如 ":50051"
}

// Server gRPC 服务包装器
type Server struct {
	cfg        ServerConfig
	grpcServer *grpc.Server
	handler    *SentinelHandler
	logger     *zap.Logger
}

// NewServer 创建 gRPC Server
func NewServer(cfg ServerConfig, handler *SentinelHandler, logger *zap.Logger) *Server {
	// 拦截器链：日志 → 恢复 → 认证（可扩展）
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

	s := &Server{
		cfg:        cfg,
		grpcServer: grpcSrv,
		handler:    handler,
		logger:     logger,
	}

	// 注册业务服务
	sentinelv1.RegisterSentinelServiceServer(grpcSrv, handler)

	// 注册健康检查
	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcSrv, healthSrv)
	healthSrv.SetServingStatus("sentinel.v1.SentinelService", grpc_health_v1.HealthCheckResponse_SERVING)

	// 注册反射（grpcurl / Evans 等工具使用）
	reflection.Register(grpcSrv)

	return s
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
