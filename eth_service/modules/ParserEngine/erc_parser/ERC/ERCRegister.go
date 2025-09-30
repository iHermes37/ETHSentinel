package ERC

import (
	"github.com/Crypto-ChainSentinel/modules/ParserEngine/erc_parser/erccommon"
	"github.com/ethereum/go-ethereum/common"
	"sync"
)

type ProtocolParser struct {
	protocolName erccommon.Protocol
	TokenAddrs   map[common.Address]erccommon.Protocol
	Configs      map[erccommon.MethodName]erccommon.EventParserFunc
}

var (
	ERCParseConfigManager map[erccommon.Protocol]ProtocolParser
	once                  sync.Once
)

func InitERCParser() {
	once.Do(func() {
		ERCParseConfigManager = make(map[erccommon.Protocol]ProtocolParser)
		// UniswapV2Protocol
		ERCParseConfigManager[erccommon.ERC20] = ProtocolParser{
			erccommon.ERC20,
			TokenAddr,
			ERC20EventsConfig,
		}
	})
}
