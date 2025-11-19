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

type DeFi struct {
	Dex Dex `json:"Dex"`
}

//-------------------uniswap--------------------------

// 定义最内层的 V2/V3 结构
type UniswapV2 struct {
	Router  string `json:"Router"`
	Factory string `json:"Factory"`
	Pair    string `json:"Pair"`
}

type UniswapV3 struct {
	Router  string `json:"Router"`
	Factory string `json:"Factory"`
	Pool    string `json:"Pool"`
}

// 定义 Uniswap 结构，里面嵌套 V2 和 V3
type Uniswap struct {
	V2 UniswapV2 `json:"V2"`
	V3 UniswapV3 `json:"V3"`
}

//-----------------SushiSwap---------------------------------------

type SushiSwap struct {
	Router  string `json:"Router"`
	Factory string `json:"Factory"`
	Pair    string `json:"Pair"`
}

type Uniswap_V3 struct {
	GraphQLDEX
}

type Sushiswap struct {
	GraphQLDEX
}

type DexPairLiquidity struct {
	PairAddress common.Address `json:"pairAddress"` // 交易对合约地址
	Token0      Token          `json:"token0"`
	Token1      Token          `json:"token1"`
	Reserve0    *big.Int       `json:"reserve0"`    // token0 储备量
	Reserve1    *big.Int       `json:"reserve1"`    // token1 储备量
	TotalSupply *big.Int       `json:"totalSupply"` // LP 总量
	Timestamp   time.Time      `json:"timestamp"`   // 事件时间
}

type DexPairLiquidityChange struct {
	PairAddress common.Address `json:"pairAddress"`
	Token0      Token          `json:"token0"`
	Token1      Token          `json:"token1"`

	DeltaReserve0 *big.Int `json:"deltaReserve0"` // token0 变化量
	DeltaReserve1 *big.Int `json:"deltaReserve1"` // token1 变化量
	DeltaLP       *big.Int `json:"deltaLP"`       // LP 总量变化
	Reserve0After *big.Int `json:"reserve0After"` // 变化后的储备量
	Reserve1After *big.Int `json:"reserve1After"`
	TotalLPAfter  *big.Int `json:"totalLPAfter"`

	Timestamp time.Time `json:"timestamp"` // 变化发生时间
	EventType string    `json:"eventType"` // "Mint" / "Burn" / "Swap" / "ManualSnapshot"
}

// ----------------------------------------

type Dex struct {
	//Uni   Uniswap   `json:"Uniswap"`
	//Sushi SushiSwap `json:"SushiSwap"`
}

type DEXProtocol string

const (
	Uniswap_V2_Protool DEXProtocol = "Uniswap"
)
