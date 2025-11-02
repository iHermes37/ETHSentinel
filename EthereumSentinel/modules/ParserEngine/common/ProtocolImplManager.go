package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

// -------------------------------策略模式------------------------------------------------

var _ ProtocolImplParser = (*ProtocolEventParseInvoker)(nil)

// 每个实现（UniswapV2、ERC20标准…）
type ProtocolImplManager struct {
	CurrentStrategy ProtocolImplParser
	impls           map[ProtocolImpl]ProtocolImplParser
}

func NewProtocolImplManager() *ProtocolImplManager {
	return &ProtocolImplManager{
		impls: make(map[ProtocolImpl]ProtocolImplParser),
	}
}

// 注册具体策略
func (m *ProtocolImplManager) RegisterStrategy(protocol ProtocolImpl, strategy ProtocolImplParser) {
	m.impls[protocol] = strategy
}

// 切换策略
func (m *ProtocolImplManager) SetStrategy(protocol ProtocolImpl) error {
	strategy, ok := m.impls[protocol]
	if !ok {
		return fmt.Errorf("unsupported protocol %v", protocol)
	}
	m.CurrentStrategy = strategy
	return nil
}

// 使用当前策略解析事件
func (m *ProtocolImplManager) HandleEvent(eventsig EventSig, log types.Log, metadata EventMetadata) (UnifiedEvent, error) {
	if m.CurrentStrategy == nil {
		return nil, fmt.Errorf("no strategy selected")
	}
	return m.CurrentStrategy.HandleEvent(eventsig, log, metadata)
}

func (p *ProtocolImplManager) GetImplementation(name ProtocolImpl) (ProtocolImplParser, error) {
	impl, ok := p.impls[name]
	if !ok {
		return nil, fmt.Errorf("unsupported implementation %s", name)
	}
	return impl, nil
}

func (p *ProtocolImplManager) ListImplementations() []ProtocolImpl {
	keys := make([]ProtocolImpl, 0, len(p.impls))
	for k := range p.impls {
		keys = append(keys, k)
	}
	return keys
}

//	func InitDexParserConfig() {
//		once.Do(func() {
//			DEXParseManager = make(map[ProtocolType]ProtocolEventParseInvoker)
//				// ---------------UniswapV2Protocol-------------------------------
//				// 创建 UniswapV2Protocol 的 invoker
//				uniswapInvoker := NewProtocolEventParseInvoker()
//				// 遍历 UniswapV2Protocol 的事件注册
//				for sig, parser := range UniswapV2Protocol.UniswapV2EventsConfig {
//					uniswapInvoker.Register(sig, parser)
//				}
//				// 注册到全局 DEX 管理器
//				DEXParseManager[UniswapV2] = *uniswapInvoker
//			},
//			)
//		}
// func InitDexStrategyManager() *ProtocolImplManager {
// 	manager := NewDexStrategyManager()

// 	// UniswapV2
// 	uniswapInvoker := NewProtocolEventParseInvoker()
// 	for sig, parser := range UniswapV2Protocol.UniswapV2EventsConfig {
// 		uniswapInvoker.Register(sig, parser)
// 	}
// 	manager.RegisterStrategy(UniswapV2, uniswapInvoker)

// 	// TODO: 注册其他DEX，比如 SushiSwap, Balancer...
// 	// manager.RegisterStrategy(SushiSwap, sushiInvoker)

// 	return manager
// }

// manager := InitDexStrategyManager()

// // 切换到 UniswapV2 策略
// manager.SetStrategy(UniswapV2)

// // 解析 Swap 事件
// event, err := manager.HandleEvent(UniswapV2_Swap, log, metadata)
// if err != nil {
// 	fmt.Println("解析失败:", err)
// } else {
// 	fmt.Println("事件解析成功:", event)
// }

// // 切换到其他DEX策略也同样简单
// // manager.SetStrategy(SushiSwap)
