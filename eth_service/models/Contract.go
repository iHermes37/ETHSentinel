package models

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

//---------------------新部署合约的监控，和与巨鲸交互的合约交互-------------------

type ERCStandard struct {
	Ercname  ContractType
	Opmethod string
	Params   map[string]interface{}
}

// 合约类型
type ContractType string

const (
	EOA           ContractType = "EOA"    // 普通地址
	KnownContract ContractType = "ERC20"  // 已知合约
	NewContract   ContractType = "ERC721" // 新部署合约
)

type ConstractInfo struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement"`

	Address      common.Address `json:"address" gorm:"type:varchar(42);uniqueIndex"` // 合约地址
	ContractType ContractType   `json:"contractType" gorm:"type:varchar(20);index"`  // 合约类型
	ContractAge  time.Duration  `json:"contractAge,omitempty" gorm:"type:bigint"`    // 合约年龄，可选（天数或区块高度）
	//IsNewProject bool         `json:"isNewProject" gorm:"type:boolean"`           // 是否为新项目

	TxHash      common.Hash `json:"txHash" gorm:"type:varchar(66);uniqueIndex"` // 关联交易哈希，唯一
	BlockNumber *big.Int    `json:"blockNumber,omitempty" gorm:"type:bigint"`   // 部署区块号
	DeployTime  time.Time   `json:"deployTime,omitempty" gorm:"type:bigint"`    // 部署合约时间戳（秒）

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
