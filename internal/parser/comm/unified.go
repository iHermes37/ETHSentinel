// Package comm — UnifiedEvent 接口 & 默认实现
// 重构要点：
//   - EventMethod 类型从原 EventSig(common.Hash) 修正为语义字符串
//   - UnifiedEventData 嵌入改为组合，避免方法命名冲突
//   - 增加 Clone() 和 String() 便于调试
package comm

import (
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/common"
)

// ─────────────────────────────────────────────
//  RefToken
// ─────────────────────────────────────────────

// RefToken 事件涉及的单个代币引用
type RefToken struct {
	Name   string   // 代币符号，如 USDC / WETH
	Amount *big.Int // 数量（wei）
}

// ─────────────────────────────────────────────
//  BaseEvent & EventMetadata
// ─────────────────────────────────────────────

// EventMetadata 链上定位元信息（不含业务数据）
type EventMetadata struct {
	TxHash           common.Hash    // 交易哈希
	ProtocolTypeVal  ProtocolType   // 协议大类
	ProtocolImplVal  ProtocolImpl   // 具体实现
	Age              time.Time      // 事件时间
	To               common.Address // 合约地址
	BlockNumber      *big.Int       // 区块高度
	OuterIndex       uint           // log 在 tx 中的下标
	TransactionIndex uint           // tx 在 block 中的下标
}

// BaseEvent 核心业务字段（协议无关）
type BaseEvent struct {
	EventType  EventMethod       // 事件语义，如 Swap / Transfer
	From       common.Address    // 发起地址
	RefTokens  []RefToken        // 涉及的代币列表
	RealValues []decimal.Decimal // 各代币的 USD 价值（可选）
}

// ─────────────────────────────────────────────
//  UnifiedEvent 接口
// ─────────────────────────────────────────────

// UnifiedEvent 所有解析结果的统一接口，对上层屏蔽协议细节。
type UnifiedEvent interface {
	// 链上定位
	GetTxHash() common.Hash
	GetBlockNumber() *big.Int
	GetTransactionIndex() uint
	GetOuterIndex() uint
	GetAge() time.Time
	GetTo() common.Address

	// 协议信息
	GetProtocolType() ProtocolType
	GetProtocolImpl() ProtocolImpl

	// 事件业务字段
	GetEventType() EventMethod
	GetBase() BaseEvent

	// 详情（类型断言后使用，例如 *SwapData / *TransferData）
	GetDetail() any

	// 调试
	String() string
}

// ─────────────────────────────────────────────
//  UnifiedEventData — 默认实现
// ─────────────────────────────────────────────

// UnifiedEventData 是 UnifiedEvent 的通用值对象实现。
// 使用组合而非嵌入指针，避免 nil 解引用。
type UnifiedEventData struct {
	Metadata  EventMetadata
	Base      BaseEvent
	DetailVal any // *SwapData / *TransferData / …
}

// ── 链上定位 ──────────────────────────────────

func (e *UnifiedEventData) GetTxHash() common.Hash        { return e.Metadata.TxHash }
func (e *UnifiedEventData) GetBlockNumber() *big.Int      { return e.Metadata.BlockNumber }
func (e *UnifiedEventData) GetTransactionIndex() uint     { return e.Metadata.TransactionIndex }
func (e *UnifiedEventData) GetOuterIndex() uint           { return e.Metadata.OuterIndex }
func (e *UnifiedEventData) GetAge() time.Time             { return e.Metadata.Age }
func (e *UnifiedEventData) GetTo() common.Address         { return e.Metadata.To }

// ── 协议信息 ──────────────────────────────────

func (e *UnifiedEventData) GetProtocolType() ProtocolType { return e.Metadata.ProtocolTypeVal }
func (e *UnifiedEventData) GetProtocolImpl() ProtocolImpl { return e.Metadata.ProtocolImplVal }

// ── 事件业务字段 ──────────────────────────────

func (e *UnifiedEventData) GetEventType() EventMethod { return e.Base.EventType }
func (e *UnifiedEventData) GetBase() BaseEvent        { return e.Base }
func (e *UnifiedEventData) GetDetail() any            { return e.DetailVal }

// ── 调试 ──────────────────────────────────────

func (e *UnifiedEventData) String() string {
	return fmt.Sprintf(
		"[%s/%s] %s tx=%s block=%s",
		e.Metadata.ProtocolTypeVal,
		e.Metadata.ProtocolImplVal,
		e.Base.EventType,
		e.Metadata.TxHash.Hex(),
		e.Metadata.BlockNumber,
	)
}

// Clone 浅拷贝，用于并发场景
func (e *UnifiedEventData) Clone() *UnifiedEventData {
	cp := *e
	return &cp
}
