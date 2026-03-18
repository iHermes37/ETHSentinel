// Package sentinel 是 ETH Sentinel 的对外 SDK 入口。
package sentinel

import (
	"context"
	"fmt"
	"math/big"

	sentinelv1 "github.com/ETHSentinel/gen/sentinel/v1"
	"github.com/ETHSentinel/internal/chain"
	"github.com/ETHSentinel/internal/conn"
	"github.com/ETHSentinel/internal/parser"
	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ETHSentinel/internal/scanner"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Event = comm.UnifiedEvent
type ParserCfg = comm.ParserCfg
type SwapData = comm.SwapData
type TransferData = comm.TransferData

const (
	ChainETH      = chain.ChainETH
	ChainBSC      = chain.ChainBSC
	ChainPolygon  = chain.ChainPolygon
	ChainArbitrum = chain.ChainArbitrum

	ProtocolTypeDEX     = comm.ProtocolTypeDEX
	ProtocolTypeToken   = comm.ProtocolTypeToken
	ProtocolTypeLending = comm.ProtocolTypeLending

	ProtocolImplUniswapV2 = comm.ProtocolImplUniswapV2
	ProtocolImplERC20     = comm.ProtocolImplERC20
	ProtocolImplERC721    = comm.ProtocolImplERC721

	EventMethodSwap     = comm.EventMethodSwap
	EventMethodTransfer = comm.EventMethodTransfer
)

type ScanResult struct {
	BlockNumber *big.Int
	ChainID     uint64
	Events      []Event
	TxCount     int
}

type Client interface {
	ScanBlock(ctx context.Context, blockNumber *big.Int, cfg ParserCfg) (*ScanResult, error)
	ScanBlocks(ctx context.Context, start, end *big.Int, cfg ParserCfg) (<-chan *ScanResult, error)
	SubscribeBlocks(ctx context.Context, cfg ParserCfg) (<-chan *ScanResult, error)
	Mempool() *MempoolClient
	Wallet() (*WalletClient, error)
	Close() error
}

// ── localClient ──────────────────────────────

type localClient struct {
	sc       *scanner.Scanner
	pool     *conn.MultiChainPool
	connMgr  *conn.Manager
	chainObj chain.Chain
	opts     *options
	logger   *zap.Logger
}

func New(opts ...Option) (Client, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}
	logger := o.logger

	chainObj, err := chain.Get(chain.ChainID(o.chainID))
	if err != nil {
		return nil, fmt.Errorf("sentinel: %w", err)
	}

	pool := conn.NewMultiChainPool(logger)
	nodeCfg := &conn.NodeConfig{
		Name:     chainObj.Name(),
		RPCURL:   chainObj.DefaultRPCURL(),
		WSURL:    chainObj.DefaultWSURL(),
		ProxyURL: o.proxyURL,
	}
	if o.rpcURL != "" {
		nodeCfg.RPCURL = o.rpcURL
	}
	if o.wsURL != "" {
		nodeCfg.WSURL = o.wsURL
	}
	pool.RegisterChain(chainObj, nodeCfg)

	connMgr := conn.NewManager(logger)
	connMgr.Register(nodeCfg)

	ctx := context.Background()
	ethClient, err := connMgr.Get(ctx, chainObj.Name(), conn.MethodRPC)
	if err != nil {
		return nil, fmt.Errorf("sentinel: connect rpc: %w", err)
	}

	engine, err := parser.NewEngine(ethClient)
	if err != nil {
		return nil, fmt.Errorf("sentinel: init engine: %w", err)
	}

	return &localClient{
		sc:       scanner.New(ethClient, engine, logger),
		pool:     pool,
		connMgr:  connMgr,
		chainObj: chainObj,
		opts:     o,
		logger:   logger,
	}, nil
}

