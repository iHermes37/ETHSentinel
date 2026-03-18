// Package scanner — 责任链（Chain of Responsibility）
//
// 重构要点：
//   - 修复原 parse_chain.go 中 ListEvents() 返回类型不匹配的问题
//   - 责任链节点直接持有 ProtocolImplParser，逐签名匹配
//   - BuildChain 接受 *comm.ActiveParsers（由 Engine.BuildActive 生成）
package scanner

import (
	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ─────────────────────────────────────────────
//  责任链节点接口
// ─────────────────────────────────────────────

// ChainNode 责任链节点
type ChainNode interface {
	Handle(log types.Log, meta comm.EventMetadata) (comm.UnifiedEvent, bool)
	SetNext(next ChainNode)
}

// ─────────────────────────────────────────────
//  EventParserNode — 具体责任链节点
// ─────────────────────────────────────────────

// EventParserNode 持有一个 ProtocolImplParser，
// 遍历其注册的事件签名，匹配则解析，否则向下传递。
type EventParserNode struct {
	parser   comm.ProtocolImplParser
	sigIndex map[common.Hash]struct{} // 快速 O(1) 签名查找
	next     ChainNode
}

func newEventParserNode(parser comm.ProtocolImplParser) *EventParserNode {
	sigs := parser.ListEventSigs()
	idx := make(map[common.Hash]struct{}, len(sigs))
	for _, sig := range sigs {
		idx[sig] = struct{}{}
	}
	return &EventParserNode{parser: parser, sigIndex: idx}
}

// Handle 若 log.Topics[0] 命中当前节点的签名表，则解析并返回；否则向下传递。
func (n *EventParserNode) Handle(log types.Log, meta comm.EventMetadata) (comm.UnifiedEvent, bool) {
	if len(log.Topics) == 0 {
		if n.next != nil {
			return n.next.Handle(log, meta)
		}
		return nil, false
	}

	topic0 := log.Topics[0]
	if _, ok := n.sigIndex[topic0]; ok {
		ev, err := n.parser.HandleEvent(comm.EventSig(topic0), log, meta)
		if err == nil {
			return ev, true
		}
		// 解析失败继续向下（可能是合约地址不在监控列表）
	}

	if n.next != nil {
		return n.next.Handle(log, meta)
	}
	return nil, false
}

func (n *EventParserNode) SetNext(next ChainNode) { n.next = next }

// ─────────────────────────────────────────────
//  BuildChain
// ─────────────────────────────────────────────

// BuildChain 根据 ActiveParsers 构建责任链并返回链头。
// ActiveParsers.Impls 的顺序决定链的优先级（先注册先匹配）。
func BuildChain(active *comm.ActiveParsers) ChainNode {
	if active == nil || len(active.Impls) == 0 {
		return nil
	}

	var head, prev ChainNode
	for _, impl := range active.Impls {
		node := newEventParserNode(impl)
		if head == nil {
			head = node
		}
		if prev != nil {
			prev.SetNext(node)
		}
		prev = node
	}
	return head
}
