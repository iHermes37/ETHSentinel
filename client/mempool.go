package sentinel

import (
	"context"
	"fmt"
	"math/big"

	sentinelv1 "github.com/ETHSentinel/gen/sentinel/v1"
	"github.com/ETHSentinel/internal/chain"
	"github.com/ETHSentinel/internal/mempool"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// PendingTx 对外暴露的 pending 交易类型
type PendingTx = mempool.PendingTx

// MempoolClient Mempool 监控客户端
type MempoolClient struct {
	// 嵌入模式
	wsClient *ethclient.Client
	wsURL    string
	// 远程模式
	mempoolStub sentinelv1.MempoolServiceClient
	// 公共
	chainID uint64
	logger  *zap.Logger
}

// FilterByMinValueWei 过滤：最小 ETH value（wei）
func FilterByMinValueWei(wei *big.Int) func(*mempool.Monitor) {
	return func(m *mempool.Monitor) {
		m.WithFilter(mempool.FilterByMinValue(wei))
	}
}

// FilterByMinValueETH 过滤：最小 ETH value（ETH 单位）
func FilterByMinValueETH(eth float64) func(*mempool.Monitor) {
	f := new(big.Float).Mul(big.NewFloat(eth), big.NewFloat(1e18))
	wei, _ := f.Int(nil)
	return func(m *mempool.Monitor) {
		m.WithFilter(mempool.FilterByMinValue(wei))
	}
}

// FilterByMinGas 过滤：最小 GasPrice（Gwei）
func FilterByMinGas(gwei int64) func(*mempool.Monitor) {
	return func(m *mempool.Monitor) {
		m.WithFilter(mempool.FilterByMinGasPrice(gwei))
	}
}

// FilterByMethod 过滤：方法签名（如 "0x38ed1739"）
func FilterByMethod(sigs ...string) func(*mempool.Monitor) {
	return func(m *mempool.Monitor) {
		m.WithFilter(mempool.FilterByMethodSig(sigs...))
	}
}

// Subscribe 订阅 pending 交易
func (m *MempoolClient) Subscribe(ctx context.Context, filters ...func(*mempool.Monitor)) (<-chan *PendingTx, error) {
	// 远程模式
	if m.mempoolStub != nil {
		return m.subscribeRemote(ctx)
	}
	// 嵌入模式
	if m.wsClient == nil {
		return nil, fmt.Errorf("mempool: ws client not initialized (use WithWSURL)")
	}
	mon := mempool.NewMonitor(m.wsURL, m.wsClient, m.chainID, m.logger)
	for _, f := range filters {
		f(mon)
	}
	return mon.Subscribe(ctx)
}

func (m *MempoolClient) subscribeRemote(ctx context.Context) (<-chan *PendingTx, error) {
	stream, err := m.mempoolStub.SubscribePending(ctx, &sentinelv1.MempoolSubscribeRequest{
		ChainId: fmt.Sprintf("%d", m.chainID),
	})
	if err != nil {
		return nil, err
	}
	out := make(chan *PendingTx, 100)
	go func() {
		defer close(out)
		for {
			_, err := stream.Recv()
			if err != nil {
				return
			}
			// TODO: 转换 proto PendingTxEvent → PendingTx
		}
	}()
	return out, nil
}

func newMempoolClient(wsURL string, wsClient *ethclient.Client, c chain.Chain, logger *zap.Logger) *MempoolClient {
	return &MempoolClient{
		wsClient: wsClient,
		wsURL:    wsURL,
		chainID:  uint64(c.ID()),
		logger:   logger,
	}
}
