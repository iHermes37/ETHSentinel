// Package transaction 提供交易构造、发送、追踪能力。
package transaction

import (
	"context"
	"fmt"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TxRequest 交易请求参数
type TxRequest struct {
	From     common.Address
	To       common.Address
	Value    *big.Int // wei
	Data     []byte
	GasLimit uint64   // 0 = 自动估算
	GasPrice *big.Int // nil = 自动获取
	Nonce    *uint64  // nil = 自动获取
	ChainID  *big.Int
}

// TxResult 交易结果
type TxResult struct {
	Hash    common.Hash
	Receipt *types.Receipt
}

// Builder 交易构造器
type Builder struct {
	client *ethclient.Client
}

// NewBuilder 创建交易构造器
func NewBuilder(client *ethclient.Client) *Builder {
	return &Builder{client: client}
}

// Build 根据 TxRequest 构造 types.Transaction
func (b *Builder) Build(ctx context.Context, req *TxRequest) (*types.Transaction, error) {
	if b.client == nil {
		return nil, fmt.Errorf("tx: ethclient is nil")
	}

	// 获取 nonce
	nonce := uint64(0)
	if req.Nonce != nil {
		nonce = *req.Nonce
	} else {
		n, err := b.client.PendingNonceAt(ctx, req.From)
		if err != nil {
			return nil, fmt.Errorf("tx: get nonce: %w", err)
		}
		nonce = n
	}

	// 获取 gasPrice
	gasPrice := req.GasPrice
	if gasPrice == nil {
		gp, err := b.client.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("tx: suggest gas price: %w", err)
		}
		gasPrice = gp
	}

	// 估算 gasLimit
	gasLimit := req.GasLimit
	if gasLimit == 0 {
		if len(req.Data) == 0 {
			gasLimit = 21000 // 普通 ETH 转账
		} else {
			callMsg := ethereum.CallMsg{
				From:     req.From,
				To:       &req.To,
				Value:    req.Value,
				Data:     req.Data,
				GasPrice: gasPrice,
			}
			estimated, err := b.client.EstimateGas(ctx, callMsg)
			if err != nil {
				gasLimit = 200000 // 估算失败时的默认值
			} else {
				gasLimit = estimated + 10000 // 留 buffer
			}
		}
	}

	value := req.Value
	if value == nil {
		value = big.NewInt(0)
	}

	tx := types.NewTransaction(nonce, req.To, value, gasLimit, gasPrice, req.Data)
	return tx, nil
}

// Send 发送已签名的交易，等待上链并返回 receipt
func Send(ctx context.Context, client *ethclient.Client, signedTx *types.Transaction) (*TxResult, error) {
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return nil, fmt.Errorf("tx: send: %w", err)
	}

	hash := signedTx.Hash()
	deadline := time.Now().Add(3 * time.Minute)

	for time.Now().Before(deadline) {
		receipt, err := client.TransactionReceipt(ctx, hash)
		if err == nil {
			return &TxResult{Hash: hash, Receipt: receipt}, nil
		}
		select {
		case <-ctx.Done():
			return &TxResult{Hash: hash}, nil
		case <-time.After(2 * time.Second):
		}
	}

	return &TxResult{Hash: hash}, nil
}
