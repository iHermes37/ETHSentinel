package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

// ------------------命令模式------------------------------

// 命令类型：事件解析函数
type EventParserFunc func(log types.Log, metadata EventMetadata) (UnifiedEvent, error)

// // 事件路由表（Invoker内部可用map管理命令）
// var EventParseConfig map[EventSig]EventParserFunc // 解析函数

// 调用者（Invoker）
type ProtocolEventParseInvoker struct {
	// ProtocolName  ProtocolType
	// ProtocolAddrs map[common.Address]ProtocolType
	Handlers map[EventSig]EventParserFunc
}

func NewProtocolEventParseInvoker() *ProtocolEventParseInvoker {
	return &ProtocolEventParseInvoker{
		Handlers: make(map[EventSig]EventParserFunc),
	}
}

// 注册事件命令
func (e *ProtocolEventParseInvoker) Register(eventsig EventSig, handler EventParserFunc) {
	e.Handlers[eventsig] = handler
}

// 执行事件命令
func (e *ProtocolEventParseInvoker) HandleEvent(eventsig EventSig, log types.Log, metadata EventMetadata) (UnifiedEvent, error) {
	if handler, ok := e.Handlers[eventsig]; ok {
		return handler(log, metadata)
	}
	return nil, fmt.Errorf("unsupported event %s", eventsig)
}

func (e *ProtocolEventParseInvoker) ListEvents() []EventSig {
	keys := make([]EventSig, 0, len(e.Handlers))
	for k := range e.Handlers {
		keys = append(keys, k)
	}
	return keys
}
