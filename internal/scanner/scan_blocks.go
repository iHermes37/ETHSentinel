package scanner

import (
	"context"
	"github.com/Crypto-ChainSentinel/test/Filter"
	"github.com/panjf2000/ants/v2"
	"log"
	"math/big"
	"sync"
	"time"

	ParserEngineCommon "github.com/Crypto-ChainSentinel/internal/parser/comm"
)

func (s *Scanner) ScanBlocks(cfg Interval) chan [][]ParserEngineCommon.UnifiedEventData {

	handlePipeline := make(chan [][]ParserEngineCommon.UnifiedEventData, 10)
	//=========================协程池===========================

	// 自定义配置创建协程池
	pool, err := ants.NewPool(100, // 协程池容量
		ants.WithExpiryDuration(30*time.Second), // 空闲协程过期时间
		ants.WithNonblocking(true),              // 非阻塞模式
		ants.WithPreAlloc(true),                 // 预分配内存
		ants.WithMaxBlockingTasks(1000),         // 最大阻塞任务数
	)
	if err != nil {
		panic(err)
	}
	defer pool.Release()

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//===========================================================

	for i := new(big.Int).Set(cfg.StartBlock); i.Cmp(cfg.EndBlock) < 0; i.Add(i, big.NewInt(1)) {
		blockNumber := new(big.Int).Set(i)
		wg.Add(1)

		_ = pool.Submit(func() {
			defer wg.Done()

			block, err := s.Client.BlockByNumber(ctx, blockNumber)
			if err != nil {
				log.Printf("读取区块 %s 失败: %v", blockNumber.String(), err)
				return
			}

			scanTxCfg := Filter.ScanTransConfig{
				BeforFilter:       cfg.FilterCfg,
				SelectedProtocols: cfg.Selected,
			}
			evlists := s.ScanBlock(block, scanTxCfg)
			handlePipeline <- evlists
		})
	}

	// 等待所有 goroutine 完成，然后关闭通道
	go func() {
		wg.Wait()
		close(handlePipeline)
		pool.Release()
	}()

	return handlePipeline
}
