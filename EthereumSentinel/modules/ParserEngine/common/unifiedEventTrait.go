package common

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type BaseEvent struct {
	EventTypeVal EventMethod    // 事件类型，例如 Swap、Deposit、Withdraw 等
	From         common.Address //事件发起地址
	TokenName    []string       // 事件涉及的代币名称（如果有的话）
	AmountVal    []*big.Int     // 事件涉及的以太坊代币数量（单位: wei）
	// PriceVal     []decimal.Decimal // 事件涉及的一个代币价格（如果有的话，单位: USD）
	RealValueVal []decimal.Decimal // 事件涉及的代币的实际价值（如果有的话，单位: USD）
}

// UnifiedEventData 封装以太坊链上事件的基础信息，用于统一事件处理
type EventMetadata struct {
	TxHashVal    common.Hash  // 交易哈希（唯一标识交易）
	ProtocolType ProtocolType //协议类型
	ProtocolImpl ProtocolImpl // 协议具体，例如 UniswapV2、Sushiswap 等
	AgeVal       time.Time    // 事件发生时间
	// FromVal      common.Address
	ToVal common.Address

	BlockNumberVal      *big.Int // 区块高度，用于排序和定位交易所在区块
	OuterIndexVal       uint     // log 在交易中的索引
	TransactionIndexVal uint     // 交易在区块中的索引位置，用于区分同区块内多笔交易
}

// UnifiedEvent 接口
type UnifiedEvent interface {
	TxHash() common.Hash        //以太坊的 tx hash，用于唯一标识触发事件的交易
	EventType() EventSig        //事件类型，比如 Swap、Deposit、Withdraw 等
	ProtocolType() ProtocolType //
	ProtocolImpl() ProtocolImpl //协议类型，比如 UniswapV2、Sushiswap 等
	Age() time.Time             //事件发生的时间

	To() common.Address //事件接收地址

	BlockNumber() *big.Int  //区块高度
	OuterIndex() uint       //log 在交易中的索引
	TransactionIndex() uint //Ethereum 交易在区块中的索引

	CoreEvent() BaseEvent // 返回基础事件信息
	Detail() any          // 双/多资产信息，具体结构根据 EventType

	// Clone() UnifiedEvent
	//HandleMs() int64      // 事件处理耗时（毫秒）
	//SetHandleMs(ms int64) // 设置处理耗时
}

type UnifiedEventData struct {
	*EventMetadata
	*BaseEvent
	DetailVal any // 双/多资产信息，具体结构根据 EventType
}

func (b *UnifiedEventData) EventType() EventSig        { return b.BaseEvent.EventTypeVal }
func (b *UnifiedEventData) TxHash() common.Hash        { return b.EventMetadata.TxHashVal }
func (b *UnifiedEventData) ProtocolType() ProtocolType { return b.EventMetadata.ProtocolType }
func (b *UnifiedEventData) ProtocolImpl() ProtocolImpl { return b.EventMetadata.ProtocolImpl }
func (b *UnifiedEventData) Age() time.Time             { return b.EventMetadata.AgeVal }

// func (b *UnifiedEventData) From() common.Address       { return b.EventMetadata.ToVal }
func (b *UnifiedEventData) To() common.Address     { return b.EventMetadata.ToVal }
func (b *UnifiedEventData) BlockNumber() *big.Int  { return b.EventMetadata.BlockNumberVal }
func (b *UnifiedEventData) OuterIndex() uint       { return b.EventMetadata.OuterIndexVal }
func (b *UnifiedEventData) TransactionIndex() uint { return b.EventMetadata.TransactionIndexVal }

func (b *UnifiedEventData) Detail() any { return b.DetailVal }
func (b *UnifiedEventData) CoreEvent() BaseEvent {
	return *b.BaseEvent
}

// func (b *UnifiedEventData) Clone() UnifiedEvent    { clone := *b; return &clone }

type TokenEvent struct {
	Protocol     string   // ERC20 or ERC721
	EventType    string   // Transfer / Approval / ApprovalForAll
	TokenAddress string   // 合约地址
	Operator     string   // 操作者地址 (spender/operator)
	From         string   // 转出地址
	To           string   // 转入地址
	TokenId      *big.Int // NFT ID (ERC721 专用，ERC20 为 nil)
	Amount       *big.Int // 转账数量 (ERC20 专用，ERC721 固定为 1)
	Approved     *bool    // 是否授权 (ERC721 的 setApprovalForAll 专用)
	TxHash       string   // 交易哈希
	BlockNumber  uint64   // 区块高度
}

type MethodName string

const (
	Transfer     MethodName = "transfer"
	TransferFrom MethodName = "transferFrom"
	Approve      MethodName = "approve"
)

type Protocol string

const (
	ERC20 Protocol = "ERC20"
)
