package comm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// 枚举类型
type EventSig common.Hash

// 枚举值
var (
	UniswapV2Swap EventSig = EventSig(crypto.Keccak256Hash([]byte("Swap(address,uint256,uint256,uint256,uint256,address)")))
)

type ProtocolType string

const (
	DEX     ProtocolType = "DEX"
	Token   ProtocolType = "Token"
	Lending ProtocolType = "LENDING"
	Unknow  ProtocolType = "Unknow"
)

type ProtocolImpl string

const (
	UniswapV2 ProtocolImpl = "UniswapV2"
	SushiSwap ProtocolImpl = "SushiSwap"
	ERC20Std  ProtocolImpl = "ERC20Standard"
	Aave      ProtocolImpl = "Aave"
)

type EventMethod string

var (
	Swap      EventMethod = "Swap"
	FlashLoan EventMethod = "FlashLoan"
)

type SwapData struct {
	FromToken   common.Address `json:"from_token"`
	ToToken     common.Address `json:"to_token"`
	FromAmount  *big.Int       `json:"from_amount"`
	ToAmount    *big.Int       `json:"to_amount"`
	Sender      common.Address `json:"sender"`
	Recipient   common.Address `json:"recipient"`
	Description *string        `json:"description,omitempty"` // 可选字段
}

type TransactionDirection string

const (
	DirectionIn   TransactionDirection = "IN"   // 入金
	DirectionOut  TransactionDirection = "OUT"  // 出金
	DirectionSwap TransactionDirection = "SWAP" // 兑换
)

// 行为类型
type TransactionType string

const (
	TypeTransfer TransactionType = "Transfer"
	TypeSwap     TransactionType = "Swap"
	TypeDeposit  TransactionType = "Deposit"
	TypeWithdraw TransactionType = "Withdraw"
	TypeContract TransactionType = "ContractCall"
)
