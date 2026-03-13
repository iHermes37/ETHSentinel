package comm

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
)

// -------------------------------策略模式------------------------------------------------
//var ProtocolImplParser = (*ProtocolImplParser)(nil)

// 协议类型层（DEX,  Lending等）管理
type ProtocolTypeParser interface {
	GetImpl(name ProtocolImplName) (ProtocolImplParser, error)
	ListImpls() []ProtocolImplName
}

// 每个实现（UniswapV2、ERC20标准…）
type ProtocolImplManager struct {
	Cur   ProtocolImplParser
	Impls map[ProtocolImplName]ProtocolImplParser
}

func NewProtocolImplManager() *ProtocolImplManager {
	return &ProtocolImplManager{
		Impls: make(map[ProtocolImplName]ProtocolImplParser),
	}
}

//===================================================================================================

// 注册具体策略
func (m *ProtocolImplManager) RegisterStrategy(protocol ProtocolImplName, strategy ProtocolImplParser) {
	m.Impls[protocol] = strategy
}

// 切换策略
func (m *ProtocolImplManager) SetStrategy(protocol ProtocolImplName) error {
	strategy, ok := m.Impls[protocol]
	if !ok {
		return fmt.Errorf("unsupported protocol %v", protocol)
	}
	m.Cur = strategy
	return nil
}

// 使用当前策略解析事件
func (m *ProtocolImplManager) HandleEvent(eventsig EventSig, log types.Log, metadata EventMetadata) (UnifiedEvent, error) {
	if m.Cur == nil {
		return nil, fmt.Errorf("no strategy selected")
	}
	return m.Cur.HandleEvent(eventsig, log, metadata)
}

// ===============================
func (p *ProtocolImplManager) GetImpl(name ProtocolImplName) (ProtocolImplParser, error) {
	Impl, ok := p.Impls[name]
	if !ok {
		return nil, fmt.Errorf("unsupported Implementation %s", name)
	}
	return Impl, nil
}

func (p *ProtocolImplManager) ListImpls() []ProtocolImplName {
	keys := make([]ProtocolImplName, 0, len(p.Impls))
	for k := range p.Impls {
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
