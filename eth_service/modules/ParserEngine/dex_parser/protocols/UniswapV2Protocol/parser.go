package UniswapV2Protocol

import (
	"fmt"
	abligens "github.com/Crypto-ChainSentinel/modules/ParserEngine/dex_parser/abigens"
	dexcommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/dex_parser/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

var pairAddr []common.Address

type contractMap map[common.Address]*abligens.UniswappairFilterer

var UniswapFilterers = func(client *ethclient.Client, pairaddrs []common.Address) contractMap {
	contractfiltermap := make(contractMap)
	for _, pairaddr := range pairaddrs {
		filter, _ := abligens.NewUniswappairFilterer(pairaddr, client)
		contractfiltermap[pairaddr] = filter
	}

	return contractfiltermap
}

func ParseSwapEvent(log types.Log, metadata dexcommon.EventMetadata) (dexcommon.UnifiedEvent, error) {
	contractMap := UniswapFilterers(client, pairAddrs)
	//if filterer == nil {
	//	return nil, fmt.Errorf("filterer 未初始化")
	//}
	// 类型断言为具体类型
	u2, ok := contractMap[log.Address]
	if !ok {
		return nil, fmt.Errorf("filterer 类型错误")
	}

	fmt.Printf("hhhhhhhhhh,success")
	swapEvent, err := u2.ParseSwap(log)
	if err != nil {
		return nil, err
	}

	return &UniswapV2_SwapEvent{
		EventMetadata: metadata,
		Event:         *swapEvent, // 注意：Event 是值类型
	}, nil
}
