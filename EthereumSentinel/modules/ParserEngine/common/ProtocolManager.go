package common

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
)

// 每个具体实现（UniswapV2、ERC20等）用命令模式解析事件
type ProtocolImplParser interface {
	HandleEvent(eventsig EventSig, log types.Log, metadata EventMetadata) (UnifiedEvent, error)
	ListEvents() []EventSig
	ListAddr() map[ProtocolType]common.Address
}

// 协议类型层（DEX, ERC20, Lending等）管理
type ProtocolTypeParser interface {
	GetImplementation(name ProtocolImpl) (ProtocolImplParser, error)
	ListImplementations() []ProtocolImpl
}

// 全局统一管理器
type ProtocolManager struct {
	protocols map[ProtocolType]ProtocolTypeParser
}

func NewProtocolManager() *ProtocolManager {
	return &ProtocolManager{protocols: make(map[ProtocolType]ProtocolTypeParser)}
}

func (m *ProtocolManager) RegisterProtocol(protocolType ProtocolType, parser ProtocolTypeParser) {
	m.protocols[protocolType] = parser
}

func (m *ProtocolManager) GetProtocol(protocolType ProtocolType) (ProtocolTypeParser, error) {
	parser, ok := m.protocols[protocolType]
	if !ok {
		return nil, fmt.Errorf("unsupported protocol %v", protocolType)
	}
	return parser, nil
}

// func InitProtocolManager() *ProtocolManager {
// 	pm := NewProtocolManager()

// 	// --------- DEX 协议 ----------
// 	dexManager := NewProtocolImplManager()

// 	// 注册 UniswapV2 事件解析
// 	uniswapInvoker := NewProtocolEventParseInvoker()
// 	uniswapInvoker.Register(UniswapV2Swap, func(log types.Log, meta EventMetadata) (UnifiedEvent, error) {
// 		return fmt.Sprintf("UniswapV2 Swap parsed at block %d", meta.BlockNumber), nil
// 	})
// 	dexManager.RegisterStrategy(UniswapV2, uniswapInvoker)

// 	// 注册 SushiSwap 事件解析
// 	sushiInvoker := NewProtocolEventParseInvoker()
// 	sushiInvoker.Register(UniswapV2_Swap, func(log types.Log, meta EventMetadata) (UnifiedEvent, error) {
// 		return fmt.Sprintf("SushiSwap Swap parsed at block %d", meta.BlockNumber), nil
// 	})
// 	dexManager.RegisterStrategy(SushiSwap, sushiInvoker)

// 	pm.RegisterProtocol(DEX, dexManager)

// 	// --------- ERC20 协议 ----------
// 	erc20Manager := NewProtocolImplManager()

// 	erc20Invoker := NewProtocolEventParseInvoker()
// 	erc20Invoker.Register(ERC20_Transfer, func(log types.Log, meta EventMetadata) (UnifiedEvent, error) {
// 		return fmt.Sprintf("ERC20 Transfer parsed tx %s", meta.TxHash), nil
// 	})
// 	erc20Manager.RegisterStrategy(ERC20Std, erc20Invoker)

// 	pm.RegisterProtocol(ERC20, erc20Manager)

// 	return pm
// }
