// Package sentinel — Functional Options 配置模式
package sentinel

import (
	"go.uber.org/zap"
	"time"
)

// options SDK 内部配置项
type options struct {
	// 连接
	rpcURL   string
	wsURL    string
	proxyURL string

	// 行为
	workerPoolSize int
	dialTimeout    time.Duration

	// 日志
	logger *zap.Logger

	// 多链
	chainID  uint64

	// 钱包
	mnemonic string

	// gRPC（当使用远程模式时）
	grpcAddr    string
	grpcTimeout time.Duration
}

func defaultOptions() *options {
	logger, _ := zap.NewProduction()
	return &options{
		workerPoolSize: 100,
		chainID:        1, // ETH 主网
		dialTimeout:    10 * time.Second,
		grpcTimeout:    30 * time.Second,
		logger:         logger,
	}
}

// Option SDK 配置函数类型
type Option func(*options)

// WithRPCURL 设置以太坊 RPC 节点 URL
func WithRPCURL(url string) Option {
	return func(o *options) { o.rpcURL = url }
}

// WithWSURL 设置以太坊 WebSocket 节点 URL（用于实时订阅）
func WithWSURL(url string) Option {
	return func(o *options) { o.wsURL = url }
}

// WithProxy 设置代理地址（HTTP / SOCKS5）
func WithProxy(proxyURL string) Option {
	return func(o *options) { o.proxyURL = proxyURL }
}

// WithWorkerPoolSize 设置并发扫块的协程池大小（默认 100）
func WithWorkerPoolSize(size int) Option {
	return func(o *options) { o.workerPoolSize = size }
}

// WithLogger 注入自定义 zap 日志器
func WithLogger(logger *zap.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// WithGRPCAddr 设置远程 gRPC 服务地址（若使用 gRPC 模式）
func WithGRPCAddr(addr string) Option {
	return func(o *options) { o.grpcAddr = addr }
}

// WithDialTimeout 设置连接超时
func WithDialTimeout(d time.Duration) Option {
	return func(o *options) { o.dialTimeout = d }
}

// WithChainID 设置目标链 ID（默认 1 = ETH 主网）
func WithChainID(id uint64) Option {
	return func(o *options) { o.chainID = id }
}

// WithMnemonic 设置钱包助记词（用于 Wallet 功能）
func WithMnemonic(mnemonic string) Option {
	return func(o *options) { o.mnemonic = mnemonic }
}
