package models

import (
	"github.com/ethereum/go-ethereum/common"
)

// ==================合约类型==================
type ContractType int

const (
	DEX ContractType = iota // dex类型合约
	Lend
)

type ConstractProto int

const (
	Uniswap_V2 ConstractProto = iota
)

type ConstractMethod int

const (
	Burn ConstractMethod = iota
	Mint
	Swap
)

type ContractTable string

const (
	Defi_Table            ContractTable = "Know"
	MonitorContract_Table ContractTable = "UnKnow"
)

// ============已知合约====================
type ConstractInfo struct {
	FactoryAddr common.Address
	Category    ContractType
	Protocol    ConstractProto
	RouterAddr  common.Address
	Method      ConstractMethod
}

// ===========未知合约=========================
type UnKnowConstractInfo struct {
	FactoryAddr string `json:"address"`      // 存储为字符串
	ContractAge int64  `json:"contract_age"` // 存储为秒数
	TxHash      string `json:"tx_hash"`      // 存储为字符串
	BlockNumber string `json:"block_number"` // 存储为字符串
}
