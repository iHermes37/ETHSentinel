// Package chain 定义多链抽象层，让上层代码不感知具体链。
// 所有 EVM 兼容链（ETH/BSC/Polygon/Arbitrum）共用同一套接口。
package chain

// ChainID 链 ID 类型
type ChainID uint64

// 预定义主流 EVM 链 ID
const (
	ChainETH      ChainID = 1
	ChainBSC      ChainID = 56
	ChainPolygon  ChainID = 137
	ChainArbitrum ChainID = 42161
	ChainOptimism ChainID = 10
	ChainBase     ChainID = 8453
)

// Chain EVM 链抽象接口
type Chain interface {
	ID() ChainID
	Name() string    // "ethereum" / "bsc" / "polygon"
	Symbol() string  // "ETH" / "BNB" / "MATIC"
	IsEVM() bool
	// 默认节点配置（可被用户覆盖）
	DefaultRPCURL() string
	DefaultWSURL() string
}

// BaseChain 通用 EVM 链实现
type BaseChain struct {
	id         ChainID
	name       string
	symbol     string
	defaultRPC string
	defaultWS  string
}

func NewBaseChain(id ChainID, name, symbol, rpc, ws string) *BaseChain {
	return &BaseChain{id: id, name: name, symbol: symbol, defaultRPC: rpc, defaultWS: ws}
}

func (c *BaseChain) ID() ChainID        { return c.id }
func (c *BaseChain) Name() string       { return c.name }
func (c *BaseChain) Symbol() string     { return c.symbol }
func (c *BaseChain) IsEVM() bool        { return true }
func (c *BaseChain) DefaultRPCURL() string { return c.defaultRPC }
func (c *BaseChain) DefaultWSURL() string  { return c.defaultWS }
