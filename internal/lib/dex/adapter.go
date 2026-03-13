package dex

import (
	"fmt"
	"github.com/Crypto-ChainSentinel/internal/lib/dex/uniswap_v2"
	"github.com/Crypto-ChainSentinel/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type DexAdapter struct {
	client *ethclient.Client
}

func (da *DexAdapter) SelectDexPair(d *types.DEXProtocol, pairAddr common.Address) (*uniswap_v2.Uniswappair, error) {
	if *d == types.Uniswap_V2_Protool {
		return uniswap_v2.NewUniswappair(pairAddr, da.client)
	}
	return nil, fmt.Errorf("unsupported dex")
}

func (da *DexAdapter) GetFactoryAddress(d *types.DEXProtocol) common.Address {
	switch *d {
	case types.Uniswap_V2_Protool:
		return common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f") // ETH 主网
	//case types.Pancake_V2_Protocol:
	//	return common.HexToAddress("0xBCfCcbde45cE874adCB698cC183deBcF17952812")
	// ...
	default:
		panic("unsupported DEX")
	}
}
