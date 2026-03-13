package types

import (
	"github.com/Crypto-ChainSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/datatypes"
	"math/big"
	"time"
)

// ===============================================
type KnowedDefiTable struct {
	ID           int64                 `gorm:"primaryKey;autoIncrement"`
	ProtocolName comm.ProtocolTypeName `gorm:"type:varchar(64);not null;index"`
	ProtocolImpl comm.ProtocolImplName `gorm:"type:varchar(64);not null;index"`
	EventName    comm.EventMethod      `gorm:"type:varchar(64);not null;index"`
	Address      common.Address        `gorm:"type:char(42);not null;index"`
}

//=================================================

type WhaleAssetsTable struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"` // 数据库主键
	Address   common.Address `json:"address" gorm:"type:varchar(42);uniqueIndex"`
	TokenName string
	Balance   string
}

//================================================

// 字段：交易时间、地址、交易哈希、交易类型、目标地址
type WhaleTxTable struct {
	ID          int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	TxHash      common.Hash     `json:"txHash" gorm:"type:varchar(66);uniqueIndex"`
	Address     *common.Address `json:"address" gorm:"type:varchar(42);index"` // 巨鲸或用户地址
	BlockNumber *big.Int        `json:"block_number" gorm:"type:varchar(42);index"`
	From        *common.Address `json:"from,omitempty" gorm:"type:varchar(42)"`
	To          *common.Address `json:"to,omitempty" gorm:"type:varchar(42)"`
	Amount      *big.Int
	Time        time.Time `json:"time" gorm:"type:datetime;index"`

	// 存储为JSON
	Overview datatypes.JSON `json:"overview" gorm:"type:json"`
	EventLog datatypes.JSON `json:"eventLog" gorm:"type:json"`
}

//===========================================
