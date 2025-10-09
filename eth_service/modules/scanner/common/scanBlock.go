package common

import (
	"context"
	"fmt"
	"github.com/Crypto-ChainSentinel/modules/ConnectionManager"
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"sync"
)

func ScanBlock(block *types.Block, cfg ScanConfig) []ParserEngineCommon.UnifiedEvent {
	client := ConnectionManager.InfuraConn()
	ctx := context.Background()

	resultPipeline := make(chan ParserEngineCommon.UnifiedEvent, len(block.Transactions()))
	var wg sync.WaitGroup

	for _, trans := range block.Transactions() {
		txhash := trans.Hash()
		tx := trans
		wg.Add(1)
		go func(tx *types.Transaction, txhash common.Hash) {
			defer wg.Done()
			tranreceipt, err := client.TransactionReceipt(ctx, txhash)
			if err != nil {
				fmt.Println("Receipt error:", err)
				return
			}
			if !ParserFilter(tx, FilterSetting{}) {
				ev := ParseTranByLog(tranreceipt, *cfg.SelectedProtocols)
				resultPipeline <- ev
			}

		}(tx, txhash)
	}

	// 等待所有 goroutine 完成，然后关闭通道
	go func() {
		wg.Wait()
		close(resultPipeline)
	}()

	// 收集结果
	var evlist []ParserEngineCommon.UnifiedEvent
	for ev := range resultPipeline {
		evlist = append(evlist, ev)
	}
	return evlist
}
