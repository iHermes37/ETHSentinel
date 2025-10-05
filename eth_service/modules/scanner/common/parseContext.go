package common

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 把所有 log 解析可能需要的上下文，统一封装成一个结构体 `ParseContext`，每个解析器只取自己关心的字段。
type ParseContext struct {
	Log      types.Log
	Metadata EventMetadata
	Client   *ethclient.Client
	Filterer interface{}            // 协议专用 filterer，可以做类型断言
	Extra    map[string]interface{} // 灵活扩展
}

// 所有解析器都实现这个接口
type Parser interface {
	Parse(ctx ParseContext) (UnifiedEvent, error)
}
