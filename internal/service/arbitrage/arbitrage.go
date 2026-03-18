// Package arbitrage 套利监控服务。
//
// 职责：
//   - 接收 scanner 产出的 Swap 事件流
//   - 在同一块内寻找同一代币的循环路径（三角套利 / 跨 DEX 套利）
//   - 产出 ArbitrageOpportunity
package arbitrage

import (
	"context"

	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// ─────────────────────────────────────────────
//  领域类型
// ─────────────────────────────────────────────

// StrategyType 套利策略类型
type StrategyType string

const (
	StrategyTriangular StrategyType = "triangular" // 三角套利（同 DEX）
	StrategyCrossDEX   StrategyType = "cross_dex"  // 跨 DEX 套利
)

// ArbitrageOpportunity 识别到的套利机会
type ArbitrageOpportunity struct {
	Strategy    StrategyType
	Path        []common.Address // 代币路径
	ProfitUSD   decimal.Decimal
	Confidence  float64
	BlockNumber string
	TxHashes    []common.Hash
}

// ─────────────────────────────────────────────
//  Config
// ─────────────────────────────────────────────

// Config 套利监控配置
type Config struct {
	MinProfitUSD decimal.Decimal // 最小利润阈值（USD）
	Strategies   []StrategyType  // 启用的策略（空 = 全部）
}

// DefaultConfig 默认配置
func DefaultConfig() Config {
	return Config{
		MinProfitUSD: decimal.NewFromInt(500), // 500U 起
		Strategies:   []StrategyType{StrategyTriangular, StrategyCrossDEX},
	}
}

// ─────────────────────────────────────────────
//  Service
// ─────────────────────────────────────────────

// Service 套利监控服务
type Service struct {
	cfg    Config
	logger *zap.Logger
}

// NewService 创建套利监控服务
func NewService(cfg Config, logger *zap.Logger) *Service {
	return &Service{cfg: cfg, logger: logger}
}

// ProcessBlock 分析一个区块内的 Swap 事件，寻找套利机会
func (s *Service) ProcessBlock(blockNum string, events []comm.UnifiedEvent) []ArbitrageOpportunity {
	// 过滤出 Swap 事件
	var swaps []comm.UnifiedEvent
	for _, ev := range events {
		if ev.GetEventType() == comm.EventMethodSwap {
			swaps = append(swaps, ev)
		}
	}
	if len(swaps) < 2 {
		return nil
	}

	var opportunities []ArbitrageOpportunity

	// 策略：三角套利检测
	if s.isEnabled(StrategyTriangular) {
		ops := s.detectTriangular(blockNum, swaps)
		opportunities = append(opportunities, ops...)
	}

	// 策略：跨 DEX 套利检测
	if s.isEnabled(StrategyCrossDEX) {
		ops := s.detectCrossDEX(blockNum, swaps)
		opportunities = append(opportunities, ops...)
	}

	return opportunities
}

// Run 持续消费事件 channel
func (s *Service) Run(ctx context.Context, in <-chan []comm.UnifiedEvent) <-chan ArbitrageOpportunity {
	out := make(chan ArbitrageOpportunity, 50)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case events, ok := <-in:
				if !ok {
					return
				}
				if len(events) == 0 {
					continue
				}
				blockNum := events[0].GetBlockNumber().String()
				for _, op := range s.ProcessBlock(blockNum, events) {
					select {
					case out <- op:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()
	return out
}

// ─────────────────────────────────────────────
//  内部检测逻辑（骨架，待填充图算法）
// ─────────────────────────────────────────────

func (s *Service) detectTriangular(blockNum string, swaps []comm.UnifiedEvent) []ArbitrageOpportunity {
	// TODO: 构建有向图，寻找长度为 3 的环路
	// 参考原 internal/service/core/arbitrageur/strategy/triangular.go
	_ = blockNum
	_ = swaps
	return nil
}

func (s *Service) detectCrossDEX(blockNum string, swaps []comm.UnifiedEvent) []ArbitrageOpportunity {
	// TODO: 按代币对分组，比较不同 DEX 的价格差异
	// 参考原 internal/service/core/arbitrageur/strategy/crossdex.go
	_ = blockNum
	_ = swaps
	return nil
}

func (s *Service) isEnabled(strategy StrategyType) bool {
	if len(s.cfg.Strategies) == 0 {
		return true
	}
	for _, st := range s.cfg.Strategies {
		if st == strategy {
			return true
		}
	}
	return false
}
