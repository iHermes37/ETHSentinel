package mempool

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
)

type Monitor interface {
	CollectPendingTx()
	MonitorWhaleRefTx()
}

type Radar struct {
	Client *ethclient.Client
}

func (r *Radar) CollectPendingTx(ctx context.Context) (<-chan *types.Transaction, error) {
	client, err := rpc.Dial("wss://mainnet.infura.io/ws/v3/0d79a9c32c814e1da6133850f6fa1128")
	if err != nil {
		return nil, err
	}
	gc := gethclient.New(client)

	hashes := make(chan common.Hash, 100)
	sub, err := gc.SubscribePendingTransactions(ctx, hashes)
	if err != nil {
		return nil, err
	}

	txCh := make(chan *types.Transaction, 100)

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-ctx.Done():
				return
			case hash := <-hashes:
				tx, isPending, err := r.Client.TransactionByHash(ctx, hash)
				if err != nil || !isPending {
					continue
				}
				select {
				case txCh <- tx:
				case <-ctx.Done():
					return
				}
			case err := <-sub.Err():
				log.Println("subscription error:", err)
			}
		}
	}()

	return txCh, nil
}

// =======================================

func (r *Radar) MonitorWhaleRefTx() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 程序退出时停止监听

	txCh, err := r.CollectPendingTx(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// 循环接收 pending 交易
	for tx := range txCh {
		fmt.Printf("Pending tx: %s\n", tx.Hash().Hex())
		// 可以解析 tx.Data，检查是否是 swap / token transfer 等
	}
}

func (r *Radar) MonitorDexRefTx() {

}

func (r *Radar) MonitorLoanRefTx() {

}
