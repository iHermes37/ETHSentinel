package protocols

import (
	dexcommon "github.com/CryptoQuantX/chain_monitor/modules/parserEngine/dex_parser/common"
	"github.com/CryptoQuantX/chain_monitor/modules/parserEngine/dex_parser/protocols/UniswapV2Protocol"
	"github.com/ethereum/go-ethereum/common"
	"strings"
	"sync"
)

// 通过字符串解析协议
func GetProtocolAddr(s string) []common.Address {
	switch strings.ToLower(s) {
	case "uniswapv2":
		return UniswapV2Protocol.UniswapV2_Address
	//case "sushiswap":
	//	return Sushiswap, nil
	//case "curve":
	//	return Curve, nil
	default:
		return nil
	}
}

// 封装每个协议的配置
type ProtocolParsers struct {
	Addr    []common.Address
	Configs []dexcommon.EventParseConfig
}

var (
	DEXParseConfigManager map[dexcommon.Protocol]ProtocolParsers
	once                  sync.Once
)

func InitEventConfig() {
	once.Do(func() {
		DEXParseConfigManager = make(map[dexcommon.Protocol]ProtocolParsers)
		// UniswapV2Protocol
		DEXParseConfigManager[dexcommon.UniswapV2] = ProtocolParsers{
			Addr:    GetProtocolAddr("UniswapV2"),
			Configs: UniswapV2Protocol.Configs,
		}
	},
	)
}
