package dexcommon

import (
	"github.com/ethereum/go-ethereum/common"
)

type Protocol int

// UnifiedEvent 接口
type UnifiedEvent interface {
	EventType() EventType //事件类型，比如 Swap、Deposit、Withdraw 等
	TxHash() common.Hash  //以太坊的 tx hash，用于唯一标识触发事件的交易
	BlockNumber() uint64  //区块高度

	SetSwapData(data SwapData) //解析出来的 swap 交易信息（from_token、to_token、amount 等）；
	SwapDataIsParsed() bool    //是否已经解析了 SwapData

	OuterIndex() int64         //log 在交易中的索引
	TransactionIndex() *uint64 //Ethereum 交易在区块中的索引

	Clone() UnifiedEvent

	//HandleMs() int64      // 事件处理耗时（毫秒）
	//SetHandleMs(ms int64) // 设置处理耗时
}

// EventMetadata 封装以太坊链上事件的基础信息，用于统一事件处理
type EventMetadata struct {
	EventTypeVal        EventType   // 事件类型，例如 Swap、Deposit、Withdraw 等
	ProtocolVal         Protocol    // 协议类型，例如 UniswapV2、Sushiswap 等
	TxHashVal           common.Hash // 交易哈希（唯一标识交易）
	BlockNumberVal      uint64      // 区块高度，用于排序和定位交易所在区块
	OuterIndexVal       int64       // log 在交易中的索引
	TransactionIndexVal *uint64     // 交易在区块中的索引位置，用于区分同区块内多笔交易
	SwapDataVal         SwapData    // swap 相关数据，例如 from_token、to_token、数量等
	SwapParsed          bool        // swap 数据是否已经解析完成
}

func (b *EventMetadata) EventType() EventType { return b.EventTypeVal }
func (b *EventMetadata) TxHash() common.Hash  { return b.TxHashVal }
func (b *EventMetadata) BlockNumber() uint64  { return b.BlockNumberVal }
func (b *EventMetadata) SetSwapData(data SwapData) {
	b.SwapDataVal = data
	b.SwapParsed = true
}
func (b *EventMetadata) SwapDataIsParsed() bool    { return b.SwapParsed }
func (b *EventMetadata) OuterIndex() int64         { return b.OuterIndexVal }
func (b *EventMetadata) TransactionIndex() *uint64 { return b.TransactionIndexVal }
func (b *EventMetadata) Clone() UnifiedEvent       { clone := *b; return &clone }

func (b *EventMetadata) NewEventMetadata() {

}