func (c *localClient) ScanBlock(ctx context.Context, blockNumber *big.Int, cfg ParserCfg) (*ScanResult, error) {
	active, err := c.sc.Engine().BuildActive(cfg)
	if err != nil {
		return nil, err
	}
	res, err := c.sc.ScanBlock(ctx, blockNumber, scanner.ScanBlockCfg{Active: active})
	if err != nil {
		return nil, err
	}
	return &ScanResult{BlockNumber: res.BlockNumber, ChainID: c.opts.chainID, Events: res.Events, TxCount: res.TxCount}, nil
}

func (c *localClient) ScanBlocks(ctx context.Context, start, end *big.Int, cfg ParserCfg) (<-chan *ScanResult, error) {
	ch, err := c.sc.ScanBlocks(ctx, scanner.ScanBlocksCfg{
		StartBlock: start, EndBlock: end, ParserCfg: cfg, WorkerPoolSize: c.opts.workerPoolSize,
	})
	if err != nil {
		return nil, err
	}
	out := make(chan *ScanResult, 10)
	go func() {
		defer close(out)
		for res := range ch {
			out <- &ScanResult{BlockNumber: res.BlockNumber, ChainID: c.opts.chainID, Events: res.Events, TxCount: res.TxCount}
		}
	}()
	return out, nil
}

