// Package sentinel 是 ETH Sentinel 的对外入口。
//
// 两种使用模式：
//
//  1. 嵌入模式（直接集成到业务进程）：
//     client, err := client.New(client.WithRPCURL("https://..."), client.WithWSURL("wss://..."))
//
//  2. 远程 gRPC 模式（连接独立部署的 Sentinel Server）：
//     client, err := client.NewRemote("localhost:50051")
//
// 两种模式暴露完全相同的 API 接口。
package sentinel

import (
	"context"
	"github.com/ETHSentinel/internal/parser"
	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ETHSentinel/internal/scanner"
	"math/big"

	sentinelv1 "github.com/ETHSentinel/gen/sentinel/v1"
	"github.com/ETHSentinel/internal/conn"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ─────────────────────────────────────────────
//  公开类型别名（调用方无需 import internal）
// ─────────────────────────────────────────────

// Event 统一链上事件（等同 comm.UnifiedEvent）
type Event = comm.UnifiedEvent

// ParserCfg 解析配置（等同 comm.ParserCfg）
type ParserCfg = comm.ParserCfg

// ProtocolType / ProtocolImpl / EventMethod 常量透出
const (
	ProtocolTypeDEX     = comm.ProtocolTypeDEX
	ProtocolTypeToken   = comm.ProtocolTypeToken
	ProtocolTypeLending = comm.ProtocolTypeLending

	ProtocolImplUniswapV2 = comm.ProtocolImplUniswapV2
	ProtocolImplERC20     = comm.ProtocolImplERC20
	ProtocolImplERC721    = comm.ProtocolImplERC721

	EventMethodSwap     = comm.EventMethodSwap
	EventMethodTransfer = comm.EventMethodTransfer
)

// ─────────────────────────────────────────────
//  ScanResult — 扫描结果
// ─────────────────────────────────────────────

// ScanResult 单块扫描结果
type ScanResult struct {
	BlockNumber *big.Int
	Events      []Event
	TxCount     int
}

// ─────────────────────────────────────────────
//  Client 接口
// ─────────────────────────────────────────────

// Client SDK 核心接口
type Client interface {
	// ScanBlock 扫描单个区块
	ScanBlock(ctx context.Context, blockNumber *big.Int, cfg ParserCfg) (*ScanResult, error)

	// ScanBlocks 扫描区间，流式返回每个块的结果
	ScanBlocks(ctx context.Context, start, end *big.Int, cfg ParserCfg) (<-chan *ScanResult, error)

	// SubscribeBlocks 实时订阅新块事件（需要 WS 连接）
	SubscribeBlocks(ctx context.Context, cfg ParserCfg) (<-chan *ScanResult, error)

	// Close 释放资源
	Close() error
}

// ─────────────────────────────────────────────
//  localClient — 嵌入模式实现
// ─────────────────────────────────────────────

type localClient struct {
	sc      *scanner.Scanner
	connMgr *conn.Manager
	opts    *options
	logger  *zap.Logger
}

// New 创建嵌入模式 SDK 客户端（直接连接以太坊节点）
func New(opts ...Option) (Client, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	logger := o.logger

	// 初始化连接管理器
	connMgr := conn.NewManager(logger)
	nodeCfg := &conn.NodeConfig{
		Name:     "default",
		RPCURL:   o.rpcURL,
		WSURL:    o.wsURL,
		ProxyURL: o.proxyURL,
	}
	connMgr.Register(nodeCfg)

	// 获取 RPC 连接
	ctx := context.Background()
	ethClient, err := connMgr.Get(ctx, "default", conn.MethodRPC)
	if err != nil {
		return nil, err
	}

	// 初始化解析引擎
	engine, err := parser.NewEngine(ethClient)
	if err != nil {
		return nil, err
	}

	sc := scanner.New(ethClient, engine, logger)

	return &localClient{
		sc:      sc,
		connMgr: connMgr,
		opts:    o,
		logger:  logger,
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
	return &ScanResult{
		BlockNumber: res.BlockNumber,
		Events:      res.Events,
		TxCount:     res.TxCount,
	}, nil
}

func (c *localClient) ScanBlocks(ctx context.Context, start, end *big.Int, cfg ParserCfg) (<-chan *ScanResult, error) {
	ch, err := c.sc.ScanBlocks(ctx, scanner.ScanBlocksCfg{
		StartBlock:     start,
		EndBlock:       end,
		ParserCfg:      cfg,
		WorkerPoolSize: c.opts.workerPoolSize,
	})
	if err != nil {
		return nil, err
	}
	out := make(chan *ScanResult, 10)
	go func() {
		defer close(out)
		for res := range ch {
			out <- &ScanResult{
				BlockNumber: res.BlockNumber,
				Events:      res.Events,
				TxCount:     res.TxCount,
			}
		}
	}()
	return out, nil
}

func (c *localClient) SubscribeBlocks(ctx context.Context, cfg ParserCfg) (<-chan *ScanResult, error) {
	// TODO: 使用 ethclient.SubscribeNewHead + ScanBlock 实现实时订阅
	return nil, nil
}

func (c *localClient) Close() error {
	c.connMgr.Close()
	return nil
}

// ─────────────────────────────────────────────
//  remoteClient — gRPC 远程模式实现
// ─────────────────────────────────────────────

type remoteClient struct {
	conn   *grpc.ClientConn
	stub   sentinelv1.SentinelServiceClient
	opts   *options
	logger *zap.Logger
}

// NewRemote 创建 gRPC 远程模式客户端（连接独立部署的 Sentinel Server）
func NewRemote(grpcAddr string, opts ...Option) (Client, error) {
	o := defaultOptions()
	o.grpcAddr = grpcAddr
	for _, opt := range opts {
		opt(o)
	}

	grpcConn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &remoteClient{
		conn:   grpcConn,
		stub:   sentinelv1.NewSentinelServiceClient(grpcConn),
		opts:   o,
		logger: o.logger,
	}, nil
}

func (c *remoteClient) ScanBlock(ctx context.Context, blockNumber *big.Int, cfg ParserCfg) (*ScanResult, error) {
	resp, err := c.stub.ScanBlock(ctx, &sentinelv1.ScanBlockRequest{
		BlockNumber:       blockNumber.String(),
		SelectedProtocols: parserCfgToProto(cfg),
	})
	if err != nil {
		return nil, err
	}
	bn := new(big.Int)
	bn.SetString(resp.BlockNumber, 10)
	return &ScanResult{
		BlockNumber: bn,
		TxCount:     int(resp.TxCount),
		// Events 转换省略（需要实现 fromProtoEvent）
	}, nil
}

func (c *remoteClient) ScanBlocks(ctx context.Context, start, end *big.Int, cfg ParserCfg) (<-chan *ScanResult, error) {
	stream, err := c.stub.ScanBlocks(ctx, &sentinelv1.ScanBlocksRequest{
		StartBlock:        start.String(),
		EndBlock:          end.String(),
		SelectedProtocols: parserCfgToProto(cfg),
		WorkerPoolSize:    int32(c.opts.workerPoolSize),
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
			out <- &ScanResult{
				BlockNumber: bn,
				TxCount:     int(resp.TxCount),
			}
		}
	}()
	return out, nil
}

func (c *remoteClient) SubscribeBlocks(ctx context.Context, cfg ParserCfg) (<-chan *ScanResult, error) {
	stream, err := c.stub.SubscribeBlocks(ctx, &sentinelv1.SubscribeRequest{
		SelectedProtocols: parserCfgToProto(cfg),
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

func (c *remoteClient) Close() error {
	return c.conn.Close()
}

// ─────────────────────────────────────────────
//  辅助：ParserCfg → proto selectors
// ─────────────────────────────────────────────

func parserCfgToProto(cfg ParserCfg) []*sentinelv1.ProtocolSelector {
	var selectors []*sentinelv1.ProtocolSelector
	for pt, implCfg := range cfg {
		for impl, methods := range implCfg {
			sel := &sentinelv1.ProtocolSelector{
				ProtocolType: internalTypeToProto(pt),
				ProtocolImpl: internalImplToProto(impl),
			}
			for _, m := range methods {
				sel.Events = append(sel.Events, internalMethodToProto(m))
			}
			selectors = append(selectors, sel)
		}
	}
	return selectors
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
