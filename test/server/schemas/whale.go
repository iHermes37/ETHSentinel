package schemas

import (
	"math/big"
	"time"

	ParserEngine "github.com/Crypto-ChainSentinel/internal/parser/comm"
)

// 可用于动态查询巨鲸信息
type WhaleQueryParams struct {
	Address   *string    `json:"address,omitempty"`   // 巨鲸地址，可选，支持精确或模糊查询
	FirstSeen *time.Time `json:"firstSeen,omitempty"` // 精确查询首次发现时间，可选
	Note      *string    `json:"note,omitempty"`      // 备注内容，可选，支持模糊匹配
	CreatedAt *time.Time `json:"createdAt"`           // 创建时间
	UpdatedAt *time.Time `json:"updatedAt"`           // 更新时间
}

// 响应给前端的巨鲸信息结构体
type WhaleResponse struct {
	Address   string    `json:"address"`        // 巨鲸地址
	FirstSeen time.Time `json:"firstSeen"`      // 第一次发现时间
	Note      *string   `json:"note,omitempty"` // 可选备注信息
	CreatedAt time.Time `json:"createdAt"`      // 创建时间
	UpdatedAt time.Time `json:"updatedAt"`      // 更新时间
}

// 动态查询巨鲸交易参数
type WhaleTradeParams struct {
	// ===== BaseWhaleTransaction =====
	Address   *string    `json:"address,omitempty" form:"address"`     // 巨鲸地址
	Type      *string    `json:"type,omitempty" form:"type"`           // 交易类型
	TxHash    *string    `json:"txHash,omitempty" form:"txHash"`       // 交易哈希
	To        *string    `json:"to,omitempty" form:"to"`               // 目标地址
	StartTime *time.Time `json:"startTime,omitempty" form:"startTime"` // 起始时间
	EndTime   *time.Time `json:"endTime,omitempty" form:"endTime"`     // 结束时间

	// ===== DeFiTxDetail =====
	Exchange     *string  `json:"exchange,omitempty" form:"exchange"`   // DEX 或借贷协议
	Direction    *string  `json:"direction,omitempty" form:"direction"` // IN / OUT
	TokenIn      *string  `json:"tokenIn,omitempty" form:"tokenIn"`     // 输入代币
	TokenOut     *string  `json:"tokenOut,omitempty" form:"tokenOut"`   // 输出代币
	AmountInMin  *float64 `json:"amountInMin,omitempty" form:"amountInMin"`
	AmountInMax  *float64 `json:"amountInMax,omitempty" form:"amountInMax"`
	AmountOutMin *float64 `json:"amountOutMin,omitempty" form:"amountOutMin"`
	AmountOutMax *float64 `json:"amountOutMax,omitempty" form:"amountOutMax"`

	// ===== ERC20TxDetail =====
	ERC20Token     *string  `json:"erc20Token,omitempty" form:"erc20Token"`
	ERC20AmountMin *float64 `json:"erc20AmountMin,omitempty" form:"erc20AmountMin"`
	ERC20AmountMax *float64 `json:"erc20AmountMax,omitempty" form:"erc20AmountMax"`

	// ===== UserTxDetail =====
	UserAsset     *string  `json:"userAsset,omitempty" form:"userAsset"`
	UserAmountMin *float64 `json:"userAmountMin,omitempty" form:"userAmountMin"`
	UserAmountMax *float64 `json:"userAmountMax,omitempty" form:"userAmountMax"`
}

// 前端统一响应结构
type WhaleTradeResponse struct {
	// 核心公共字段
	Time    time.Time `json:"time"`    // 交易时间
	Address string    `json:"address"` // 发起地址（巨鲸或用户）
	Type    string    `json:"type"`    // 交易类型：DeFi / ERC20 / Transfer
	TxHash  string    `json:"txHash"`  // 交易哈希
	To      string    `json:"to,omitempty"`

	// 通用资产信息（方便前端展示）
	AssetsIn   []string  `json:"assetsIn,omitempty"`   // 输入资产的符号列表
	AmountsIn  []float64 `json:"amountsIn,omitempty"`  // 对应输入资产的数量
	AssetsOut  []string  `json:"assetsOut,omitempty"`  // 输出资产的符号列表
	AmountsOut []float64 `json:"amountsOut,omitempty"` // 对应输出资产的数量

	// 协议/代币/方向（不同类型细节）
	Exchange  string `json:"exchange,omitempty"`  // DeFi 协议（如 Uniswap）
	Token     string `json:"token,omitempty"`     // ERC20 代币
	Asset     string `json:"asset,omitempty"`     // 普通转账资产
	Amount    string `json:"amount,omitempty"`    // ERC20 或转账金额（string 避免精度丢失）
	Direction string `json:"direction,omitempty"` // IN / OUT

	// 元信息
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WhaleDetectionMethod int

const (
	HoldingsAnalysis WhaleDetectionMethod = iota
	ChainScan
	TransactionPattern
)

type CapturedWhale struct {
	Method    WhaleDetectionMethod
	Balance   *int
	Threshold *int
}

// ===================================
type WhaleAssetsResponse struct {
	ID        int64  `json:"id"`
	Address   string `json:"address"` // 使用字符串
	TokenName string `json:"token_name"`
	Balance   string `json:"balance"`
}

type WhaleTrackSettingRequest struct {
	RealMonitor bool
	StartBlock  *big.Int
	EndBlock    *big.Int
	Selected    *map[ParserEngine.ProtocolTypeName][]ParserEngine.ProtocolImplName
}
