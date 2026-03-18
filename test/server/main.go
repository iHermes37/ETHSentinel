// Package grpc — SentinelService 的 gRPC 处理器实现
package grpc

import (
	"context"
	"fmt"
	"math/big"

	sentinelv1 "github.com/ETHSentinel/gen/sentinel/v1"
	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ETHSentinel/internal/scanner"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SentinelHandler 实现 sentinelv1.SentinelServiceServer
type SentinelHandler struct {
	sentinelv1.UnimplementedSentinelServiceServer
	scanner *scanner.Scanner
	client  *ethclient.Client
	logger  *zap.Logger
}

// NewSentinelHandler 创建处理器
func NewSentinelHandler(sc *scanner.Scanner, client *ethclient.Client, logger *zap.Logger) *SentinelHandler {
	return &SentinelHandler{scanner: sc, client: client, logger: logger}
}

// ─────────────────────────────────────────────
//  ScanBlock — 一元 RPC
// ─────────────────────────────────────────────

func (h *SentinelHandler) ScanBlock(ctx context.Context, req *sentinelv1.ScanBlockRequest) (*sentinelv1.ScanBlockResponse, error) {
	blockNumber := new(big.Int)
	if _, ok := blockNumber.SetString(req.BlockNumber, 10); !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid block number: %q", req.BlockNumber)
	}

	active, err := h.buildActive(req.SelectedProtocols)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "build parsers: %v", err)
	}

	res, err := h.scanner.ScanBlock(ctx, blockNumber, scanner.ScanBlockCfg{Active: active})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "scan block: %v", err)
	}

	pbEvents := make([]*sentinelv1.UnifiedEvent, 0, len(res.Events))
	for _, ev := range res.Events {
		pbEvents = append(pbEvents, toProtoEvent(ev))
	}

	return &sentinelv1.ScanBlockResponse{
		BlockNumber: res.BlockNumber.String(),
		Events:      pbEvents,
		TxCount:     uint32(res.TxCount),
		EventCount:  uint32(len(res.Events)),
	}, nil
}

// ─────────────────────────────────────────────
//  ScanBlocks — 服务端流式 RPC
// ─────────────────────────────────────────────

func (h *SentinelHandler) ScanBlocks(req *sentinelv1.ScanBlocksRequest, stream sentinelv1.SentinelService_ScanBlocksServer) error {
	startBlock := new(big.Int)
	endBlock := new(big.Int)
	if _, ok := startBlock.SetString(req.StartBlock, 10); !ok {
		return status.Errorf(codes.InvalidArgument, "invalid start block: %q", req.StartBlock)
	}
	if _, ok := endBlock.SetString(req.EndBlock, 10); !ok {
		return status.Errorf(codes.InvalidArgument, "invalid end block: %q", req.EndBlock)
	}

	parserCfg := protoSelectorsToParserCfg(req.SelectedProtocols)
	poolSize := int(req.WorkerPoolSize)

	ch, err := h.scanner.ScanBlocks(stream.Context(), scanner.ScanBlocksCfg{
		StartBlock:     startBlock,
		EndBlock:       endBlock,
		ParserCfg:      parserCfg,
		WorkerPoolSize: poolSize,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "start scan: %v", err)
	}

	for res := range ch {
		if res.Err != nil {
			h.logger.Warn("scan block error", zap.String("block", res.BlockNumber.String()), zap.Error(res.Err))
			continue
		}
		pbEvents := make([]*sentinelv1.UnifiedEvent, 0, len(res.Events))
		for _, ev := range res.Events {
			pbEvents = append(pbEvents, toProtoEvent(ev))
		}
		if err := stream.Send(&sentinelv1.ScanBlockResponse{
			BlockNumber: res.BlockNumber.String(),
			Events:      pbEvents,
			TxCount:     uint32(res.TxCount),
			EventCount:  uint32(len(res.Events)),
		}); err != nil {
			return err
		}
	}
	return nil
}

// ─────────────────────────────────────────────
//  SubscribeBlocks — 实时订阅新块（服务端流式）
// ─────────────────────────────────────────────

