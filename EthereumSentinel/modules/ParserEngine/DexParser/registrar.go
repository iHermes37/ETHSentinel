package dexparser

import (
	"github.com/Crypto-ChainSentinel/modules/ParserEngine/DexParser/uniswapv2"
	"github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
)

type UniswapV2Registrar struct{}

func (UniswapV2Registrar) ProtocolName() common.ProtocolType { return common.DEX }

func (UniswapV2Registrar) Register(manager *common.ProtocolImplManager) {
	invoker := common.NewProtocolEventParseInvoker()
	for sig, parser := range uniswapv2.UniswapV2EventsConfig {
		invoker.Register(sig, parser)
	}
	manager.RegisterStrategy("UniswapV2", invoker)
}

// 显式暴露注册函数
func NewUniswapV2Registrar() common.ProtocolRegistrar {
	return UniswapV2Registrar{}
}
