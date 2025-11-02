package ParserEngine

import (
	dexparser "github.com/Crypto-ChainSentinel/modules/ParserEngine/DexParser"
	"github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
)

func ParserEngine() *common.ProtocolManager {
	registry := common.GetRegistry()

	// 显式注册所有协议实现（可从配置文件动态加载）
	registry.RegisterProtocol(dexparser.NewUniswapV2Registrar())
	// registry.RegisterProtocol(dexparser.NewSushiSwapRegistrar())

	pm := registry.InitProtocolManager()

	return pm
}

// 使用示例
// func main() {
// 	protocolManager := ParserEngine()
// 	// 选择协议类型
// 	dexParser, _ := protocolManager.GetProtocol(DEX)
// 	// 选择具体实现
// 	uniswapParser, _ := dexParser.GetImplementation(UniswapV2)
// 	// 使用命令模式解析事件
// 	event, err := uniswapParser.HandleEvent(UniswapV2_Swap, types.Log{}, EventMetadata{BlockNumber: 12345})
// 	if err != nil {
// 		fmt.Println("解析失败:", err)
// 	} else {
// 		fmt.Println("事件解析成功:", event)
// 	}
// 	// 切换到 ERC20 协议
// 	erc20Parser, _ := protocolManager.GetProtocol(ERC20)
// 	erc20Impl, _ := erc20Parser.GetImplementation(ERC20Std)
// 	event2, _ := erc20Impl.HandleEvent(ERC20_Transfer, types.Log{}, EventMetadata{TxHash: "0xabc123"})
// 	fmt.Println("ERC20 事件解析:", event2)
// }
