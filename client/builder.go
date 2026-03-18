// Package sentinel — Fluent Builder，提供链式 API 构建扫描任务
//
// 使用示例：
//
//	results, err := sentinel.NewScanBuilder(client).
//	    FromBlock(20000000).
//	    ToBlock(20001000).
//	    WithDEX(sentinel.ProtocolImplUniswapV2, sentinel.EventMethodSwap).
//	    WithToken(sentinel.ProtocolImplERC20, sentinel.EventMethodTransfer).
//	    Stream(ctx)
package sentinel

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ETHSentinel/internal/parser/comm"
)

// ScanBuilder 链式扫描配置构建器
type ScanBuilder struct {
	client     Client
	startBlock *big.Int
	endBlock   *big.Int
	cfg        comm.ParserCfg
	errs       []error
}

// NewScanBuilder 创建扫描构建器
func NewScanBuilder(client Client) *ScanBuilder {
	return &ScanBuilder{
		client: client,
		cfg:    make(comm.ParserCfg),
	}
}

// FromBlock 设置起始区块（整数）
func (b *ScanBuilder) FromBlock(n int64) *ScanBuilder {
	b.startBlock = big.NewInt(n)
	return b
}

// FromBlockBig 设置起始区块（*big.Int）
func (b *ScanBuilder) FromBlockBig(n *big.Int) *ScanBuilder {
	b.startBlock = new(big.Int).Set(n)
	return b
}

// ToBlock 设置结束区块（整数，不含）
func (b *ScanBuilder) ToBlock(n int64) *ScanBuilder {
	b.endBlock = big.NewInt(n)
	return b
}

// ToBlockBig 设置结束区块（*big.Int，不含）
func (b *ScanBuilder) ToBlockBig(n *big.Int) *ScanBuilder {
	b.endBlock = new(big.Int).Set(n)
	return b
}

// WithDEX 添加 DEX 协议解析（可传多个事件方法，不传 = 全部事件）
func (b *ScanBuilder) WithDEX(impl comm.ProtocolImpl, methods ...comm.EventMethod) *ScanBuilder {
	return b.addImpl(comm.ProtocolTypeDEX, impl, methods)
}

// WithToken 添加 Token 协议解析
func (b *ScanBuilder) WithToken(impl comm.ProtocolImpl, methods ...comm.EventMethod) *ScanBuilder {
	return b.addImpl(comm.ProtocolTypeToken, impl, methods)
}

// WithLending 添加借贷协议解析
func (b *ScanBuilder) WithLending(impl comm.ProtocolImpl, methods ...comm.EventMethod) *ScanBuilder {
	return b.addImpl(comm.ProtocolTypeLending, impl, methods)
}

func (b *ScanBuilder) addImpl(pt comm.ProtocolType, impl comm.ProtocolImpl, methods []comm.EventMethod) *ScanBuilder {
	if _, ok := b.cfg[pt]; !ok {
		b.cfg[pt] = make(comm.ImplConfig)
	}
	b.cfg[pt][impl] = methods
	return b
}

// Build 验证并构建 ParserCfg（不执行扫描）
func (b *ScanBuilder) Build() (comm.ParserCfg, error) {
	if len(b.errs) > 0 {
		return nil, fmt.Errorf("scan builder: %v", b.errs)
	}
	return b.cfg, nil
}

// ScanOne 扫描单个区块
func (b *ScanBuilder) ScanOne(ctx context.Context, blockNumber int64) (*ScanResult, error) {
	cfg, err := b.Build()
	if err != nil {
		return nil, err
	}
	return b.client.ScanBlock(ctx, big.NewInt(blockNumber), cfg)
}

// Stream 扫描区间，返回流式 channel
func (b *ScanBuilder) Stream(ctx context.Context) (<-chan *ScanResult, error) {
	if b.startBlock == nil || b.endBlock == nil {
		return nil, fmt.Errorf("scan builder: start/end block not set")
	}
	cfg, err := b.Build()
	if err != nil {
		return nil, err
	}
	return b.client.ScanBlocks(ctx, b.startBlock, b.endBlock, cfg)
}

// Subscribe 实时订阅新块
func (b *ScanBuilder) Subscribe(ctx context.Context) (<-chan *ScanResult, error) {
	cfg, err := b.Build()
	if err != nil {
		return nil, err
	}
	return b.client.SubscribeBlocks(ctx, cfg)
}
