package models

import (
	"context"
	"github.com/CryptoQuantX/chain_monitor/initialize"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

type BlockStruct struct {
	Header       BlockHeader
	Transactions types.Transactions
}

type BlockHeader struct {
	BlockNumber *big.Int
	Timestamp   uint64
	Miner       common.Address
	GasLimit    uint64
	GasUsed     uint64
	Number      *big.Int
}

func (b *BlockStruct) ParseEthBlock() (*BlockStruct, *ethclient.Client) {
	ethClient := initialize.InfuraConn("https://mainnet.infura.io/v3/0d79a9c32c814e1da6133850f6fa1128", "http://192.168.248.215:7890")

	header, err := ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("获取最新区块失败: %v", err)
	}

	b.Header.BlockNumber = header.Number
	b.Header.Timestamp = header.Time
	b.Header.Miner = header.Coinbase
	b.Header.GasLimit = header.GasLimit
	b.Header.GasUsed = header.GasUsed
	b.Header.Number = header.Number

	block, err := ethClient.BlockByNumber(context.Background(), b.Header.BlockNumber)
	if err != nil {
		log.Fatalf("获取区块失败: %v", err)
	}

	b.Transactions = block.Transactions()

	return b, ethClient
}
