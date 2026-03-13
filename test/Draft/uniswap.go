package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 1. 连接以太坊节点
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/YOUR_PROJECT_ID")
	if err != nil {
		log.Fatal(err)
	}

	// 2. UniswapV2Protocol Pair 合约地址
	pairAddress := common.HexToAddress("0x...") // 替换为 Pair 地址

	// 3. Pair ABI 中 getReserves 函数定义
	pairABI, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"internalType":"uint112","name":"_reserve0","type":"uint112"},{"internalType":"uint112","name":"_reserve1","type":"uint112"},{"internalType":"uint32","name":"_blockTimestampLast","type":"uint32"}],"stateMutability":"view","type":"function"}]`))
	if err != nil {
		log.Fatal(err)
	}

	// 4. 构建调用
	data, err := pairABI.Pack("getReserves")
	if err != nil {
		log.Fatal(err)
	}

	// 5. 调用
	res, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &pairAddress,
		Data: data,
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 6. 解析返回值
	var reserve0, reserve1 *big.Int
	var ts uint32
	err = pairABI.UnpackIntoInterface(&[]interface{}{&reserve0, &reserve1, &ts}, "getReserves", res)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Reserve0:", reserve0)
	fmt.Println("Reserve1:", reserve1)
}
