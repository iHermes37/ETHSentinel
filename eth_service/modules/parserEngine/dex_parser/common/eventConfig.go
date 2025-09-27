package dexcommon

import (
	abligens "github.com/Crypto-ChainSentinel/modules/parserEngine/dex_parser/abigens"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventParserFunc func(log types.Log, metadata EventMetadata, filterer *abligens.UniswappairFilterer) (UnifiedEvent, error)

var EventParseConfig map[EventSig]EventParserFunc // 解析函数
