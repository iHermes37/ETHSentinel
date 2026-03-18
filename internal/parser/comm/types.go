// Package comm 定义链上解析引擎的核心类型。
// 重构要点：
//   - 将 ProtocolTypeName / ProtocolImplName 归一化为带前缀的常量
//   - EventSig 改为 [32]byte 别名（等同 common.Hash），去掉不必要的包装
//   - 统一 ParserCfg 配置树
package comm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// ─────────────────────────────────────────────
//  协议层次类型
// ─────────────────────────────────────────────

// ProtocolType 协议大类（DEX / Token / Lending …）
type ProtocolType string

const (
	ProtocolTypeDEX     ProtocolType = "DEX"
	ProtocolTypeToken   ProtocolType = "Token"
	ProtocolTypeLending ProtocolType = "Lending"
	ProtocolTypeUnknown ProtocolType = "Unknown"
)

// ProtocolImpl 具体协议实现名称
type ProtocolImpl string

const (
	ProtocolImplUniswapV2 ProtocolImpl = "UniswapV2"
	ProtocolImplSushiSwap ProtocolImpl = "SushiSwap"
	ProtocolImplERC20     ProtocolImpl = "ERC20"
	ProtocolImplERC721    ProtocolImpl = "ERC721"
	ProtocolImplAave      ProtocolImpl = "Aave"
)

// ─────────────────────────────────────────────
//  事件类型
// ─────────────────────────────────────────────

// EventSig Solidity 事件签名的 keccak256 哈希（即 log.Topics[0]）
type EventSig = common.Hash

// 预定义事件签名常量
var (
	SigUniswapV2Swap = EventSig(
		crypto.Keccak256Hash([]byte("Swap(address,uint256,uint256,uint256,uint256,address)")),
	)
	SigERC20Transfer = EventSig(
		crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")),
	)
	SigERC20Approval = EventSig(
		crypto.Keccak256Hash([]byte("Approval(address,address,uint256)")),
	)
	SigERC721Transfer = EventSig(
		crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")),
	)
)

// EventMethod 事件语义名称
type EventMethod string

const (
	EventMethodSwap      EventMethod = "Swap"
	EventMethodTransfer  EventMethod = "Transfer"
	EventMethodDeposit   EventMethod = "Deposit"
	EventMethodWithdraw  EventMethod = "Withdraw"
	EventMethodFlashLoan EventMethod = "FlashLoan"
	EventMethodApproval  EventMethod = "Approval"
)

// ─────────────────────────────────────────────
//  解析配置树
// ─────────────────────────────────────────────

// SelectedEvents 某个实现需要解析的事件列表（空 = 全部）
type SelectedEvents []EventMethod

// ImplConfig 单个协议实现的解析配置
type ImplConfig map[ProtocolImpl]SelectedEvents

// ParserCfg 完整解析配置，按协议大类分组
// 示例:
//   ParserCfg{
//     ProtocolTypeDEX: {ProtocolImplUniswapV2: {EventMethodSwap}},
//     ProtocolTypeToken: {ProtocolImplERC20: nil},
//   }
type ParserCfg map[ProtocolType]ImplConfig

// ─────────────────────────────────────────────
//  业务数据结构（对应原 SwapData 等）
// ─────────────────────────────────────────────

// SwapData UniswapV2 Swap 核心字段
type SwapData struct {
	FromToken  common.Address `json:"from_token"`
	ToToken    common.Address `json:"to_token"`
	FromAmount *big.Int       `json:"from_amount"`
	ToAmount   *big.Int       `json:"to_amount"`
	Sender     common.Address `json:"sender"`
	Recipient  common.Address `json:"recipient"`
}

// TransferData ERC20/ETH 转账字段
type TransferData struct {
	Token  common.Address `json:"token"`
	From   common.Address `json:"from"`
	To     common.Address `json:"to"`
	Amount *big.Int       `json:"amount"`
}
