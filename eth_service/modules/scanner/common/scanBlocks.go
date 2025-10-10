package common

import (
	"context"
	"log"
	"math/big"
	"sync"

	"github.com/Crypto-ChainSentinel/modules/ConnectionManager"
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
	"github.com/ethereum/go-ethereum/common"
)

type ScanBlocksConfig struct {
	WhaleAddr  *common.Address
	StartBlock *big.Int
	EndBlock   *big.Int
	selected   *map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl
}

func ScanBlocks(cfg ScanBlocksConfig) chan [][]ParserEngineCommon.UnifiedEvent {
	var wg sync.WaitGroup
	client := ConnectionManager.InfuraConn()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//var globalEv []ParserEngineCommon.UnifiedEvent

	handlePipline := make(chan [][]ParserEngineCommon.UnifiedEvent, 10)
	sem := make(chan struct{}, 10) // 同时最多 10 个 goroutine

	for i := new(big.Int).Set(cfg.StartBlock); i.Cmp(cfg.EndBlock) < 0; i.Add(i, big.NewInt(1)) {
		blockNumber := new(big.Int).Set(i)
		wg.Add(1)
		sem <- struct{}{} // 占用一个槽位
		go func(blockNumber *big.Int) {
			defer wg.Done()
			defer func() { <-sem }() // 释放槽位

			block, err := client.BlockByNumber(ctx, blockNumber)
			if err != nil {
				log.Printf("读取区块 %d 失败: %v", blockNumber, err)
				return
			}

			scanTxCfg := ScanTransConfig{
				WhaleAddr:         cfg.WhaleAddr,
				SelectedProtocols: cfg.selected,
			}

			evlists := ScanBlock(block, scanTxCfg)
			handlePipline <- evlists
		}(blockNumber)
	}

	// 等待所有 goroutine 完成，然后关闭通道
	go func() {
		wg.Wait()
		close(handlePipline)
	}()

	// go Handle.HandleEvents(handlePipline)
	return handlePipline
}
