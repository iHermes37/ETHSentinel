package UniswapV2Protocol

import (
	"fmt"
	abligens "github.com/Crypto-ChainSentinel/modules/parserEngine/dex_parser/abigens"
	dexcommon "github.com/Crypto-ChainSentinel/modules/parserEngine/dex_parser/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var UniswapV2EventsConfig = map[dexcommon.EventSig]dexcommon.EventParserFunc{
	dexcommon.UniswapV2_Swap: ParseSwapEvent,
}

// USDC/WETH: 0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc
// DAI/WETH: 0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11
// USDT/WETH: 0x0d4a11d5EEaaC28EC3F61d100daF4d40471f1852
var ContractAddress = map[common.Address]dexcommon.Protocol{
	common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"): dexcommon.UniswapV2,
	common.HexToAddress("0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11"): dexcommon.UniswapV2,
	common.HexToAddress("0x0d4a11d5EEaaC28EC3F61d100daF4d40471f1852"): dexcommon.UniswapV2,
}

func ParseSwapEvent(log types.Log, metadata dexcommon.EventMetadata, filterer *abligens.UniswappairFilterer) (dexcommon.UnifiedEvent, error) {
	if filterer == nil {
		return nil, fmt.Errorf("filterer 未初始化")
	}

	fmt.Printf("hhhhhhhhhh,success")
	swapEvent, err := filterer.ParseSwap(log)
	if err != nil {
		return nil, err
	}

	return &UniswapV2_SwapEvent{
		EventMetadata: metadata,
		Event:         *swapEvent, // 注意：Event 是值类型
	}, nil
}