func (h *SentinelHandler) SubscribeBlocks(req *sentinelv1.SubscribeRequest, stream sentinelv1.SentinelService_SubscribeBlocksServer) error {
	ctx := stream.Context()

	// 修复错误一：SubscribeNewHead 需要 chan<- *types.Header，不是 chan interface{}
	headCh := make(chan *types.Header, 10)
	sub, err := h.client.SubscribeNewHead(ctx, headCh)
	if err != nil {
		return status.Errorf(codes.Unavailable, "subscribe new head: %v", err)
	}
	defer sub.Unsubscribe()

	active, err := h.buildActive(req.SelectedProtocols)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "build parsers: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-sub.Err():
			return status.Errorf(codes.Internal, "subscription error: %v", err)
		case header, ok := <-headCh:
			if !ok {
				return nil
			}
			res, err := h.scanner.ScanBlock(ctx, header.Number, scanner.ScanBlockCfg{Active: active})
			if err != nil {
				h.logger.Warn("subscribe: scan block error", zap.Error(err))
				continue
			}
			pbEvents := make([]*sentinelv1.UnifiedEvent, 0, len(res.Events))
			for _, ev := range res.Events {
				pbEvents = append(pbEvents, toProtoEvent(ev))
			}
			if err := stream.Send(&sentinelv1.ScanBlockResponse{
				BlockNumber: res.BlockNumber.String(),
				Events:      pbEvents,
				TxCount:     uint32(res.TxCount),
				EventCount:  uint32(len(res.Events)),
			}); err != nil {
				return err
			}
		}
	}
}

// ─────────────────────────────────────────────
//  HealthCheck
// ─────────────────────────────────────────────

func (h *SentinelHandler) HealthCheck(ctx context.Context, _ *emptypb.Empty) (*sentinelv1.HealthResponse, error) {
	latest, err := h.client.BlockNumber(ctx)
	if err != nil {
		return &sentinelv1.HealthResponse{Healthy: false}, nil
	}
	return &sentinelv1.HealthResponse{
		Healthy:      true,
		LatestBlock:  fmt.Sprintf("%d", latest),
		NodeEndpoint: "connected",
	}, nil
}

// ─────────────────────────────────────────────
//  辅助：Proto ↔ 内部类型转换
// ─────────────────────────────────────────────

func toProtoEvent(ev comm.UnifiedEvent) *sentinelv1.UnifiedEvent {
	meta := &sentinelv1.EventMetadata{
		TxHash: ev.GetTxHash().Hex(),
		// 修复错误二：string 类型不能直接转 proto enum，改用映射函数
		ProtocolType:     internalTypeToProto(ev.GetProtocolType()),
		ProtocolImpl:     internalImplToProto(ev.GetProtocolImpl()),
		Age:              timestamppb.New(ev.GetAge()),
		To:               ev.GetTo().Hex(),
		BlockNumber:      ev.GetBlockNumber().String(),
		OuterIndex:       uint32(ev.GetOuterIndex()),
		TransactionIndex: uint32(ev.GetTransactionIndex()),
	}

	base := ev.GetBase()
	pbBase := &sentinelv1.BaseEvent{
		// 修复错误三：同上，EventMethod string 不能直接转 proto enum
		EventType: internalMethodToProto(base.EventType),
		From:      base.From.Hex(),
	}
	for _, rt := range base.RefTokens {
		pbBase.RefTokens = append(pbBase.RefTokens, &sentinelv1.RefToken{
			Name:   rt.Name,
			Amount: rt.Amount.String(),
		})
	}

	pbEv := &sentinelv1.UnifiedEvent{
		Metadata: meta,
		Base:     pbBase,
	}

	switch d := ev.GetDetail().(type) {
	case *comm.SwapData:
		pbEv.Detail = &sentinelv1.UnifiedEvent_Swap{
			Swap: &sentinelv1.SwapDetail{
				FromToken:  d.FromToken.Hex(),
				ToToken:    d.ToToken.Hex(),
				FromAmount: d.FromAmount.String(),
				ToAmount:   d.ToAmount.String(),
				Sender:     d.Sender.Hex(),
				Recipient:  d.Recipient.Hex(),
			},
		}
	case *comm.TransferData:
		pbEv.Detail = &sentinelv1.UnifiedEvent_Transfer{
			Transfer: &sentinelv1.TransferDetail{
				Token:  d.Token.Hex(),
				From:   d.From.Hex(),
				To:     d.To.Hex(),
				Amount: d.Amount.String(),
			},
		}
	}

	return pbEv
}

