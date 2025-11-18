package models

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type TxProtocol string

const (
	DeFiTx  TxProtocol = "DeFi"
	TokenTx TxProtocol = "Token"
	ETHTx   TxProtocol = "ETH"
)

//// 与 DeFi 协议交互
//type DeFiTxDetail struct {
//	TokensInsymbol  []string  `json:"tokensIn"`            // 输入代币，可能多个
//	TokensOutsymbol []string  `json:"tokensOut,omitempty"` // 输出代币
//	AmountsIn       []float64 `json:"amountsIn"`
//	AmountsOut      []float64 `json:"amountsOut,omitempty"`
//	Exchange        string    `json:"exchange,omitempty"` // DEX 或借贷协议
//	//Direction       TransactionDirection `json:"direction"`
//}
//
//// 普通 Token 交互
//type TokenTxDetail struct {
//	Token  string  `json:"token"`
//	Amount float64 `json:"amount"`
//	//Direction TransactionDirection `json:"direction"` // IN / OUT
//}

// ======================================

// 核心公共字段：交易时间、地址、交易哈希、交易类型、目标地址
type BaseWhaleTxTable struct {
	ID          int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	TxHash      common.Hash     `json:"txHash" gorm:"type:varchar(66);uniqueIndex"`
	Address     *common.Address `json:"address" gorm:"type:varchar(42);index"` // 巨鲸或用户地址
	BlockNumber *big.Int
	From        *common.Address
	To          *common.Address `json:"to,omitempty" gorm:"type:varchar(42)"`
	Amount      *big.Int
	Time        time.Time `json:"time" gorm:"type:datetime;index"`
	//Type TransactionType `json:"type" gorm:"type:varchar(20)"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
