// Package grpc — gRPC 拦截器（日志 + panic 恢复）
package grpc

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ─────────────────────────────────────────────
//  Unary 拦截器
// ─────────────────────────────────────────────

// LoggingUnaryInterceptor 记录每个一元 RPC 的耗时和错误
func LoggingUnaryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(start)
		if err != nil {
			logger.Error("grpc unary error",
				zap.String("method", info.FullMethod),
				zap.Duration("elapsed", elapsed),
				zap.Error(err),
			)
		} else {
			logger.Info("grpc unary",
				zap.String("method", info.FullMethod),
				zap.Duration("elapsed", elapsed),
			)
		}
		return resp, err
	}
}

// RecoveryUnaryInterceptor 捕获 panic 并转为 Internal 错误
func RecoveryUnaryInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("grpc unary panic",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.ByteString("stack", debug.Stack()),
				)
				err = status.Errorf(codes.Internal, "internal server error: %v", r)
			}
		}()
		return handler(ctx, req)
	}
}

// ─────────────────────────────────────────────
//  Stream 拦截器
// ─────────────────────────────────────────────

// LoggingStreamInterceptor 记录流式 RPC
func LoggingStreamInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)
		elapsed := time.Since(start)
		if err != nil {
			logger.Error("grpc stream error",
				zap.String("method", info.FullMethod),
				zap.Duration("elapsed", elapsed),
				zap.Error(err),
			)
		} else {
			logger.Info("grpc stream done",
				zap.String("method", info.FullMethod),
				zap.Duration("elapsed", elapsed),
			)
		}
		return err
	}
}

// RecoveryStreamInterceptor 捕获流式 RPC 中的 panic
func RecoveryStreamInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("grpc stream panic",
					zap.String("method", info.FullMethod),
					zap.Any("panic", r),
					zap.ByteString("stack", debug.Stack()),
				)
				err = status.Errorf(codes.Internal, fmt.Sprintf("stream panic: %v", r))
			}
		}()
		return handler(srv, ss)
	}
}
