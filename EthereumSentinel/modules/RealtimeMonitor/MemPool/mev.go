package RealtimeMonitor

import (
	"context"
	"fmt"
	"log"

	connectionManager "github.com/Crypto-ChainSentinel/modules/ConnectionManager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func GetPendingTx() chan *types.Transaction {
	client, err := rpc.Dial("wss://mainnet.infura.io/ws/v3/0d79a9c32c814e1da6133850f6fa1128")
	if err != nil {
		log.Fatal("连接节点失败:", err)
	}
	defer client.Close()

	gc := gethclient.New(client)

	hashes := make(chan common.Hash, 100)
	_, err = gc.SubscribePendingTransactions(context.Background(), hashes)
	if err != nil {
		log.Printf("failed to SubscribePendingTransactions: %v", err)
	}
	log.Print("subscribed pending txs now")
	txchannel := make(chan *types.Transaction)
	for {
		select {
		case hash := <-hashes:
			log.Printf("received tx %s", hash)
			ethClient := connectionManager.InfuraConn()
			tx, isPending, err := ethClient.TransactionByHash(context.Background(), hash)
			if err != nil {
				log.Println("TransactionByHash error:", err)
				continue
			}
			txchannel <- tx
			fmt.Printf("Tx: %s, Pending: %v\n", tx.Hash().Hex(), isPending)
		}
	}

	return txchannel
}

func MonitorMempool(txpipline chan *types.Transaction) {

	for tx := range txpipline {
		if tx.To() || tx.from() {
			//巨鲸操作---解析---报警

		}
	}
}
