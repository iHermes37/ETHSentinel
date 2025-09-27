package protocols

import (
	dexcommon "github.com/Crypto-ChainSentinel/modules/parserEngine/dex_parser/common"
	"github.com/Crypto-ChainSentinel/modules/parserEngine/dex_parser/protocols/UniswapV2Protocol"
	"github.com/ethereum/go-ethereum/common"
	"sync"
)

// 封装每个协议的配置(该协议的所有pair地址【找到对应的协议】+该协议所有的解析器配置)
type ProtocolParser struct {
	protocolName  dexcommon.Protocol
	ContractAddrs map[common.Address]dexcommon.Protocol
	Configs       map[dexcommon.EventSig]dexcommon.EventParserFunc
}

var (
	DEXParseConfigManager map[dexcommon.Protocol]ProtocolParser
	once                  sync.Once
)

func InitEventConfig() {
	once.Do(func() {
		DEXParseConfigManager = make(map[dexcommon.Protocol]ProtocolParser)
		// UniswapV2Protocol
		DEXParseConfigManager[dexcommon.UniswapV2] = ProtocolParser{
			dexcommon.UniswapV2,
			UniswapV2Protocol.ContractAddress,
			UniswapV2Protocol.UniswapV2EventsConfig,
		}
	},
	)
}
