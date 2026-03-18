// Package wallet 提供统一的钱包接口，整合 HD 钱包、keystore、交易发送。
package wallet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ETHSentinel/internal/wallet/hd"
	"github.com/ETHSentinel/internal/wallet/transaction"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// Manager 钱包管理器，统一对外入口
type Manager struct {
	wallet  *hd.Wallet
	client  *ethclient.Client
	chainID *big.Int
	logger  *zap.Logger
	builder *transaction.Builder
}

// NewManager 从助记词创建钱包管理器
func NewManager(mnemonic string, client *ethclient.Client, chainID *big.Int, logger *zap.Logger) (*Manager, error) {
	w, err := hd.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, fmt.Errorf("wallet: %w", err)
	}
	return &Manager{
		wallet:  w,
		client:  client,
		chainID: chainID,
		logger:  logger,
		builder: transaction.NewBuilder(client),
	}, nil
}

// NewManagerFromSeed 从种子创建钱包管理器
func NewManagerFromSeed(seed []byte, client *ethclient.Client, chainID *big.Int, logger *zap.Logger) (*Manager, error) {
	w, err := hd.NewFromSeed(seed)
	if err != nil {
		return nil, fmt.Errorf("wallet: %w", err)
	}
	return &Manager{
		wallet:  w,
		client:  client,
		chainID: chainID,
		logger:  logger,
		builder: transaction.NewBuilder(client),
	}, nil
}

// ─────────────────────────────────────────────
//  账户管理
// ─────────────────────────────────────────────

// DeriveAccount 派生第 index 个账户
func (m *Manager) DeriveAccount(index uint32) (accounts.Account, error) {
	return m.wallet.DeriveAt(index, true)
}

// Accounts 返回所有已派生账户
func (m *Manager) Accounts() []accounts.Account {
	return m.wallet.Accounts()
}

// Address 返回第 index 个账户地址
func (m *Manager) Address(index uint32) (common.Address, error) {
	acc, err := m.DeriveAccount(index)
	if err != nil {
		return common.Address{}, err
	}
	return acc.Address, nil
}

// ─────────────────────────────────────────────
//  余额查询
// ─────────────────────────────────────────────

// ETHBalance 查询 ETH 余额（wei）
func (m *Manager) ETHBalance(ctx context.Context, addr common.Address) (*big.Int, error) {
	return m.client.BalanceAt(ctx, addr, nil)
}

// ─────────────────────────────────────────────
//  交易
// ─────────────────────────────────────────────

// Send 构造、签名并发送交易
func (m *Manager) Send(ctx context.Context, accountIndex uint32, req *transaction.TxRequest) (*transaction.TxResult, error) {
	acc, err := m.DeriveAccount(accountIndex)
	if err != nil {
		return nil, err
	}
	req.From = acc.Address
	req.ChainID = m.chainID

	tx, err := m.builder.Build(ctx, req)
	if err != nil {
		return nil, err
	}

	signed, err := m.wallet.SignTx(acc, tx, m.chainID)
	if err != nil {
		return nil, err
	}

	result, err := transaction.Send(ctx, m.client, signed)
	if err != nil {
		return nil, err
	}

	m.logger.Info("transaction sent",
		zap.String("hash", result.Hash.Hex()),
		zap.String("from", acc.Address.Hex()),
		zap.String("to", req.To.Hex()),
	)
	return result, nil
}

// SignMessage 签名消息（EIP-191）
func (m *Manager) SignMessage(accountIndex uint32, message []byte) ([]byte, error) {
	acc, err := m.DeriveAccount(accountIndex)
	if err != nil {
		return nil, err
	}
	return m.wallet.SignHash(acc, accounts.TextHash(message))
}

// SignTx 仅签名，不发送
func (m *Manager) SignTx(accountIndex uint32, tx *types.Transaction) (*types.Transaction, error) {
	acc, err := m.DeriveAccount(accountIndex)
	if err != nil {
		return nil, err
	}
	return m.wallet.SignTx(acc, tx, m.chainID)
}

// GenerateMnemonic 生成新助记词（工具函数）
func GenerateMnemonic() (string, error) {
	return hd.GenerateMnemonic(128)
}
