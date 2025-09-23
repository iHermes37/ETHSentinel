package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	// 1. 连接 Infura WebSocket 节点
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/YOUR_PROJECT_ID")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Uniswap Pair 合约地址
	pairAddress := common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc") // USDC/WETH V2

	// 3. 实例化 Pair 合约
	pair, err := uniswapv2.NewUniswapV2Pair(pairAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 调用 getReserves
	reserves, err := pair.GetReserves(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Reserve0:", reserves.Reserve0)
	fmt.Println("Reserve1:", reserves.Reserve1)
	fmt.Println("LastUpdate:", reserves.BlockTimestampLast)
}
