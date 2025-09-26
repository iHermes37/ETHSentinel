package dexcommon

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// 枚举类型
type EventType int

// 枚举值
const (
	UniswapV2_SwapBuy EventType = iota
	UniswapV2_SwapSell
	UniswapV2_Mint
	UniswapV2_PairCreated
	UniswapV2_Burn

	// Common events
	BlockMeta
	Unknown
)

// String 方法
func (e EventType) String() string {
	names := []string{
		"UniswapV2_SwapBuy",     // 0
		"UniswapV2_SwapSell",    // 1
		"UniswapV2_Mint",        // 2
		"UniswapV2_PairCreated", // 3
		"UniswapV2_Burn",        // 4
	}
	if int(e) < 0 || int(e) >= len(names) {
		return "Unknown"
	}
	return names[e]
}

// JSON 序列化/反序列化
func (e EventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *EventType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	names := map[string]EventType{}
	for i := UniswapV2_SwapBuy; i <= Unknown; i++ {
		names[i.String()] = i
	}
	if val, ok := names[s]; ok {
		*e = val
		return nil
	}
	*e = Unknown
	return nil
}

type SwapData struct {
	FromToken   common.Address `json:"from_token"`
	ToToken     common.Address `json:"to_token"`
	FromAmount  *big.Int       `json:"from_amount"`
	ToAmount    *big.Int       `json:"to_amount"`
	Sender      common.Address `json:"sender"`
	Recipient   common.Address `json:"recipient"`
	Description *string        `json:"description,omitempty"` // 可选字段
}
