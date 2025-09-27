package dexcommon

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// 枚举类型
type EventSig common.Hash

// 枚举值
var (
	UniswapV2_Swap EventSig = EventSig(crypto.Keccak256Hash([]byte("Swap(address,uint256,uint256,uint256,uint256,address)")))
)

type Protocol string

const (
	UniswapV2 Protocol = "UniswapV2"
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
