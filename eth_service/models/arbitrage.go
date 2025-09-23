package models

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

type Token struct {
	Id     common.Address `json:"id"`
	Symbol string         `json:"symbol"`
}

type Pair struct {
	Token0      Token            `json:"token0"` // 代币0的合约地址或符号
	Token1      Token            `json:"token1"` // 代币1的合约地址或符号
	PoolAddr    common.Address   `json:"id"`     // Pool 合约地址
	Fee         *uint32          `json:"feeTier,omitempty"`
	PairReserve *DexPairReserver `json:"pairReserve,omitempty"`
}

type Pairs struct {
	Pair []Pair `json:"pools,pairs"`
}

type DexPairReserver struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
}

type ArbitrageDirection int

const (
	None ArbitrageDirection = iota
	AtoB
	BtoA
)

func (d ArbitrageDirection) String() string {
	return [...]string{"None", "A->B", "B->A"}[d]
}

// 套利机会的结果数据结构，目的是把整个跨DEX套利过程的数据都明确记录下来。
type ArbitrageOpportunity struct {
	X big.Int // 借出的基准代币数量（闪电贷借的数量）
	Y float64 // 第一个 DEX 交易后得到的另一种代币数量
	Z float64 // 第二个 DEX 交易后兑换回基准代币数量
	//Profit          float64   // 净利润 = Z - X（扣除手续费、滑点后的收益）
	EstimatedProfit float64   `json:"estimatedProfit"`  // 预期利润（单位可为ETH或USDT）
	ProfitRatio     float64   `json:"profitRatio"`      // 预期收益率
	GasCost         float64   `json:"gasCost"`          // 预估手续费
	NetProfit       float64   `json:"netProfit"`        // 扣手续费后的净收益
	Slippage        float64   `json:"slippage"`         // 预期滑点
	RiskLevel       string    `json:"riskLevel"`        // 风险等级：低/中/高
	TxHash          *string   `json:"txHash,omitempty"` // 若已执行，返回交易哈希
	Timestamp       time.Time `json:"timestamp"`        // 发现时间
	Status          string    `json:"status"`           // 状态：PENDING/EXECUTED/FAILED
}

// 交易对数据结构
type CrossPairData struct {
	// DEX A 的信息
	Pair_DexA Pair
	// DEX B 的信息
	Pair_DexB   Pair
	Direction   ArbitrageDirection   // 套利方向 "AtoB" 或 "BtoA"
	PullToken   Token                // 决定买入的基准代币
	Opportunity ArbitrageOpportunity // 当前套利机会
}
