package comm

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ------------------命令模式------------------------------
// 每个具体实现（UniswapV2）用命令模式解析事件
type ProtocolImplParser interface {
	HandleEvent(eventsig EventSig, log types.Log, metadata EventMetadata) (UnifiedEvent, error)
	//ListEvents() []EventSig
	SetEvents(ms []EventMethod)
	ListEvents() map[EventMethod]common.Address
}

// 命令类型：事件解析函数
type EventParserFunc func(log types.Log, metadata EventMetadata) (UnifiedEvent, error)

// // 事件路由表（Invoker内部可用map管理命令）
// var EventParseConfig map[EventSig]EventParserFunc // 解析函数

// 调用者（Invoker）
type EventParseInvoker struct {
	Name       ProtocolImplName
	Handlers   map[EventSig]EventParserFunc
	NeedParser []MethodName
}

func NewEventParseInvoker() *EventParseInvoker {
	return &EventParseInvoker{
		Handlers: make(map[EventSig]EventParserFunc),
	}
}

// 注册事件解析命令
func (e *EventParseInvoker) Register(handlers map[EventSig]EventParserFunc) {
	e.Handlers = handlers
}

//======================================================

func (e *EventParseInvoker) HandleEvent(eventsig EventSig, log types.Log, metadata EventMetadata) (UnifiedEvent, error) {
	if handler, ok := e.Handlers[eventsig]; ok {
		return handler(log, metadata)
	}
	return nil, fmt.Errorf("unsupported event %s", eventsig)
}

func (e *EventParseInvoker) ListEvents() []EventSig {
	keys := make([]EventSig, 0, len(e.Handlers))
	for k := range e.Handlers {
		keys = append(keys, k)
	}
	return keys
}

func (e *EventParseInvoker) SetEvents(ms []EventMethod) {

}
