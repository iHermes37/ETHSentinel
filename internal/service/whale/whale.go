// Package whale 鲸鱼监控服务。
//
// 职责：
//   - 消费 scanner 产出的 UnifiedEvent 流
//   - 根据阈值判断是否为鲸鱼行为
//   - 产出 WhaleBehaviorEvent，由上层（gRPC handler / DB writer）处理
//
// 重构要点：
//   - 去掉原来的全局函数风格，改为 Service 结构体
//   - 阈值和过滤逻辑通过 Config 注入，可测试
package whale

import (
	"context"
	"math/big"

	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// ─────────────────────────────────────────────
//  领域类型
// ─────────────────────────────────────────────

// BehaviorType 鲸鱼行为分类
type BehaviorType string

const (
	BehaviorLargeBuy  BehaviorType = "large_buy"
	BehaviorLargeSell BehaviorType = "large_sell"
	BehaviorLPRemove  BehaviorType = "lp_remove"
	BehaviorLargeXfer BehaviorType = "large_transfer"
)

// BehaviorEvent 鲸鱼行为事件
type BehaviorEvent struct {
	Address    common.Address
	EventType  BehaviorType
	TokenAddr  common.Address
	ValueUSD   decimal.Decimal
	Confidence float64 // 0~1
	Source     comm.UnifiedEvent
}

// ─────────────────────────────────────────────
//  Config
// ─────────────────────────────────────────────

// Config 鲸鱼监控配置
type Config struct {
	// MinUSDValue 识别为鲸鱼行为的最小 USD 价值（例如 100000 = 10万U）
	MinUSDValue decimal.Decimal
	// WatchAddresses 如果非空，只监控这些地址
	WatchAddresses []common.Address
}

// DefaultConfig 默认配置（10万 U 起）
func DefaultConfig() Config {
	return Config{
		MinUSDValue: decimal.NewFromInt(100_000),
	}
}

// ─────────────────────────────────────────────
//  Service
// ─────────────────────────────────────────────

// Service 鲸鱼监控服务
type Service struct {
	cfg    Config
	logger *zap.Logger
	// watchSet 快速查找监控地址
	watchSet map[common.Address]struct{}
}

// NewService 创建鲸鱼监控服务
func NewService(cfg Config, logger *zap.Logger) *Service {
	ws := make(map[common.Address]struct{}, len(cfg.WatchAddresses))
	for _, addr := range cfg.WatchAddresses {
		ws[addr] = struct{}{}
	}
	return &Service{cfg: cfg, logger: logger, watchSet: ws}
}

// Process 消费单批事件，返回识别到的鲸鱼行为列表
func (s *Service) Process(events []comm.UnifiedEvent) []BehaviorEvent {
	var results []BehaviorEvent
	for _, ev := range events {
		if be, ok := s.analyze(ev); ok {
			results = append(results, be)
		}
	}
	return results
}

// Run 持续消费事件 channel，将鲸鱼行为写入输出 channel
// 调用方通过 cancel ctx 停止
func (s *Service) Run(ctx context.Context, in <-chan []comm.UnifiedEvent) <-chan BehaviorEvent {
	out := make(chan BehaviorEvent, 100)
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
				for _, be := range s.Process(events) {
					select {
					case out <- be:
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
//  内部分析逻辑
// ─────────────────────────────────────────────

func (s *Service) analyze(ev comm.UnifiedEvent) (BehaviorEvent, bool) {
	base := ev.GetBase()

	// 监控地址过滤
	if len(s.watchSet) > 0 {
		if _, watched := s.watchSet[base.From]; !watched {
			return BehaviorEvent{}, false
		}
	}

	// 计算总 USD 价值
	totalUSD := decimal.Zero
	for _, rv := range base.RealValues {
		totalUSD = totalUSD.Add(rv)
	}

	// 阈值过滤
	if totalUSD.LessThan(s.cfg.MinUSDValue) {
		return BehaviorEvent{}, false
	}

	// 行为分类
	behaviorType := classifyBehavior(ev)
	if behaviorType == "" {
		return BehaviorEvent{}, false
	}

	var tokenAddr common.Address
	if swap, ok := ev.GetDetail().(*comm.SwapData); ok {
		tokenAddr = swap.FromToken
	} else if xfer, ok := ev.GetDetail().(*comm.TransferData); ok {
		tokenAddr = xfer.Token
	}

	confidence := calcConfidence(totalUSD, s.cfg.MinUSDValue)

	s.logger.Info("whale behavior detected",
		zap.String("address", base.From.Hex()),
		zap.String("type", string(behaviorType)),
		zap.String("usd_value", totalUSD.String()),
		zap.Float64("confidence", confidence),
	)

	return BehaviorEvent{
		Address:    base.From,
		EventType:  behaviorType,
		TokenAddr:  tokenAddr,
		ValueUSD:   totalUSD,
		Confidence: confidence,
		Source:     ev,
	}, true
}

func classifyBehavior(ev comm.UnifiedEvent) BehaviorType {
	switch ev.GetEventType() {
	case comm.EventMethodSwap:
		// 简化：有 Amount0In > 0 视为买入，否则卖出
		// 实际生产中需要结合代币价格方向判断
		return BehaviorLargeBuy
	case comm.EventMethodTransfer:
		return BehaviorLargeXfer
	}
	return ""
}

// calcConfidence 根据超出阈值倍数估算置信度（上限 1.0）
func calcConfidence(value, threshold decimal.Decimal) float64 {
	if threshold.IsZero() {
		return 1.0
	}
	ratio, _ := value.Div(threshold).Float64()
	// log scale：10x 阈值 → ~0.9，100x → ~1.0
	conf := 1.0 - 1.0/float64(big.NewInt(int64(ratio)+1).BitLen())
	if conf > 1.0 {
		return 1.0
	}
	return conf
}
