package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"time"
)

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

// 核心公共字段：交易时间、地址、交易哈希、交易类型、目标地址
type BaseWhaleTransaction struct {
	Time    time.Time       `json:"time" gorm:"type:datetime;index"`
	Address common.Address  `json:"address" gorm:"type:varchar(42);index"` // 巨鲸或用户地址
	Type    TransactionType `json:"type" gorm:"type:varchar(20)"`
	TxHash  common.Hash     `json:"txHash" gorm:"type:varchar(66);uniqueIndex"`
	To      *common.Address `json:"to,omitempty" gorm:"type:varchar(42)"`
}

// 与 DeFi 协议交互
type DeFiTxDetail struct {
	TokensInsymbol  []string             `json:"tokensIn"`            // 输入代币，可能多个
	TokensOutsymbol []string             `json:"tokensOut,omitempty"` // 输出代币
	AmountsIn       []float64            `json:"amountsIn"`
	AmountsOut      []float64            `json:"amountsOut,omitempty"`
	Exchange        string               `json:"exchange,omitempty"` // DEX 或借贷协议
	Direction       TransactionDirection `json:"direction"`
}

// 普通 ERC20 交互
type ERC20TxDetail struct {
	Token     string               `json:"token"`
	Amount    float64              `json:"amount"`
	Direction TransactionDirection `json:"direction"` // IN / OUT
}

// 与普通用户的转账
type UserTxDetail struct {
	Asset     string               `json:"asset"` // ETH 或 Token
	Amount    float64              `json:"amount"`
	Direction TransactionDirection `json:"direction"`
}

// 巨鲸交易明细
type WhaleTransaction struct {
	ID          int64                `json:"id" gorm:"primaryKey;autoIncrement"`
	Base        BaseWhaleTransaction `json:"base"`
	DeFiDetail  *DeFiTxDetail        `json:"defi,omitempty"`
	ERC20Detail *ERC20TxDetail       `json:"erc20,omitempty"`
	UserDetail  *UserTxDetail        `json:"user,omitempty"`
	CreatedAt   time.Time            `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time            `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (w *WhaleTransaction) BuildWhaleTransaction(tx *types.Transaction, from common.Address, Dt *DeFiTxDetail) error {
	w.Base = BaseWhaleTransaction{tx.Time(), from, "", tx.Hash(), tx.To()}
	w.DeFiDetail = Dt
	return nil
}
