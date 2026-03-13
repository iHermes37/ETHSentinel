package main

import (
	"context"
	"fmt"
	"github.com/Crypto-ChainSentinel/utils"
	"log"

	connectionManager "github.com/Crypto-ChainSentinel/internal/ConnectionManager"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"

	//"github.com/Crypto-ETHSentinel/internal/parser/dex/ERC"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// func cal() {
// 	// 事件签名字符串
// 	eventSig := "Transfer(address,address,uint256)"
// 	// 计算 Keccak256 哈希
// 	hash := crypto.Keccak256Hash([]byte(eventSig))
// 	fmt.Println(hash.Hex()) // 这就是 topics[0]
// }

// func GetEthBlock() *types.Block {
// 	ethClient := connectionManager.InfuraConn()

// 	header, err := ethClient.HeaderByNumber(context.Background(), nil)
// 	if err != nil {
// 		log.Fatalf("获取最新区块失败: %v", err)
// 	}

// 	block, err := ethClient.BlockByNumber(context.Background(), header.Number)
// 	if err != nil {
// 		log.Fatalf("获取区块失败: %v", err)
// 	}

// 	log.Println("获取区块成功")
// 	log.Println(block.Hash)

// 	return block
// }

// 模拟 EventMetadata
// mockMetadata 模拟 EventMetadata
// func mockMetadata() dexcommon.EventMetadata {
// 	txIndex := uint64(0)
// 	return dexcommon.EventMetadata{
// 		EventTypeVal:        dexcommon.UniswapV2_SwapBuy,
// 		ProtocolVal:         dexcommon.UniswapV2,
// 		TxHashVal:           tools.HexToHash("0xdeadbeef"),
// 		BlockNumberVal:      123456,
// 		OuterIndexVal:       0,
// 		TransactionIndexVal: &txIndex,
// 		SwapParsed:          false,
// 	}
// }

// mockLog 构造一个符合 UniswapV2 Swap 事件的 Log
func mockLog(address common.Address) types.Log {
	// Swap 事件的 topic[0] 是 Swap 事件签名 keccak256("Swap(address,uint256,uint256,address)")
	swapSig := common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613d6a0e06d5b69f7c24c7c3f4a45d2")
	// 简单模拟 Data 字段，可以是 ABI 编码后的数量、token 地址等
	data := common.FromHex("00000000000000000000000000000000000000000000000000000000000003e800000000000000000000000000000000000000000000000000000000000001f4")
	return types.Log{
		Address: address,
		Topics:  []common.Hash{swapSig},
		Data:    data,
		// 以下字段可选填
		BlockNumber:    123456,
		TxHash:         common.HexToHash("0xdeadbeef"),
		TxIndex:        0,
		BlockHash:      common.HexToHash("0xabc123"),
		BlockTimestamp: uint64(time.Now().Unix()),
		Index:          0,
		Removed:        false,
	}
}

// func TestUniswapV2SwapParsing() {
// 	// 1. 初始化全局 DexEventParsers
// 	protocols.InitEventConfig()

// 	// 2. 创建 DexEventParser，指定只解析 UniswapV2 和 SwapBuy 事件
// 	parser := &dexcore.MyEventParser{
// 		EventConfigs: make(map[dexcommon.Protocol]protocols.ProtocolParsers),
// 	}

// 	needDexs := []dexcommon.Protocol{dexcommon.UniswapV2}
// 	needEvent := dexcommon.EventFilter{
// 		FilterEvent: []dexcommon.EventType{
// 			dexcommon.UniswapV2_SwapBuy,
// 		},
// 	}

// 	// 3. 初始化解析器
// 	parser.NewDexEventParser(needDexs, needEvent)
// 	client, err := ethclient.Dial("https://mainnet.infura.io/v3/0d79a9c32c814e1da6133850f6fa1128")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 4. 遍历配置并调用解析函数
// 	for dex, configGroup := range parser.EventConfigs {
// 		fmt.Printf("解析 Dex: %v\n", dex)
// 		for _, cfg := range configGroup.Configs {
// 			log := mockLog(cfg.ContractAddress)
// 			metadata := mockMetadata()

// 			// 这里初始化 filterer
// 			filterer, err := abligens.NewUniswappairFilterer(cfg.ContractAddress, client)
// 			if err != nil {
// 				fmt.Printf("filterer 初始化失败: %v\n", err)
// 				continue
// 			}

// 			event, err := cfg.Parser(log, metadata, filterer) // abligens 这里暂时传 nil
// 			if err != nil {
// 				fmt.Printf("解析失败: %+v\n", event)
// 			} else {
// 				fmt.Printf("解析成功: %+v\n", event)
// 			}
// 		}
// 	}
// }

func MonitorMempool() {
	client, err := rpc.Dial("wss://mainnet.infura.io/ws/v3/0d79a9c32c814e1da6133850f6fa1128")
	if err != nil {
		log.Fatal("连接节点失败:", err)
	}
	defer client.Close()

	gc := gethclient.New(client)

	// 用于接收完整的 pending 交易
	//transactions := make(chan *types.Transaction, 100)
	//_, err = gc.SubscribeFullPendingTransactions(context.Background(), transactions)
	//if err != nil {
	//	log.Fatal("订阅 pending 交易失败:", err)
	//}
	////defer sub.Unsubscribe()
	//fmt.Println("✅ 已订阅 pending transactions...")
	//for {
	//	select {
	//	case tx := <-transactions:
	//		from, _ := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	//		fmt.Printf("Pending Tx from: %s, to: %v, value: %s\n",
	//			from.Hex(),
	//			func() string {
	//				if tx.To() != nil {
	//					return tx.To().Hex()
	//				}
	//				return "ContractCreation"
	//			}(),
	//			tx.Value().String())
	//	}
	//}

	//txHashes := make(chan string)
	//_, err = client.Subscribe(context.Background(), "eth", txHashes, "newPendingTransactions")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for txHash := range txHashes {
	//	fmt.Println("Pending tx hash:", txHash)
	//}

	hashes := make(chan common.Hash, 100)
	_, err = gc.SubscribePendingTransactions(context.Background(), hashes)
	if err != nil {
		log.Printf("failed to SubscribePendingTransactions: %v", err)
		return
	}
	log.Print("subscribed pending txs now")
	for {
		select {
		case hash := <-hashes:
			log.Printf("received tx %s", hash)
		}
	}

}

//func watch() {
//	wss := "wss://mainnet.infura.io/ws/v3/0d79a9c32c814e1da6133850f6fa1128"
//
//	rc, err := rpc.Dial(wss)
//	if err != nil {
//		log.Printf("failed to dial: %v", err)
//		return
//	}
//	log.Printf("connected to %s", wss)
//	gc := gethclient.New(rc)
//
//	transactions := make(chan *types.Transaction, 100)
//	_, err = gc.SubscribeFullPendingTransactions(context.Background(), transactions)
//	if err != nil {
//		log.Printf("failed to SubscribePendingTransactions: %v", err)
//		return
//	}
//	log.Print("subscribed pending txs now")
//	for {
//		select {
//		case transaction := <-transactions:
//			// 这里的transaction是完整数据，可以直接使用
//			txBytes, err := transaction.MarshalJSON()
//			if err != nil {
//				continue
//			}
//			log.Printf("received tx %s", string(txBytes))
//		}
//	}
//}

//func main() {
//	//TestUniswapV2SwapParsing()
//	//block := GetEthBlock()
//	//commonParser.ParseBlock(block)
//	//cal()
//	fmt.Println("Hello, World!")
//	//MonitorMempool()
//
//	//go watch()
//	//signalChan := make(chan os.Signal, 1)
//	//signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
//	//<-signalChan
//
//}

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

// https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD

func main() {
	txChan := GetPendingTx() // 启动订阅并获取 channel
	cli := connectionManager.InfuraConn()

	for tx := range txChan {
		fmt.Printf("New pending tx: %s from %s to %v, value: %s\n",
			tx.Hash().Hex(),
			utils.Parsefrom(cli, tx).Hex(),
			func() string {
				if tx.To() != nil {
					return tx.To().Hex()
				}
				return "ContractCreation"
			}(),
			tx.Value().String(),
		)
	}
}
