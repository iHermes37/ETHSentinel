// Package grpc — MempoolService gRPC 处理器
package grpc

import (
	"math/big"

	sentinelv1 "github.com/ETHSentinel/gen/sentinel/v1"
	"github.com/ETHSentinel/internal/chain"
	"github.com/ETHSentinel/internal/conn"
	"github.com/ETHSentinel/internal/mempool"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MempoolHandler 实现 sentinelv1.MempoolServiceServer
type MempoolHandler struct {
	sentinelv1.UnimplementedMempoolServiceServer
	pool   *conn.MultiChainPool
	logger *zap.Logger
}

// NewMempoolHandler 创建 Mempool 处理器
func NewMempoolHandler(pool *conn.MultiChainPool, logger *zap.Logger) *MempoolHandler {
	return &MempoolHandler{pool: pool, logger: logger}
}

// SubscribePending 实时订阅 pending 交易（服务端流式）
func (h *MempoolHandler) SubscribePending(
	req *sentinelv1.MempoolSubscribeRequest,
	stream sentinelv1.MempoolService_SubscribePendingServer,
) error {
	ctx := stream.Context()

	chainID, ok := parseChainID(req.ChainId)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "invalid chain_id: %q", req.ChainId)
	}

	c, err := chain.Get(chainID)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "unsupported chain: %v", err)
	}

	wsClient, err := h.pool.GetWS(ctx, chainID)
	if err != nil {
		return status.Errorf(codes.Unavailable, "get ws client: %v", err)
	}

	mon := mempool.NewMonitor(c.DefaultWSURL(), wsClient, uint64(chainID), h.logger)

	// 应用过滤器
	if req.MinValueWei != "" {
		minVal := new(big.Int)
		if _, ok := minVal.SetString(req.MinValueWei, 10); ok {
			mon.WithFilter(mempool.FilterByMinValue(minVal))
		}
	}
	if req.MinGasGwei > 0 {
		mon.WithFilter(mempool.FilterByMinGasPrice(req.MinGasGwei))
	}
	if len(req.MethodSigs) > 0 {
		mon.WithFilter(mempool.FilterByMethodSig(req.MethodSigs...))
	}

	pendingCh, err := mon.Subscribe(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "subscribe mempool: %v", err)
	}

	for ptx := range pendingCh {
		to := ""
		if ptx.Tx.To() != nil {
			to = ptx.Tx.To().Hex()
		}
		gasPrice := ""
		if ptx.GasPrice != nil {
			gasPrice = ptx.GasPrice.String()
		}
		gasTip := ""
		if ptx.GasTipCap != nil {
			gasTip = ptx.GasTipCap.String()
		}

		if err := stream.Send(&sentinelv1.PendingTxEvent{
			TxHash:   ptx.Tx.Hash().Hex(),
			From:     ptx.From.Hex(),
			To:       to,
			Value:    ptx.Tx.Value().String(),
			GasPrice: gasPrice,
			GasTip:   gasTip,
			GasLimit: ptx.Tx.Gas(),
			Input:    ptx.Tx.Data(),
			ChainId:  req.ChainId,
			SeenAt:   timestamppb.New(ptx.SeenAt),
		}); err != nil {
			return err
		}
	}
	return nil
}

func parseChainID(s string) (chain.ChainID, bool) {
	n := new(big.Int)
	if _, ok := n.SetString(s, 10); !ok {
		return 0, false
	}
	return chain.ChainID(n.Uint64()), true
}
