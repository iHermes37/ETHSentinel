package UniswapV2Protocol

import (
	abligens "github.com/CryptoQuantX/chain_monitor/modules/parserEngine/dex_parser/abigens"
	dexcommon "github.com/CryptoQuantX/chain_monitor/modules/parserEngine/dex_parser/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var Configs = []dexcommon.EventParseConfig{
	{
		ContractAddress: UniswapV2_Address[0], // UniswapV2Protocol pair
		Protocol:        dexcommon.UniswapV2,
		EventType:       dexcommon.UniswapV2_SwapBuy,
		Parser:          ParseSwapEvent,
	},
}

var UniswapV2_Address = []common.Address{
	common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"),
	common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"),
}

func ParseSwapEvent(log types.Log, metadata dexcommon.EventMetadata, filterer *abligens.UniswappairFilterer) (dexcommon.UnifiedEvent, error) {
	swapEvent, err := filterer.ParseSwap(log)
	if err != nil {
		return nil, err
	}

	return &UniswapV2_SwapEvent{
		EventMetadata: metadata,
		Event:         *swapEvent, // 注意：Event 是值类型
	}, nil
}
