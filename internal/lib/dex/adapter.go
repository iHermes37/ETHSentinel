package dex

import (
	"fmt"

	uniswap_v2 "github.com/ETHSentinel/internal/lib/dex/uniswap_v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// DEXProtocol 协议类型（从原 types 包内联，避免循环依赖）
type DEXProtocol string

const (
	UniswapV2Protocol DEXProtocol = "Uniswap"
)

type DexAdapter struct {
	client *ethclient.Client
}

func (da *DexAdapter) SelectDexPair(d *DEXProtocol, pairAddr common.Address) (*uniswap_v2.Uniswappair, error) {
	if *d == UniswapV2Protocol {
		return uniswap_v2.NewUniswappair(pairAddr, da.client)
	}
	return nil, fmt.Errorf("unsupported dex: %s", *d)
}

func (da *DexAdapter) GetFactoryAddress(d *DEXProtocol) common.Address {
	switch *d {
	case UniswapV2Protocol:
		return common.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f") // ETH 主网
	default:
		panic(fmt.Sprintf("unsupported DEX: %s", *d))
	}
}