func protoSelectorsToParserCfg(selectors []*sentinelv1.ProtocolSelector) comm.ParserCfg {
	cfg := make(comm.ParserCfg)
	for _, sel := range selectors {
		pt := protoTypeToInternal(sel.ProtocolType)
		pi := protoImplToInternal(sel.ProtocolImpl)
		if _, ok := cfg[pt]; !ok {
			cfg[pt] = make(comm.ImplConfig)
		}
		var methods []comm.EventMethod
		for _, em := range sel.Events {
			methods = append(methods, protoMethodToInternal(em))
		}
		cfg[pt][pi] = methods
	}
	return cfg
}

func (h *SentinelHandler) buildActive(selectors []*sentinelv1.ProtocolSelector) (*comm.ActiveParsers, error) {
	_ = selectors
	// TODO: 接入 parser.Engine
	return &comm.ActiveParsers{}, nil
}

// ─────────────────────────────────────────────
//  Proto enum → 内部类型
// ─────────────────────────────────────────────

func protoTypeToInternal(pt sentinelv1.ProtocolType) comm.ProtocolType {
	switch pt {
	case sentinelv1.ProtocolType_PROTOCOL_TYPE_DEX:
		return comm.ProtocolTypeDEX
	case sentinelv1.ProtocolType_PROTOCOL_TYPE_TOKEN:
		return comm.ProtocolTypeToken
	case sentinelv1.ProtocolType_PROTOCOL_TYPE_LENDING:
		return comm.ProtocolTypeLending
	default:
		return comm.ProtocolTypeUnknown
	}
}

func protoImplToInternal(pi sentinelv1.ProtocolImpl) comm.ProtocolImpl {
	switch pi {
	case sentinelv1.ProtocolImpl_PROTOCOL_IMPL_UNISWAP_V2:
		return comm.ProtocolImplUniswapV2
	case sentinelv1.ProtocolImpl_PROTOCOL_IMPL_SUSHISWAP:
		return comm.ProtocolImplSushiSwap
	case sentinelv1.ProtocolImpl_PROTOCOL_IMPL_ERC20:
		return comm.ProtocolImplERC20
	case sentinelv1.ProtocolImpl_PROTOCOL_IMPL_ERC721:
		return comm.ProtocolImplERC721
	default:
		return comm.ProtocolImplERC20
	}
}

func protoMethodToInternal(em sentinelv1.EventMethod) comm.EventMethod {
	switch em {
	case sentinelv1.EventMethod_EVENT_METHOD_SWAP:
		return comm.EventMethodSwap
	case sentinelv1.EventMethod_EVENT_METHOD_TRANSFER:
		return comm.EventMethodTransfer
	default:
		return comm.EventMethodTransfer
	}
}

// ─────────────────────────────────────────────
//  内部类型 → Proto enum
// ─────────────────────────────────────────────

func internalTypeToProto(pt comm.ProtocolType) sentinelv1.ProtocolType {
	switch pt {
	case comm.ProtocolTypeDEX:
		return sentinelv1.ProtocolType_PROTOCOL_TYPE_DEX
	case comm.ProtocolTypeToken:
		return sentinelv1.ProtocolType_PROTOCOL_TYPE_TOKEN
	case comm.ProtocolTypeLending:
		return sentinelv1.ProtocolType_PROTOCOL_TYPE_LENDING
	default:
		return sentinelv1.ProtocolType_PROTOCOL_TYPE_UNSPECIFIED
	}
}

func internalImplToProto(pi comm.ProtocolImpl) sentinelv1.ProtocolImpl {
	switch pi {
	case comm.ProtocolImplUniswapV2:
		return sentinelv1.ProtocolImpl_PROTOCOL_IMPL_UNISWAP_V2
	case comm.ProtocolImplSushiSwap:
		return sentinelv1.ProtocolImpl_PROTOCOL_IMPL_SUSHISWAP
	case comm.ProtocolImplERC20:
		return sentinelv1.ProtocolImpl_PROTOCOL_IMPL_ERC20
	case comm.ProtocolImplERC721:
		return sentinelv1.ProtocolImpl_PROTOCOL_IMPL_ERC721
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
	case comm.EventMethodDeposit:
		return sentinelv1.EventMethod_EVENT_METHOD_DEPOSIT
	case comm.EventMethodWithdraw:
		return sentinelv1.EventMethod_EVENT_METHOD_WITHDRAW
	case comm.EventMethodFlashLoan:
		return sentinelv1.EventMethod_EVENT_METHOD_FLASHLOAN
	default:
		return sentinelv1.EventMethod_EVENT_METHOD_UNSPECIFIED
	}
}
