package scanner

import (
	"context"
	"fmt"
	"sync"

	ParserEngineCommon "github.com/Crypto-ChainSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (s *Scanner) ScanBlock(
	block *types.Block,

	EthCh chan<- *types.Transaction,
	TokenCh chan<- *types.Receipt,
	DefiCh chan<- *types.Receipt,
	NewContractCh chan<- *types.Receipt,

) [][]ParserEngineCommon.UnifiedEvent {

	ctx := context.Background()

	resultPipeline := make(chan []ParserEngineCommon.UnifiedEvent, len(block.Transactions()))
	var wg sync.WaitGroup

	for _, tran := range block.Transactions() {
		txhash := tran.Hash()
		tx := tran
		to := tran.To()
		wg.Add(1)
		go func(tx *types.Transaction, txhash common.Hash) {
			defer wg.Done()

			tran_receipt, err := s.Client.TransactionReceipt(ctx, txhash)
			if err != nil {
				fmt.Println("Receipt error:", err)
				return
			}

			if to != nil {
				if len(tx.Data()) == 0 {
					// 普通转账
					EthCh <- tx
				} else if IsToken(to) {
					// 普通代币转账
					TokenCh <- tran_receipt
				} else if IsDefi(to) {
					// DeFi 交互
					DefiCh <- tran_receipt
				}
			} else {
				// 合约创建部署
				NewContractCh <- tran_receipt
			}

		}(tx, txhash)
	}

	// 等待所有 goroutine 完成，然后关闭通道
	go func() {
		wg.Wait()
		close(resultPipeline)
	}()

	// 收集结果
	var evlist [][]ParserEngineCommon.UnifiedEvent
	for ev := range resultPipeline {
		evlist = append(evlist, ev)
	}

	return evlist
}
