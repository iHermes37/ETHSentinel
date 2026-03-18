// Package comm — EventParseInvoker（命令模式）
// 重构要点：
//   - 将 ListEvents() 返回 []EventSig（原返回 map，导致 parse_chain.go 里 for-range 类型冲突）
//   - SetEvents 支持按语义名称过滤（不再是空实现）
//   - Register 改为接受 map，与原 registrar.go 兼容
package comm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

// ─────────────────────────────────────────────
//  接口定义
// ─────────────────────────────────────────────

// ProtocolTypeParser 协议大类管理器接口（DEX / Token / Lending …）
type ProtocolTypeParser interface {
	GetImpl(name ProtocolImpl) (ProtocolImplParser, error)
	ListImpls() []ProtocolImpl
}

// ProtocolImplParser 单个协议实现的解析接口（命令模式 Invoker）
type ProtocolImplParser interface {
	// HandleEvent 解析一条 log
	HandleEvent(sig EventSig, log types.Log, meta EventMetadata) (UnifiedEvent, error)
	// ListEventSigs 返回该实现所有已注册的事件签名
	ListEventSigs() []EventSig
	// SetFilter 设置只解析哪些语义事件（空 = 全部）
	SetFilter(methods []EventMethod)
}

// ─────────────────────────────────────────────
//  EventParserFunc — 命令类型
// ─────────────────────────────────────────────

// EventParserFunc 单个事件的解析函数签名
type EventParserFunc func(log types.Log, meta EventMetadata) (UnifiedEvent, error)

// ─────────────────────────────────────────────
//  EventParseInvoker — 命令模式 Invoker
// ─────────────────────────────────────────────

// EventParseInvoker 持有一组 EventSig → EventParserFunc 的路由表，
// 并根据 log.Topics[0] 分发到对应解析函数。
type EventParseInvoker struct {
	Name        ProtocolImpl
	allHandlers map[EventSig]EventParserFunc // 全量注册
	active      map[EventSig]EventParserFunc // 过滤后生效的
	// sigToMethod 用于 SetFilter 按语义名称过滤
	sigToMethod map[EventSig]EventMethod
}

// NewEventParseInvoker 创建空的 Invoker
func NewEventParseInvoker(name ProtocolImpl) *EventParseInvoker {
	return &EventParseInvoker{
		Name:        name,
		allHandlers: make(map[EventSig]EventParserFunc),
		active:      make(map[EventSig]EventParserFunc),
		sigToMethod: make(map[EventSig]EventMethod),
	}
}

// Register 批量注册事件签名 → 解析函数映射
func (e *EventParseInvoker) Register(handlers map[EventSig]EventParserFunc) {
	for sig, fn := range handlers {
		e.allHandlers[sig] = fn
		e.active[sig] = fn
	}
}

// RegisterOne 单条注册，附带语义方法名（用于 SetFilter）
func (e *EventParseInvoker) RegisterOne(sig EventSig, method EventMethod, fn EventParserFunc) {
	e.allHandlers[sig] = fn
	e.active[sig] = fn
	e.sigToMethod[sig] = method
}

// SetFilter 过滤只解析指定语义事件；空列表 = 解析全部
func (e *EventParseInvoker) SetFilter(methods []EventMethod) {
	if len(methods) == 0 {
		// 恢复全部
		for sig, fn := range e.allHandlers {
			e.active[sig] = fn
		}
		return
	}
	wantSet := make(map[EventMethod]struct{}, len(methods))
	for _, m := range methods {
		wantSet[m] = struct{}{}
	}
	e.active = make(map[EventSig]EventParserFunc)
	for sig, fn := range e.allHandlers {
		if method, ok := e.sigToMethod[sig]; ok {
			if _, want := wantSet[method]; want {
				e.active[sig] = fn
			}
		}
	}
}

// HandleEvent 根据事件签名路由到对应解析函数
func (e *EventParseInvoker) HandleEvent(sig EventSig, log types.Log, meta EventMetadata) (UnifiedEvent, error) {
	fn, ok := e.active[sig]
	if !ok {
		return nil, fmt.Errorf("invoker[%s]: no handler for sig %s", e.Name, sig.Hex())
	}
	return fn(log, meta)
}

// ListEventSigs 返回当前生效的事件签名列表
func (e *EventParseInvoker) ListEventSigs() []EventSig {
	sigs := make([]EventSig, 0, len(e.active))
	for sig := range e.active {
		sigs = append(sigs, sig)
	}
	return sigs
}
