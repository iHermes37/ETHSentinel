package main

import (
	"fmt"
	dexcommon "github.com/CryptoQuantX/chain_monitor/modules/parserEngine/dex_parser/common"
	dexcore "github.com/CryptoQuantX/chain_monitor/modules/parserEngine/dex_parser/core"
	"github.com/CryptoQuantX/chain_monitor/modules/parserEngine/dex_parser/protocols"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// 模拟 EventMetadata
// mockMetadata 模拟 EventMetadata
func mockMetadata() dexcommon.EventMetadata {
	txIndex := uint64(0)
	return dexcommon.EventMetadata{
		EventTypeVal:        dexcommon.UniswapV2_SwapBuy,
		ProtocolVal:         dexcommon.UniswapV2,
		TxHashVal:           common.HexToHash("0xdeadbeef"),
		BlockNumberVal:      123456,
		OuterIndexVal:       0,
		TransactionIndexVal: &txIndex,
		SwapParsed:          false,
	}
}

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

func TestUniswapV2SwapParsing() {
	// 1. 初始化全局 DexEventParsers
	protocols.InitEventConfig()

	// 2. 创建 DexEventParser，指定只解析 UniswapV2 和 SwapBuy 事件
	parser := &dexcore.MyEventParser{
		EventConfigs: make(map[dexcommon.Protocol]protocols.ProtocolParsers),
	}

	needDexs := []dexcommon.Protocol{dexcommon.UniswapV2}
	needEvent := dexcommon.EventFilter{
		FilterEvent: []dexcommon.EventType{
			dexcommon.UniswapV2_SwapBuy,
		},
	}

	// 3. 初始化解析器
	parser.NewDexEventParser(needDexs, needEvent)

	// 4. 遍历配置并调用解析函数
	for dex, configGroup := range parser.EventConfigs {
		fmt.Printf("解析 Dex: %v\n", dex)
		for _, cfg := range configGroup.Configs {
			log := mockLog(cfg.ContractAddress)
			metadata := mockMetadata()

			event, err := cfg.Parser(log, metadata, nil) // abligens 这里暂时传 nil
			if err != nil {
				fmt.Printf("解析失败: %+v\n", event)
			} else {
				fmt.Printf("解析成功: %+v\n", event)
			}
		}
	}
}

func main() {
	TestUniswapV2SwapParsing()
}
