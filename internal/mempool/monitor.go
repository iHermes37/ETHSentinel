// Package mempool 提供以太坊 Mempool（交易池）监控能力。
// 基于你的 radar.go 重构，支持多链、依赖注入、过滤器。
package mempool

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap"
)

// PendingTx 封装 pending 交易及附加信息
type PendingTx struct {
	Tx        *types.Transaction
	From      common.Address
	SeenAt    time.Time
	ChainID   uint64
	GasPrice  *big.Int
	GasTipCap *big.Int // EIP-1559
}

// Monitor Mempool 监控器
type Monitor struct {
	wsURL   string          // WS 节点（必须是 WS，gethclient 订阅需要）
	client  *ethclient.Client
	filters []FilterFunc
	logger  *zap.Logger
	chainID uint64
}

// FilterFunc 过滤函数，返回 true 表示保留该交易
type FilterFunc func(tx *types.Transaction) bool

// NewMonitor 创建 Mempool 监控器
func NewMonitor(wsURL string, client *ethclient.Client, chainID uint64, logger *zap.Logger) *Monitor {
	return &Monitor{
		wsURL:   wsURL,
		client:  client,
		chainID: chainID,
		logger:  logger,
	}
}

// WithFilter 添加过滤器（链式调用）
func (m *Monitor) WithFilter(f FilterFunc) *Monitor {
	m.filters = append(m.filters, f)
	return m
}

// Subscribe 订阅 pending 交易，通过 channel 流式输出
// 底层使用 gethclient.SubscribePendingTransactions（需要 WS 连接）
func (m *Monitor) Subscribe(ctx context.Context) (<-chan *PendingTx, error) {
	// gethclient 需要单独建立 rpc 连接
	rpcClient, err := rpc.DialWebsocket(ctx, m.wsURL, "")
	if err != nil {
		return nil, err
	}
	gc := gethclient.New(rpcClient)

	hashes := make(chan common.Hash, 200)
	sub, err := gc.SubscribePendingTransactions(ctx, hashes)
	if err != nil {
		rpcClient.Close()
		return nil, err
	}

	out := make(chan *PendingTx, 100)

	go func() {
		defer close(out)
		defer sub.Unsubscribe()
		defer rpcClient.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case err := <-sub.Err():
				m.logger.Warn("mempool subscription error", zap.Error(err))
				return
			case hash := <-hashes:
				tx, isPending, err := m.client.TransactionByHash(ctx, hash)
				if err != nil || !isPending {
					continue
				}

				// 应用过滤器
				if !m.passFilters(tx) {
					continue
				}

				// 获取发送方地址
				from, _ := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)

				pendingTx := &PendingTx{
					Tx:       tx,
					From:     from,
					SeenAt:   time.Now(),
					ChainID:  m.chainID,
					GasPrice: tx.GasPrice(),
				}
				if tx.Type() == types.DynamicFeeTxType {
					pendingTx.GasTipCap = tx.GasTipCap()
				}

				select {
				case out <- pendingTx:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	m.logger.Info("mempool monitor started", zap.Uint64("chain_id", m.chainID))
	return out, nil
}

func (m *Monitor) passFilters(tx *types.Transaction) bool {
	for _, f := range m.filters {
		if !f(tx) {
			return false
		}
	}
	return true
}

// MonitorWhaleRefTx 监控鲸鱼相关 pending 交易（保留原有接口）
func (m *Monitor) MonitorWhaleRefTx(ctx context.Context) {
	ch, err := m.Subscribe(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for ptx := range ch {
		m.logger.Info("pending tx",
			zap.String("hash", ptx.Tx.Hash().Hex()),
			zap.String("from", ptx.From.Hex()),
			zap.String("gas_price", ptx.GasPrice.String()),
		)
	}
}
