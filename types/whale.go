他package models

import (
	"github.com/Crypto-ETHSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type Holder struct {
	Address    common.Address `json:"address"`    // 持仓人地址
	Balance    *big.Int       `json:"balance"`    // 持仓数量
	ShareRatio float64        `json:"shareRatio"` // 持仓占比（百分比或小数）
}

type TokenHolders struct {
	Symbol  string         `json:"symbol"`  // 代币符号
	Address common.Address `json:"address"` // 代币合约地址
	Holders []Holder       `json:"holders"` // 持仓人列表，按持仓数量排行
}

// ===========================
type Topics struct {
	hash common.Hash
	from common.Address
	to   common.Address
}

type EventLog struct {
	comm.ProtocolImplName
	comm.EventMethod `json:"event_name,omitempty"`
	topic            Topics
}

// 交易路径概览
// 用户 ETH → WETH → 1inch路由 → SushiSwap池 → LUNC → 1inch路由 → 用户钱包
type TradePathOverview struct {
}

//==========================================
type WhaleAddress struct{
	ID          int64          `gorm:"primaryKey"`
	Address     common.Address `gorm:"type:varchar(42);uniqueIndex"`
	Chain       string         `gorm:"index"`   // eth、bsc、arb...
	Label       string         // 交易所 / 机构 / 个人
	Tag         string         // 巨鲸类型（做市商、套利机器人等）
	RiskLevel   int            // 风险等级
	IsContract  bool           // 是否合约地址
	FirstSeenAt time.Time
	LastActive  time.Time
}


type WhaleAssetSnapshot struct {
	ID        int64
	Address   common.Address `gorm:"index"`
	Chain     string         `gorm:"index"`

	TokenAddr common.Address `gorm:"index"`
	TokenName string
	Symbol    string

	Balance   decimal.Decimal
	UsdValue  decimal.Decimal

	BlockNum  uint64
	UpdatedAt time.Time
}


type WhaleTxRecord struct {
	ID          int64
	TxHash      common.Hash `gorm:"uniqueIndex"`

	Address     common.Address `gorm:"index"`
	Chain       string

	TokenAddr   common.Address
	Amount      decimal.Decimal
	UsdValue    decimal.Decimal

	TxType      string // swap transfer stake mint burn
	Dex         string // uniswap pancake etc

	BlockNum    uint64
	Timestamp   time.Time
}


type WhaleBehaviorEvent struct {
	ID          int64
	Address     common.Address
	EventType   string   // 大额买入、大额卖出、LP撤出等
	TokenAddr   common.Address
	ValueUSD    decimal.Decimal
	Confidence  float64
	Timestamp   time.Time
}
