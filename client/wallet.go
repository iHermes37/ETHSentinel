package sentinel

import (
	"context"
	"math/big"

	sentinelv1 "github.com/ETHSentinel/gen/sentinel/v1"
	"github.com/ETHSentinel/internal/wallet"
	"github.com/ETHSentinel/internal/wallet/transaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// WalletClient 钱包客户端
type WalletClient struct {
	// 嵌入模式
	mgr *wallet.Manager
	// 远程模式
	walletStub sentinelv1.WalletServiceClient
	// 公共
	mnemonic string
	client   *ethclient.Client
	chainID  *big.Int
	logger   *zap.Logger
}

// TxRequest 交易请求
type TxRequest = transaction.TxRequest

// TxResult 交易结果
type TxResult = transaction.TxResult

func newWalletClient(mnemonic string, client *ethclient.Client, chainID *big.Int, logger *zap.Logger) (*WalletClient, error) {
	mgr, err := wallet.NewManager(mnemonic, client, chainID, logger)
	if err != nil {
		return nil, err
	}
	return &WalletClient{
		mgr:      mgr,
		mnemonic: mnemonic,
		client:   client,
		chainID:  chainID,
		logger:   logger,
	}, nil
}

// Address 获取第 index 个账户地址
func (w *WalletClient) Address(index uint32) (common.Address, error) {
	if w.walletStub != nil {
		resp, err := w.walletStub.DeriveAccount(context.Background(), &sentinelv1.DeriveAccountRequest{
			Mnemonic: w.mnemonic, AccountIndex: index,
		})
		if err != nil {
			return common.Address{}, err
		}
		return common.HexToAddress(resp.Address), nil
	}
	return w.mgr.Address(index)
}

// ETHBalance 查询 ETH 余额（wei）
func (w *WalletClient) ETHBalance(ctx context.Context, addr common.Address) (*big.Int, error) {
	if w.walletStub != nil {
		resp, err := w.walletStub.GetBalance(ctx, &sentinelv1.BalanceRequest{
			Address: addr.Hex(),
			ChainId: w.chainID.String(),
		})
		if err != nil {
			return nil, err
		}
		bal := new(big.Int)
		bal.SetString(resp.BalanceWei, 10)
		return bal, nil
	}
	return w.mgr.ETHBalance(ctx, addr)
}

// Send 发送交易
func (w *WalletClient) Send(ctx context.Context, accountIndex uint32, to common.Address, valueWei *big.Int, data []byte) (*TxResult, error) {
	if w.walletStub != nil {
		resp, err := w.walletStub.SendTransaction(ctx, &sentinelv1.SendTxRequest{
			Mnemonic:     w.mnemonic,
			AccountIndex: accountIndex,
			To:           to.Hex(),
			ValueWei:     valueWei.String(),
			Data:         data,
			ChainId:      w.chainID.String(),
		})
		if err != nil {
			return nil, err
		}
		return &TxResult{Hash: common.HexToHash(resp.TxHash)}, nil
	}
	return w.mgr.Send(ctx, accountIndex, &TxRequest{
		To: to, Value: valueWei, Data: data, ChainID: w.chainID,
	})
}

// SignMessage 签名消息（EIP-191）
func (w *WalletClient) SignMessage(accountIndex uint32, message []byte) ([]byte, error) {
	if w.walletStub != nil {
		resp, err := w.walletStub.SignMessage(context.Background(), &sentinelv1.SignMessageRequest{
			Mnemonic: w.mnemonic, AccountIndex: accountIndex, Message: message,
		})
		if err != nil {
			return nil, err
		}
		return common.FromHex(resp.Signature), nil
	}
	return w.mgr.SignMessage(accountIndex, message)
}

// Mnemonic 返回助记词
func (w *WalletClient) Mnemonic() string { return w.mnemonic }

// GenerateMnemonic 生成新助记词（12个单词）
func GenerateMnemonic() (string, error) {
	return wallet.GenerateMnemonic()
}
