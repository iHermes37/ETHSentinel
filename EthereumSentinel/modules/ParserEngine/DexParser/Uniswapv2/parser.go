package uniswapv2

import (
	"fmt"
	"math/big"

	abligens "github.com/Crypto-ChainSentinel/modules/ParserEngine/DexParser/abigens"
	dexcommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
	connmanager "github.com/Crypto-ChainSentinel/modules/connmanager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

// USDC/WETH: 0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc
// DAI/WETH: 0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11
// USDT/WETH: 0x0d4a11d5EEaaC28EC3F61d100daF4d40471f1852
var ContractAddress = map[common.Address]dexcommon.ProtocolImpl{
	common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"): dexcommon.UniswapV2,
	common.HexToAddress("0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11"): dexcommon.UniswapV2,
	common.HexToAddress("0x0d4a11d5EEaaC28EC3F61d100daF4d40471f1852"): dexcommon.UniswapV2,
}

var pairAddrs = func() []common.Address {
	var pairs []common.Address
	for pairAddr, _ := range ContractAddress {
		pairs = append(pairs, pairAddr)
	}
	return pairs
}

type contractMap map[common.Address]*abligens.UniswappairFilterer

var UniswapFilterers = func(client *ethclient.Client, pairaddrs []common.Address) contractMap {
	contractfiltermap := make(contractMap)
	for _, pairaddr := range pairaddrs {
		filter, _ := abligens.NewUniswappairFilterer(pairaddr, client)
		contractfiltermap[pairaddr] = filter
	}
	return contractfiltermap
}

// ----------------------------------------------------------------

// 事件路由表（Invoker内部可用map管理命令）
var UniswapV2EventsConfig = map[dexcommon.EventSig]dexcommon.EventParserFunc{
	dexcommon.UniswapV2Swap: ParseSwapEvent,
}

func ParseSwapEvent(log types.Log, metadata dexcommon.EventMetadata) (dexcommon.UnifiedEvent, error) {
	client := connmanager.InfuraConn()
	contractMap := UniswapFilterers(client, pairAddrs())
	//if filterer == nil {
	//	return nil, fmt.Errorf("filterer 未初始化")
	//}
	// 类型断言为具体类型
	u2, ok := contractMap[log.Address]
	if !ok {
		return nil, fmt.Errorf("filterer 类型错误")
	}

	pair, _ := abligens.NewUniswappair(log.Address, client)

	token0Addr, _ := pair.Token0(nil) // token0 地址
	token1Addr, _ := pair.Token1(nil) // token1 地址

	fmt.Printf("hhhhhhhhhh,success")
	swapEvent, err := u2.ParseSwap(log)
	if err != nil {
		return nil, err
	}

	return &dexcommon.UnifiedEventData{
		EventMetadata: &metadata,
		BaseEvent: &dexcommon.BaseEvent{
			EventTypeVal: dexcommon.Swap,
			From:         swapEvent.Sender,
			TokenName:    []string{GetTokenName(token0Addr), GetTokenName(token1Addr)},
			AmountVal:    []*big.Int{swapEvent.Amount0In, swapEvent.Amount1Out},
			RealValueVal: 
		},
		DetailVal: swapEvent,
	}, nil

}

func GetTokenName(addr common.Address) string {

}

func GetTokenRealVal(addr common.Address) decimal.Decimal {

}
