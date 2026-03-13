package main

import "github.com/Crypto-ChainSentinel/internal/parser/comm"

//type ProtocolRegistrar interface {
//	ProtocolName() ProtocolType
//	Register(manager *ProtocolImplManager)
//}
//
//type ProtocolRegistry struct {
//	registrars []ProtocolRegistrar
//}
//
////==========================================
//
//var (
//	instance *ProtocolRegistry
//	once     sync.Once
//)
//
//// 单例获取
//func GetRegistry() *ProtocolRegistry {
//	once.Do(func() {
//		instance = &ProtocolRegistry{}
//	})
//	return instance
//}

//=================================================\

// 注册协议
func (r *ProtocolRegistry) RegisterProtocol(pr ProtocolRegistrar) {
	r.registrars = append(r.registrars, pr)
}

func (r *ProtocolRegistry) InitProtocolManager() *comm.ProtocolManager {
	pm := comm.NewProtocolManager()
	for _, reg := range r.registrars {
		implManager := comm.NewProtocolImplManager()
		reg.Register(implManager)
		pm.RegisterProtocol(reg.ProtocolName(), implManager)
	}
	return pm
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
// 	// -------- ERC20 协议 ----------
// 	erc20Manager := NewProtocolImplManager()
// 	erc20Invoker := NewProtocolEventParseInvoker()
// 	erc20Invoker.Register(ERC20_Transfer, func(log types.Log, meta EventMetadata) (UnifiedEvent, error) {
// 		return fmt.Sprintf("ERC20 Transfer parsed tx %s", meta.TxHash), nil
// 	})
// 	erc20Manager.RegisterStrategy(ERC20Std, erc20Invoker)
// 	pm.RegisterProtocol(ERC20, erc20Manager)
// 	return pm
// }
