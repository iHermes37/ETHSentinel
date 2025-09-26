package dexcommon

import (
	abligens "github.com/CryptoQuantX/chain_monitor/modules/parserEngine/dex_parser/abigens"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventParserFunc func(log types.Log, metadata EventMetadata, filterer *abligens.UniswappairFilterer) (dexcommon.UnifiedEvent, error)

type EventParseConfig struct {
	ContractAddress common.Address // 合约地址
	Protocol        Protocol
	EventType       EventType
	Parser          EventParserFunc // 解析函数
}

const (
	UniswapV2 Protocol = iota
)
