package main

import (
	"fmt"

	//"github.com/Crypto-ChainSentinel/modules/ParserEngine/dex_parser/ERC"

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
// 		TxHashVal:           common.HexToHash("0xdeadbeef"),
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

// func findmempool() {
// 	// 使用 WebSocket 连接节点
// 	client, err := rpc.Dial("wss://mainnet.infura.io/ws/v3/YOUR_INFURA_KEY")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	gc := gethclient.New(rc)

// 	transactions := make(chan *types.Transaction, 100)
// 	pendingTransactions, err := gc.SubscribeFullPendingTransactions(context.Background(), transactions)
// 	if err != nil {
// 		return
// 	}

// 	for tx := range transactions {
// 		txBytes, _ := tx.MarshalJSON()
// 		log.Printf("Received tx: %s", string(txBytes))
// 	}

// 	// 创建 channel 接收交易哈希
// 	txHashes := make(chan string)

// 	// 使用 Subscribe 方法订阅 newPendingTransactions
// 	sub, err := client.Subscribe(context.Background(), "eth", txHashes, "newPendingTransactions")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("✅ 已订阅 pending transactions...")

// 	// 循环监听
// 	for {
// 		select {
// 		case err := <-sub.Err():
// 			log.Println("订阅错误:", err)
// 		case txHash := <-txHashes:
// 			fmt.Println("Pending Tx Hash:", txHash)
// 			// 如果需要完整交易对象，可用 ethclient.TransactionByHash 查询
// 		}
// 	}

// }

func main() {
	//TestUniswapV2SwapParsing()
	//block := GetEthBlock()
	//commonParser.ParseBlock(block)
	//cal()
	fmt.Println("Hello, World!")
}
