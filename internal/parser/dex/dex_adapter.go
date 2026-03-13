package dex

import (
	"github.com/Crypto-ChainSentinel/internal/parser/comm"
	uniswap_v2 "github.com/Crypto-ChainSentinel/internal/parser/dex/uniswap_v2"
)

func DexAdapter() comm.ProtocolTypeParser {
	impl_mgr := comm.NewProtocolImplManager()
	impl_mgr.RegisterStrategy(comm.UniswapV2, uniswap_v2.Register())

	var protocolTypeParser comm.ProtocolTypeParser
	protocolTypeParser = impl_mgr

	return protocolTypeParser
}
