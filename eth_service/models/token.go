package models

import (
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
