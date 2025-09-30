package dexcommon

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type EventParserFunc func(log types.Log, metadata EventMetadata) (UnifiedEvent, error)

var EventParseConfig map[EventSig]EventParserFunc // 解析函数