func (c *localClient) SubscribeBlocks(_ context.Context, _ ParserCfg) (<-chan *ScanResult, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (c *localClient) Mempool() *MempoolClient {
	wsURL := c.opts.wsURL
	if wsURL == "" {
		wsURL = c.chainObj.DefaultWSURL()
	}
	wsClient, _ := c.pool.GetWS(context.Background(), c.chainObj.ID())
	return newMempoolClient(wsURL, wsClient, c.chainObj, c.logger)
}

func (c *localClient) Wallet() (*WalletClient, error) {
	if c.opts.mnemonic == "" {
		return nil, fmt.Errorf("sentinel: wallet requires mnemonic, use WithMnemonic()")
	}
	rpcClient, err := c.pool.GetRPC(context.Background(), c.chainObj.ID())
	if err != nil {
		return nil, err
	}
	return newWalletClient(c.opts.mnemonic, rpcClient, new(big.Int).SetUint64(c.opts.chainID), c.logger)
}

func (c *localClient) Close() error {
	c.pool.Close()
	c.connMgr.Close()
	return nil
}

// ── remoteClient ─────────────────────────────

type remoteClient struct {
	grpcConn     *grpc.ClientConn
	sentinelStub sentinelv1.SentinelServiceClient
	mempoolStub  sentinelv1.MempoolServiceClient
	walletStub   sentinelv1.WalletServiceClient
	opts         *options
	logger       *zap.Logger
}

func NewRemote(grpcAddr string, opts ...Option) (Client, error) {
	o := defaultOptions()
	o.grpcAddr = grpcAddr
	for _, opt := range opts {
		opt(o)
	}
	gc, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &remoteClient{
		grpcConn:     gc,
		sentinelStub: sentinelv1.NewSentinelServiceClient(gc),
		mempoolStub:  sentinelv1.NewMempoolServiceClient(gc),
		walletStub:   sentinelv1.NewWalletServiceClient(gc),
		opts:         o,
		logger:       o.logger,
	}, nil
}

func (c *remoteClient) ScanBlock(ctx context.Context, blockNumber *big.Int, cfg ParserCfg) (*ScanResult, error) {
	resp, err := c.sentinelStub.ScanBlock(ctx, &sentinelv1.ScanBlockRequest{
		BlockNumber: blockNumber.String(), SelectedProtocols: parserCfgToProto(cfg),
	})
	if err != nil {
		return nil, err
	}
	bn := new(big.Int)
	bn.SetString(resp.BlockNumber, 10)
	return &ScanResult{BlockNumber: bn, TxCount: int(resp.TxCount)}, nil
}

func (c *remoteClient) ScanBlocks(ctx context.Context, start, end *big.Int, cfg ParserCfg) (<-chan *ScanResult, error) {
	stream, err := c.sentinelStub.ScanBlocks(ctx, &sentinelv1.ScanBlocksRequest{
		StartBlock: start.String(), EndBlock: end.String(),
		SelectedProtocols: parserCfgToProto(cfg), WorkerPoolSize: int32(c.opts.workerPoolSize),
	})
	if err != nil {
		return nil, err
	}
	out := make(chan *ScanResult, 10)
	go func() {
		defer close(out)
		for {
			resp, err := stream.Recv()
			if err != nil {
				return
			}
			bn := new(big.Int)
			bn.SetString(resp.BlockNumber, 10)
			out <- &ScanResult{BlockNumber: bn, TxCount: int(resp.TxCount)}
		}
	}()
	return out, nil
}

func (c *remoteClient) SubscribeBlocks(ctx context.Context, cfg ParserCfg) (<-chan *ScanResult, error) {
	stream, err := c.sentinelStub.SubscribeBlocks(ctx, &sentinelv1.SubscribeRequest{SelectedProtocols: parserCfgToProto(cfg)})
	if err != nil {
		return nil, err
	}
	out := make(chan *ScanResult, 10)
	go func() {
		defer close(out)
		for {
			resp, err := stream.Recv()
			if err != nil {
				return
			}
			bn := new(big.Int)
			bn.SetString(resp.BlockNumber, 10)
			out <- &ScanResult{BlockNumber: bn, TxCount: int(resp.TxCount)}
		}
	}()
	return out, nil
}

func (c *remoteClient) Mempool() *MempoolClient {
	return &MempoolClient{mempoolStub: c.mempoolStub, chainID: c.opts.chainID, logger: c.logger}
}

func (c *remoteClient) Wallet() (*WalletClient, error) {
	if c.opts.mnemonic == "" {
		return nil, fmt.Errorf("sentinel: wallet requires mnemonic")
	}
	return &WalletClient{
		mnemonic:   c.opts.mnemonic,
		walletStub: c.walletStub,
		chainID:    new(big.Int).SetUint64(c.opts.chainID),
		logger:     c.logger,
	}, nil
}

func (c *remoteClient) Close() error { return c.grpcConn.Close() }

// ── Proto 转换 ────────────────────────────────

func parserCfgToProto(cfg ParserCfg) []*sentinelv1.ProtocolSelector {
	var out []*sentinelv1.ProtocolSelector
	for pt, implCfg := range cfg {
		for impl, methods := range implCfg {
			sel := &sentinelv1.ProtocolSelector{
				ProtocolType: internalTypeToProto(pt),
				ProtocolImpl: internalImplToProto(impl),
			}
			for _, m := range methods {
				sel.Events = append(sel.Events, internalMethodToProto(m))
			}
			out = append(out, sel)
		}
	}
	return out
}

func internalTypeToProto(pt comm.ProtocolType) sentinelv1.ProtocolType {
	switch pt {
	case comm.ProtocolTypeDEX:
		return sentinelv1.ProtocolType_PROTOCOL_TYPE_DEX
	case comm.ProtocolTypeToken:
		return sentinelv1.ProtocolType_PROTOCOL_TYPE_TOKEN
	default:
		return sentinelv1.ProtocolType_PROTOCOL_TYPE_UNSPECIFIED
	}
}

func internalImplToProto(pi comm.ProtocolImpl) sentinelv1.ProtocolImpl {
	switch pi {
	case comm.ProtocolImplUniswapV2:
		return sentinelv1.ProtocolImpl_PROTOCOL_IMPL_UNISWAP_V2
	case comm.ProtocolImplERC20:
		return sentinelv1.ProtocolImpl_PROTOCOL_IMPL_ERC20
	default:
		return sentinelv1.ProtocolImpl_PROTOCOL_IMPL_UNSPECIFIED
	}
}

func internalMethodToProto(m comm.EventMethod) sentinelv1.EventMethod {
	switch m {
	case comm.EventMethodSwap:
		return sentinelv1.EventMethod_EVENT_METHOD_SWAP
	case comm.EventMethodTransfer:
		return sentinelv1.EventMethod_EVENT_METHOD_TRANSFER
	default:
		return sentinelv1.EventMethod_EVENT_METHOD_UNSPECIFIED
	}
}
