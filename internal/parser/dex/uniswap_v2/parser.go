// Package uniswapv2 实现 UniswapV2 协议的事件解析。
//
// 重构要点：
//   - 修复原代码中 ParseSwapEvent 每次调用都重新创建 filterer 的性能问题
//   - filterer 在 Parser 初始化时一次性创建并缓存
//   - Parser 通过构造函数注入 ethclient，而非在函数内部隐式获取连接
//   - 实现 comm.ProtocolImplParser 接口
package uniswapv2

import (
	"fmt"
	"math/big"
	"sync"

	ablibgens "github.com/ETHSentinel/internal/lib/dex/uniswap_v2"
	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// ─────────────────────────────────────────────
//  预定义合约地址（可通过配置文件替换）
// ─────────────────────────────────────────────

// KnownPairs 已知的 UniswapV2 Pair 地址列表
// 实际生产中建议从配置/数据库动态加载
var KnownPairs = []common.Address{
	common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc"), // USDC/WETH
	common.HexToAddress("0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11"), // DAI/WETH
	common.HexToAddress("0x0d4a11d5EEaaC28EC3F61d100daF4d40471f1852"), // USDT/WETH
}

// ─────────────────────────────────────────────
//  Parser — 实现 comm.ProtocolImplParser
// ─────────────────────────────────────────────

// pairCache 缓存一个 Pair 合约相关的信息，避免重复链上调用
type pairCache struct {
	filterer *ablibgens.UniswappairFilterer
	token0   common.Address
	token1   common.Address
}

// Parser UniswapV2 事件解析器
type Parser struct {
	client  *ethclient.Client
	invoker *comm.EventParseInvoker
	logger  *zap.Logger

	// filtererCache key=pair合约地址
	mu            sync.RWMutex
	filtererCache map[common.Address]*pairCache
}

// NewParser 创建 UniswapV2 解析器（client 稍后通过 WithClient 注入，
// 支持懒初始化——在实际处理第一笔交易时才连接）
//func NewParser() *Parser {
//	p := &Parser{
//		invoker:       comm.NewEventParseInvoker(comm.ProtocolImplUniswapV2),
//		filtererCache: make(map[common.Address]*pairCache),
//		logger:        zap.NewNop(),
//	}
//	// 注册事件处理函数
//	p.invoker.RegisterOne(comm.SigUniswapV2Swap, comm.EventMethodSwap, p.parseSwap)
//	return p
//}

func NewParser(client ...*ethclient.Client) *Parser {
	p := &Parser{
		invoker:       comm.NewEventParseInvoker(comm.ProtocolImplUniswapV2),
		filtererCache: make(map[common.Address]*pairCache),
		logger:        zap.NewNop(),
	}
	if len(client) > 0 && client[0] != nil {
		p.client = client[0]
	}
	p.invoker.RegisterOne(comm.SigUniswapV2Swap, comm.EventMethodSwap, p.parseSwap)
	return p
}

// WithClient 注入以太坊客户端（可在 Scanner 初始化阶段调用）
func (p *Parser) WithClient(client *ethclient.Client) *Parser {
	p.client = client
	return p
}

// WithLogger 注入日志
func (p *Parser) WithLogger(logger *zap.Logger) *Parser {
	p.logger = logger
	return p
}

// WithKnownPairs 预热 filterer 缓存（可选，提升首次解析速度）
func (p *Parser) WithKnownPairs(pairs []common.Address) *Parser {
	for _, addr := range pairs {
		_, _ = p.ensureCache(addr) // 忽略错误，懒加载兜底
	}
	return p
}

// ── comm.ProtocolImplParser 接口实现 ──────────

func (p *Parser) HandleEvent(sig comm.EventSig, log types.Log, meta comm.EventMetadata) (comm.UnifiedEvent, error) {
	return p.invoker.HandleEvent(sig, log, meta)
}

func (p *Parser) ListEventSigs() []comm.EventSig {
	return p.invoker.ListEventSigs()
}

func (p *Parser) SetFilter(methods []comm.EventMethod) {
	p.invoker.SetFilter(methods)
}

// ─────────────────────────────────────────────
//  内部解析实现
// ─────────────────────────────────────────────

// parseSwap 解析单条 Swap log
func (p *Parser) parseSwap(log types.Log, meta comm.EventMetadata) (comm.UnifiedEvent, error) {
	cache, err := p.ensureCache(log.Address)
	if err != nil {
		return nil, fmt.Errorf("uniswapv2: get filterer for %s: %w", log.Address.Hex(), err)
	}

	swapEvent, err := cache.filterer.ParseSwap(log)
	if err != nil {
		return nil, fmt.Errorf("uniswapv2: parse swap log: %w", err)
	}

	token0Name := p.tokenSymbol(cache.token0)
	token1Name := p.tokenSymbol(cache.token1)

	// 确定方向：amount0In > 0 表示用 token0 换 token1
	var fromToken, toToken common.Address
	var fromAmt, toAmt *big.Int
	if swapEvent.Amount0In.Sign() > 0 {
		fromToken, toToken = cache.token0, cache.token1
		fromAmt, toAmt = swapEvent.Amount0In, swapEvent.Amount1Out
	} else {
		fromToken, toToken = cache.token1, cache.token0
		fromAmt, toAmt = swapEvent.Amount1In, swapEvent.Amount0Out
	}

	return &comm.UnifiedEventData{
		Metadata: comm.EventMetadata{
			TxHash:           meta.TxHash,
			ProtocolTypeVal:  comm.ProtocolTypeDEX,
			ProtocolImplVal:  comm.ProtocolImplUniswapV2,
			Age:              meta.Age,
			To:               log.Address,
			BlockNumber:      meta.BlockNumber,
			OuterIndex:       meta.OuterIndex,
			TransactionIndex: meta.TransactionIndex,
		},
		Base: comm.BaseEvent{
			EventType: comm.EventMethodSwap,
			From:      swapEvent.Sender,
			RefTokens: []comm.RefToken{
				{Name: token0Name, Amount: swapEvent.Amount0In},
				{Name: token1Name, Amount: swapEvent.Amount1Out},
			},
			RealValues: []decimal.Decimal{}, // 价格预言机可在此填充
		},
		DetailVal: &comm.SwapData{
			FromToken:  fromToken,
			ToToken:    toToken,
			FromAmount: fromAmt,
			ToAmount:   toAmt,
			Sender:     swapEvent.Sender,
			Recipient:  swapEvent.To,
		},
	}, nil
}

// ensureCache 懒加载 pair 的 filterer 和 token 信息
func (p *Parser) ensureCache(addr common.Address) (*pairCache, error) {
	p.mu.RLock()
	if cache, ok := p.filtererCache[addr]; ok {
		p.mu.RUnlock()
		return cache, nil
	}
	p.mu.RUnlock()

	if p.client == nil {
		return nil, fmt.Errorf("uniswapv2: ethclient not initialized (call WithClient first)")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// double-check
	if cache, ok := p.filtererCache[addr]; ok {
		return cache, nil
	}

	filterer, err := ablibgens.NewUniswappairFilterer(addr, p.client)
	if err != nil {
		return nil, err
	}

	pair, err := ablibgens.NewUniswappair(addr, p.client)
	if err != nil {
		return nil, err
	}

	token0, err := pair.Token0(nil)
	if err != nil {
		return nil, err
	}
	token1, err := pair.Token1(nil)
	if err != nil {
		return nil, err
	}

	cache := &pairCache{filterer: filterer, token0: token0, token1: token1}
	p.filtererCache[addr] = cache

	p.logger.Debug("uniswapv2: pair cached",
		zap.String("pair", addr.Hex()),
		zap.String("token0", token0.Hex()),
		zap.String("token1", token1.Hex()),
	)
	return cache, nil
}

// tokenSymbol 尝试获取代币符号，失败则返回地址缩写
func (p *Parser) tokenSymbol(addr common.Address) string {
	// TODO: 接入代币元数据缓存（如 token list / Redis）
	hex := addr.Hex()
	if len(hex) > 10 {
		return hex[:6] + "…" + hex[len(hex)-4:]
	}
	return hex
}
