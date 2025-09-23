package model

import (
	"github.com/CryptoQuantX/chain_monitor/models"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

// 套利策略类型
type ArbitrageType string

const (
	CrossDex ArbitrageType = "CROSS_DEX" // 跨DEX套利
	Triangle ArbitrageType = "TRIANGLE"  // 三角套利
)

type FlashLoanDetail struct {
	lendingAgreement string
	loanAmount       float64
	LoanToken        string
}

// 前端请求参数
type ArbitrageOpportunityParams struct {
	Token0      common.Address `json:"token0"`      // 代币0 地址
	Token1      common.Address `json:"token1"`      // 代币1 地址
	DexA        string         `json:"dexA"`        // 第一个DEX名
	DexB        string         `json:"dexB"`        // 第二个DEX名
	Direction   string         `json:"direction"`   // "AtoB" 或 "BtoA" 或 "Any"
	BorrowToken common.Address `json:"borrowToken"` // 基准借入代币
	Amount      float64        `json:"amount"`      // 借入数量
	MinProfit   float64        `json:"minProfit"`   // 最小利润阈值
}

// 前端响应参数
type ArbitrageOpportunityResponse struct {
	PairDexA  models.Pair `json:"pairDexA"`  // DEX A 的交易对信息
	PairDexB  models.Pair `json:"pairDexB"`  // DEX B 的交易对信息
	Direction string      `json:"direction"` // 套利方向 "A->B" or "B->A"

	BorrowToken models.Token `json:"borrowToken"` // 借入代币
	X           float64      `json:"x"`           // 借入数量
	Y           float64      `json:"y"`           // DEX A 兑换后数量
	Z           float64      `json:"z"`           // DEX B 再兑换回来数量
	//Profit      float64      `json:"profit"`      // 最终净利润
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

// ------------------------------------------------------------------------------------------
type DexLiquidityQueryParams struct {
	PairAddress *common.Address `json:"pairAddress,omitempty"` // 可选，单个交易对
	Token0      *common.Address `json:"token0,omitempty"`      // 可选
	Token1      *common.Address `json:"token1,omitempty"`      // 可选
	StartTime   *time.Time      `json:"startTime,omitempty"`   // 查询起始时间
	EndTime     *time.Time      `json:"endTime,omitempty"`     // 查询结束时间
	IntervalSec int             `json:"intervalSec,omitempty"` // 分段时间间隔（秒），用于汇总快照
	Limit       int             `json:"limit,omitempty"`       // 最大条数
}

type DexLiquiditySnapshot struct {
	Timestamp time.Time `json:"timestamp"` // 时间点
	Reserve0  *big.Int  `json:"reserve0"`  // token0 储备量
	Reserve1  *big.Int  `json:"reserve1"`  // token1 储备量
	TotalLP   *big.Int  `json:"totalLP"`   // LP 总量
}

type DexLiquidityResponse struct {
	PairAddress common.Address         `json:"pairAddress"` // 交易对地址
	Token0      models.Token           `json:"token0"`
	Token1      models.Token           `json:"token1"`
	Snapshots   []DexLiquiditySnapshot `json:"snapshots"` // 按时间点的流动性快照
}

type DexLiquidityChangeResponse struct {
	PairAddress common.Address                  `json:"pairAddress"`
	Token0      models.Token                    `json:"token0"`
	Token1      models.Token                    `json:"token1"`
	Changes     []models.DexPairLiquidityChange `json:"changes"` // 变化量事件序列
}
