// Package scanner 提供以太坊区块扫描能力。
//
// 重构要点：
//   - Scanner 通过构造函数注入依赖（ethclient、parser.Engine、logger）
//   - 去掉 Init() 方法，改为 New 时一次性完成初始化
//   - 提供 ScanBlock（单块）/ ScanBlocks（区间流式）两个核心方法
//   - 交易分类逻辑抽取为独立方法，可单独测试
package scanner

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ETHSentinel/internal/parser"
	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
)

// ─────────────────────────────────────────────
//  BlockResult — 单块扫描结果
// ─────────────────────────────────────────────

// BlockResult 一个区块的所有解析结果
type BlockResult struct {
	BlockNumber *big.Int
	Events      []comm.UnifiedEvent // 所有交易解析出的事件（扁平化）
	TxCount     int
	Err         error
}

// ─────────────────────────────────────────────
//  Scanner
// ─────────────────────────────────────────────

// Scanner 以太坊区块扫描器
type Scanner struct {
	client *ethclient.Client
	engine *parser.Engine
	logger *zap.Logger
	// addrBook 已知合约地址到分类的映射（可选，用于快速分类）
	addrBook map[common.Address]TxCategory
}

// New 创建扫描器
func New(client *ethclient.Client, engine *parser.Engine, logger *zap.Logger) *Scanner {
	return &Scanner{
		client:   client,
		engine:   engine,
		logger:   logger,
		addrBook: make(map[common.Address]TxCategory),
	}
}

// RegisterAddress 注册已知合约地址的分类（用于快速分类交易）
func (s *Scanner) RegisterAddress(addr common.Address, category TxCategory) {
	s.addrBook[addr] = category
}

// ─────────────────────────────────────────────
//  ScanBlock — 单块扫描
// ─────────────────────────────────────────────

// ScanBlock 扫描单个区块，返回所有解析出的链上事件。
// blockNumber 为 nil 时扫描最新块。
func (s *Scanner) ScanBlock(ctx context.Context, blockNumber *big.Int, cfg ScanBlockCfg) (*BlockResult, error) {
	block, err := s.client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("scanner: fetch block %v: %w", blockNumber, err)
	}

	// 若没有提供 ActiveParsers，尝试用默认配置
	active := cfg.Active
	chain := BuildChain(active)

	result := &BlockResult{
		BlockNumber: block.Number(),
		TxCount:     len(block.Transactions()),
	}

	// 并发获取 receipt 并解析事件
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i, tx := range block.Transactions() {
		txCopy := tx
		txIdx := uint(i)
		wg.Add(1)

		go func() {
			defer wg.Done()

			receipt, err := s.client.TransactionReceipt(ctx, txCopy.Hash())
			if err != nil {
				s.logger.Warn("scanner: get receipt failed",
					zap.String("tx", txCopy.Hash().Hex()),
					zap.Error(err),
				)
				return
			}

			evs := s.parseReceipt(receipt, txIdx, block.Time(), chain)

			mu.Lock()
			result.Events = append(result.Events, evs...)
			mu.Unlock()
		}()
	}
	wg.Wait()

	return result, nil
}

// ─────────────────────────────────────────────
//  ScanBlocks — 区间流式扫描
// ─────────────────────────────────────────────

// ScanBlocks 扫描 [startBlock, endBlock) 区间，通过 channel 流式返回结果。
// 调用方通过 for result := range ch 消费，cancel ctx 可提前终止。
func (s *Scanner) ScanBlocks(ctx context.Context, cfg ScanBlocksCfg) (<-chan *BlockResult, error) {
	if cfg.StartBlock == nil || cfg.EndBlock == nil {
		return nil, fmt.Errorf("scanner: start/end block must not be nil")
	}
	if cfg.StartBlock.Cmp(cfg.EndBlock) >= 0 {
		return nil, fmt.Errorf("scanner: start block must be less than end block")
	}

	poolSize := cfg.WorkerPoolSize
	if poolSize <= 0 {
		poolSize = 100
	}

	// 构建 ActiveParsers
	var active *comm.ActiveParsers
	if len(cfg.ParserCfg) > 0 && s.engine != nil {
		var err error
		active, err = s.engine.BuildActive(cfg.ParserCfg)
		if err != nil {
			return nil, fmt.Errorf("scanner: build active parsers: %w", err)
		}
	}

	out := make(chan *BlockResult, poolSize)

	pool, err := ants.NewPool(poolSize,
		ants.WithExpiryDuration(30*time.Second),
		ants.WithNonblocking(false),
		ants.WithPreAlloc(true),
	)
	if err != nil {
		return nil, fmt.Errorf("scanner: create goroutine pool: %w", err)
	}

	go func() {
		defer close(out)
		defer pool.Release()

		var wg sync.WaitGroup
		blockCfg := ScanBlockCfg{Active: active}

		for i := new(big.Int).Set(cfg.StartBlock); i.Cmp(cfg.EndBlock) < 0; i.Add(i, big.NewInt(1)) {
			if ctx.Err() != nil {
				break
			}
			bn := new(big.Int).Set(i)
			wg.Add(1)

			if submitErr := pool.Submit(func() {
				defer wg.Done()
				res, err := s.ScanBlock(ctx, bn, blockCfg)
				if err != nil {
					log.Printf("scanner: block %s error: %v", bn, err)
					out <- &BlockResult{BlockNumber: bn, Err: err}
					return
				}
				out <- res
			}); submitErr != nil {
				wg.Done()
				s.logger.Error("scanner: submit task failed", zap.Error(submitErr))
			}
		}
		wg.Wait()
	}()

	return out, nil
}

// ─────────────────────────────────────────────
//  内部：解析单笔交易的所有日志
// ─────────────────────────────────────────────

func (s *Scanner) parseReceipt(
	receipt *types.Receipt,
	txIdx uint,
	blockTime uint64,
	chain ChainNode,
) []comm.UnifiedEvent {

	if chain == nil || len(receipt.Logs) == 0 {
		return nil
	}

	meta := comm.EventMetadata{
		TxHash:           receipt.TxHash,
		BlockNumber:      receipt.BlockNumber,
		TransactionIndex: txIdx,
		Age:              time.Unix(int64(blockTime), 0),
	}

	var evs []comm.UnifiedEvent
	for i, lg := range receipt.Logs {
		meta.OuterIndex = uint(i)
		if ev, ok := chain.Handle(*lg, meta); ok {
			evs = append(evs, ev)
		}
	}
	return evs
}

// ─────────────────────────────────────────────
//  内部：交易分类
// ─────────────────────────────────────────────

// ClassifyTx 判断交易类型（普通转账 / Token / DeFi / 合约部署）
func (s *Scanner) ClassifyTx(tx *types.Transaction) TxCategory {
	to := tx.To()
	if to == nil {
		return TxCategoryNewContract
	}
	if len(tx.Data()) == 0 {
		return TxCategoryETH
	}
	if cat, ok := s.addrBook[*to]; ok {
		return cat
	}
	return TxCategoryDeFi // 默认视为 DeFi 交互
}
