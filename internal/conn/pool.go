// Package conn — 多链连接池
// 支持同时管理多条 EVM 链的连接，按 ChainID 索引
package conn

import (
	"context"
	"fmt"
	"sync"

	"github.com/ETHSentinel/internal/chain"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// ChainConn 单条链的连接组
type ChainConn struct {
	RPC *ethclient.Client
	WS  *ethclient.Client
}

// MultiChainPool 多链连接池
type MultiChainPool struct {
	mu      sync.RWMutex
	conns   map[chain.ChainID]*ChainConn
	configs map[chain.ChainID]*NodeConfig
	logger  *zap.Logger
}

// NewMultiChainPool 创建多链连接池
func NewMultiChainPool(logger *zap.Logger) *MultiChainPool {
	return &MultiChainPool{
		conns:   make(map[chain.ChainID]*ChainConn),
		configs: make(map[chain.ChainID]*NodeConfig),
		logger:  logger,
	}
}

// RegisterChain 注册一条链的连接配置
// 如果不传 cfg，使用链的默认节点地址
func (p *MultiChainPool) RegisterChain(c chain.Chain, cfg ...*NodeConfig) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var nodeCfg *NodeConfig
	if len(cfg) > 0 && cfg[0] != nil {
		nodeCfg = cfg[0]
	} else {
		nodeCfg = &NodeConfig{
			Name:   c.Name(),
			RPCURL: c.DefaultRPCURL(),
			WSURL:  c.DefaultWSURL(),
		}
	}
	p.configs[c.ID()] = nodeCfg
}

// GetRPC 获取指定链的 RPC 连接（懒加载）
func (p *MultiChainPool) GetRPC(ctx context.Context, chainID chain.ChainID) (*ethclient.Client, error) {
	p.mu.RLock()
	if cc, ok := p.conns[chainID]; ok && cc.RPC != nil {
		p.mu.RUnlock()
		return cc.RPC, nil
	}
	p.mu.RUnlock()

	return p.dial(ctx, chainID, MethodRPC)
}

// GetWS 获取指定链的 WS 连接（懒加载）
func (p *MultiChainPool) GetWS(ctx context.Context, chainID chain.ChainID) (*ethclient.Client, error) {
	p.mu.RLock()
	if cc, ok := p.conns[chainID]; ok && cc.WS != nil {
		p.mu.RUnlock()
		return cc.WS, nil
	}
	p.mu.RUnlock()

	return p.dial(ctx, chainID, MethodWS)
}

func (p *MultiChainPool) dial(ctx context.Context, chainID chain.ChainID, method Method) (*ethclient.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	cfg, ok := p.configs[chainID]
	if !ok {
		return nil, fmt.Errorf("conn: chain %d not registered", chainID)
	}

	client, err := dial(ctx, cfg, method)
	if err != nil {
		return nil, fmt.Errorf("conn: dial chain %d/%s: %w", chainID, method, err)
	}

	if _, ok := p.conns[chainID]; !ok {
		p.conns[chainID] = &ChainConn{}
	}
	if method == MethodRPC {
		p.conns[chainID].RPC = client
	} else {
		p.conns[chainID].WS = client
	}

	p.logger.Info("chain connected",
		zap.Uint64("chain_id", uint64(chainID)),
		zap.String("method", string(method)),
	)
	return client, nil
}

// Close 关闭所有连接
func (p *MultiChainPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, cc := range p.conns {
		if cc.RPC != nil {
			cc.RPC.Close()
		}
		if cc.WS != nil {
			cc.WS.Close()
		}
	}
}
