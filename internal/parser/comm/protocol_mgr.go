package comm

import (
	"fmt"
)

// 全局统一管理器
type ProtocolManager struct {
	CurProtocols map[ProtocolImplName]ProtocolImplParser
	AllProtocols map[ProtocolTypeName]ProtocolTypeParser
}

func NewProtocolManager() *ProtocolManager {
	return &ProtocolManager{AllProtocols: make(map[ProtocolTypeName]ProtocolTypeParser)}
}

func (m *ProtocolManager) RegisterProtocol(protocolType ProtocolTypeName, parser ProtocolTypeParser) {
	m.AllProtocols[protocolType] = parser
}

func (m *ProtocolManager) SetCurProtocol(curs map[ProtocolTypeName]ProtocolTypeParser) {
	//m.CurProtocols = curs
}

func (m *ProtocolManager) GetProtocol(protocolType ProtocolTypeName) (ProtocolTypeParser, error) {
	parser, ok := m.AllProtocols[protocolType]
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
