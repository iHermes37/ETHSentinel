package Scanner

import (
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// 定义责任链节点接口
type EventHandlerNode interface {
	Handle(log types.Log, metadata ParserEngineCommon.EventMetadata) (ParserEngineCommon.UnifiedEvent, bool)
	SetNext(next EventHandlerNode)
}

// 实现 事件解析节点（具体实现层）
type EventParserNode struct {
	parser ParserEngineCommon.ProtocolImplParser
	next   EventHandlerNode
}

// 对每笔交易 tranreceipt.Logs 遍历责任链，第一个匹配的解析器返回结果即可。
func (n *EventParserNode) Handle(log types.Log, metadata ParserEngineCommon.EventMetadata) (ParserEngineCommon.UnifiedEvent, bool) {
	// 遍历所有事件
	for _, sig := range n.parser.ListEvents() {
		if log.Topics[0] == common.Hash(sig) {
			ev, err := n.parser.HandleEvent(sig, log, metadata)
			if err == nil {
				return ev, true
			}
		}
	}
	// 没有匹配则传递到下一节点
	if n.next != nil {
		return n.next.Handle(log, metadata)
	}
	return nil, false
}

func (n *EventParserNode) SetNext(next EventHandlerNode) {
	n.next = next
}

// selectedProtocols 是用户选择的协议类型及具体实现，例如：
//
//	map[ProtocolType][]ProtocolImpl{
//	    DEX: {UniswapV2, SushiSwap},
//	    ERC20: {ERC20Std},
//	}
//
// 构建责任链
func BuildParserChain(pm *ParserEngineCommon.ProtocolManager, selectedProtocols map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl) EventHandlerNode {
	var head, prev EventHandlerNode

	for protoType, impls := range selectedProtocols {
		ptParser, _ := pm.GetProtocol(protoType)
		for _, impl := range impls {
			parser, _ := ptParser.GetImplementation(impl)
			node := &EventParserNode{parser: parser}
			if head == nil {
				head = node
			}
			if prev != nil {
				prev.SetNext(node)
			}
			prev = node
		}
	}
	return head
}

// 使用责任链解析交易/日志
// chain := BuildParserChain(protocolManager, selectedProtocols)
// for _, log := range receipt.Logs {
// 	if ev, ok := chain.Handle(*log, EventMetadata{BlockNumber: receipt.BlockNumber}); ok {
// 		fmt.Println("解析成功:", ev)
// 	} else {
// 		fmt.Println("未匹配事件", log.Topics[0])
// 	}
// }
