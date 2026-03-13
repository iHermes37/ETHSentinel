package uniswapv2

import (
	"github.com/Crypto-ChainSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/common"
)

var ContractAddress = map[common.Address]comm.ProtocolImplName{
	common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"): comm.UniswapV2,
	common.HexToAddress("0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11"): comm.UniswapV2,
	common.HexToAddress("0x0d4a11d5EEaaC28EC3F61d100daF4d40471f1852"): comm.UniswapV2,
}

var UniswapV2EventsConfig = map[comm.EventSig]comm.EventParserFunc{
	comm.UniswapV2Swap: ParseSwapEvent,
}
