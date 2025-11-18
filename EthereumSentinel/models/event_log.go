package models

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type Topics struct {
	hash common.Hash
	from common.Address
	to   common.Address
}

type EventLog struct {
	ProtocolImpl
	EventName `json:"event_name,omitempty"`
	topic     Topics
	Data
}

// 交易路径概览
// 用户 ETH → WETH → 1inch路由 → SushiSwap池 → LUNC → 1inch路由 → 用户钱包
type TradePathOverview struct {
}

type LogTable struct {
	ID         int64       `json:"id" gorm:"primaryKey;autoIncrement"`
	TxHash     common.Hash `json:"txHash" gorm:"type:varchar(66);uniqueIndex"`
	TxProtocol TxProtocol
	Overview   TradePathOverview
	EventLog   *EventLog

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
