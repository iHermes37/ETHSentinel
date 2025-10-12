package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// 枚举类型
type EventSig common.Hash

// 枚举值
var (
	UniswapV2Swap EventSig = EventSig(crypto.Keccak256Hash([]byte("Swap(address,uint256,uint256,uint256,uint256,address)")))
)

type ProtocolType string

const (
	DEX     ProtocolType = "DEX"
	ERC     ProtocolType = "ERC"
	Lending ProtocolType = "LENDING"
	Unknow  ProtocolType = "Unknow"
)

type ProtocolImpl string

const (
	UniswapV2 ProtocolImpl = "UniswapV2"
	SushiSwap ProtocolImpl = "SushiSwap"
	ERC20Std  ProtocolImpl = "ERC20Standard"
	Aave      ProtocolImpl = "Aave"
)

type EventMethod string

var (
	Swap EventMethod = "Swap"
)

type SwapData struct {
	FromToken   common.Address `json:"from_token"`
	ToToken     common.Address `json:"to_token"`
	FromAmount  *big.Int       `json:"from_amount"`
	ToAmount    *big.Int       `json:"to_amount"`
	Sender      common.Address `json:"sender"`
	Recipient   common.Address `json:"recipient"`
	Description *string        `json:"description,omitempty"` // 可选字段
}
