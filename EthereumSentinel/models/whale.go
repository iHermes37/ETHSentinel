package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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

type TokenAmount struct {
	TokenName string
	Balance   common.Decimal
}

type Whale struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"` // 数据库主键
	Address   common.Address `json:"address" gorm:"type:varchar(42);uniqueIndex"`
	Positions []TokenAmount  `json:"amount" gorm:"type:varchar(42)"`
	// FirstSeen time.Time      `json:"firstSeen" gorm:"type:datetime;index"`    // 第一次发现时间
	Note string `json:"note,omitempty" gorm:"type:varchar(255)"` // 可选：备注，记录发现来源
	// CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	// UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
