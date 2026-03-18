package chain

import (
	"fmt"
	"sync"
)

// Registry 全局链注册表
type Registry struct {
	mu     sync.RWMutex
	chains map[ChainID]Chain
}

var defaultRegistry = &Registry{
	chains: make(map[ChainID]Chain),
}

// Register 注册一条链
func Register(c Chain) {
	defaultRegistry.mu.Lock()
	defer defaultRegistry.mu.Unlock()
	defaultRegistry.chains[c.ID()] = c
}

// Get 根据 ChainID 获取链配置
func Get(id ChainID) (Chain, error) {
	defaultRegistry.mu.RLock()
	defer defaultRegistry.mu.RUnlock()
	c, ok := defaultRegistry.chains[id]
	if !ok {
		return nil, fmt.Errorf("chain: unsupported chain ID %d", id)
	}
	return c, nil
}

// MustGet 获取链配置，不存在则 panic（用于初始化阶段）
func MustGet(id ChainID) Chain {
	c, err := Get(id)
	if err != nil {
		panic(err)
	}
	return c
}

// All 返回所有已注册的链
func All() []Chain {
	defaultRegistry.mu.RLock()
	defer defaultRegistry.mu.RUnlock()
	list := make([]Chain, 0, len(defaultRegistry.chains))
	for _, c := range defaultRegistry.chains {
		list = append(list, c)
	}
	return list
}

func init() {
	// 启动时自动注册所有内置链
	Register(Ethereum())
	Register(BSC())
	Register(Polygon())
	Register(Arbitrum())
}
