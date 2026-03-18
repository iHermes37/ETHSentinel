// Package grpc — WalletService gRPC 处理器
package grpc

import (
	"context"
	"fmt"
	"math/big"

	sentinelv1 "github.com/ETHSentinel/gen/sentinel/v1"
	"github.com/ETHSentinel/internal/conn"
	"github.com/ETHSentinel/internal/wallet"
	"github.com/ETHSentinel/internal/wallet/transaction"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WalletHandler 实现 sentinelv1.WalletServiceServer
type WalletHandler struct {
	sentinelv1.UnimplementedWalletServiceServer
	pool   *conn.MultiChainPool
	logger *zap.Logger
}

// NewWalletHandler 创建 Wallet 处理器
func NewWalletHandler(pool *conn.MultiChainPool, logger *zap.Logger) *WalletHandler {
	return &WalletHandler{pool: pool, logger: logger}
}

// CreateWallet 创建或导入钱包
func (h *WalletHandler) CreateWallet(ctx context.Context, req *sentinelv1.CreateWalletRequest) (*sentinelv1.CreateWalletResponse, error) {
	mnemonic := req.Mnemonic

	// 没有传助记词则自动生成
	if mnemonic == "" {
		var err error
		mnemonic, err = wallet.GenerateMnemonic()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "generate mnemonic: %v", err)
		}
	}

	chainID, client, err := h.getChainClient(ctx, req.ChainId)
	if err != nil {
		return nil, err
	}

	mgr, err := wallet.NewManager(mnemonic, client, chainID, h.logger)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "create wallet: %v", err)
	}

	addr, err := mgr.Address(0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "derive address: %v", err)
	}

	return &sentinelv1.CreateWalletResponse{
		Mnemonic: mnemonic,
		Address:  addr.Hex(),
	}, nil
}

// DeriveAccount 派生指定索引的账户地址
func (h *WalletHandler) DeriveAccount(ctx context.Context, req *sentinelv1.DeriveAccountRequest) (*sentinelv1.DeriveAccountResponse, error) {
	if req.Mnemonic == "" {
		return nil, status.Error(codes.InvalidArgument, "mnemonic is required")
	}

	chainID, client, err := h.getChainClient(ctx, "1")
	if err != nil {
		return nil, err
	}

	mgr, err := wallet.NewManager(req.Mnemonic, client, chainID, h.logger)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "create wallet: %v", err)
	}

	acc, err := mgr.DeriveAccount(req.AccountIndex)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "derive account: %v", err)
	}

	return &sentinelv1.DeriveAccountResponse{
		Address: acc.Address.Hex(),
		Path:    acc.URL.Path,
	}, nil
}

// GetBalance 查询账户余额
func (h *WalletHandler) GetBalance(ctx context.Context, req *sentinelv1.BalanceRequest) (*sentinelv1.BalanceResponse, error) {
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address is required")
	}

	chainID, client, err := h.getChainClient(ctx, req.ChainId)
	if err != nil {
		return nil, err
	}

	addr := common.HexToAddress(req.Address)

	// ERC20 代币余额（暂时只支持 ETH 原生余额）
	_ = chainID
	balance, err := client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get balance: %v", err)
	}

	return &sentinelv1.BalanceResponse{
		BalanceWei: balance.String(),
		Symbol:     "ETH",
	}, nil
}

// SendTransaction 构造、签名并发送交易
func (h *WalletHandler) SendTransaction(ctx context.Context, req *sentinelv1.SendTxRequest) (*sentinelv1.SendTxResponse, error) {
	if req.Mnemonic == "" {
		return nil, status.Error(codes.InvalidArgument, "mnemonic is required")
	}
	if req.To == "" {
		return nil, status.Error(codes.InvalidArgument, "to address is required")
	}

	chainID, client, err := h.getChainClient(ctx, req.ChainId)
	if err != nil {
		return nil, err
	}

	mgr, err := wallet.NewManager(req.Mnemonic, client, chainID, h.logger)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "create wallet: %v", err)
	}

	value := new(big.Int)
	if req.ValueWei != "" {
		value.SetString(req.ValueWei, 10)
	}

	txReq := &transaction.TxRequest{
		To:       common.HexToAddress(req.To),
		Value:    value,
		Data:     req.Data,
		GasLimit: req.GasLimit,
	}
	if req.GasPriceWei != "" {
		gp := new(big.Int)
		gp.SetString(req.GasPriceWei, 10)
		txReq.GasPrice = gp
	}

	result, err := mgr.Send(ctx, req.AccountIndex, txReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "send tx: %v", err)
	}

	return &sentinelv1.SendTxResponse{
		TxHash: result.Hash.Hex(),
	}, nil
}

// SignMessage 签名消息
func (h *WalletHandler) SignMessage(_ context.Context, req *sentinelv1.SignMessageRequest) (*sentinelv1.SignMessageResponse, error) {
	if req.Mnemonic == "" {
		return nil, status.Error(codes.InvalidArgument, "mnemonic is required")
	}

	// 签名不需要链上交互，传 nil client
	mgr, err := wallet.NewManager(req.Mnemonic, nil, big.NewInt(1), h.logger)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "create wallet: %v", err)
	}

	sig, err := mgr.SignMessage(req.AccountIndex, req.Message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sign message: %v", err)
	}

	return &sentinelv1.SignMessageResponse{
		Signature: fmt.Sprintf("0x%x", sig),
	}, nil
}

// ─────────────────────────────────────────────
//  内部辅助
// ─────────────────────────────────────────────

func (h *WalletHandler) getChainClient(ctx context.Context, chainIDStr string) (*big.Int, *ethclient.Client, error) {
	if chainIDStr == "" {
		chainIDStr = "1"
	}
	cid, ok := parseChainID(chainIDStr)
	if !ok {
		return nil, nil, status.Errorf(codes.InvalidArgument, "invalid chain_id: %q", chainIDStr)
	}
	client, err := h.pool.GetRPC(ctx, cid)
	if err != nil {
		return nil, nil, status.Errorf(codes.Unavailable, "get rpc client for chain %s: %v", chainIDStr, err)
	}
	return new(big.Int).SetUint64(uint64(cid)), client, nil
}
