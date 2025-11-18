package dexparser

import (
	"github.com/Crypto-ChainSentinel/modules/parse_engine/comm"
	"github.com/Crypto-ChainSentinel/modules/parse_engine/dex_parser/uniswap_v2"
)

type UniswapV2Registrar struct{}

func (UniswapV2Registrar) ProtocolName() comm.ProtocolType { return comm.DEX }

func (UniswapV2Registrar) Register(manager *comm.ProtocolImplManager) {
	invoker := comm.NewProtocolEventParseInvoker()
	for sig, parser := range uniswapv2.UniswapV2EventsConfig {
		invoker.Register(sig, parser)
	}
	manager.RegisterStrategy("UniswapV2", invoker)
}

// 显式暴露注册函数
func NewUniswapV2Registrar() comm.ProtocolRegistrar {
	return UniswapV2Registrar{}
}
