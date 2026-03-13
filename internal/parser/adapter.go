package parser

import (
	"fmt"
	"github.com/Crypto-ChainSentinel/internal/parser/comm"
	"github.com/Crypto-ChainSentinel/internal/parser/dex"
)

type ParserEngine interface {
	RegisterAllParser()
	SetParser(impl comm.ProtocolImplName, event comm.EventMethod)
}

type Parser struct {
	Mgr comm.ProtocolManager
}

func (p *Parser) RegisterAllParser() {
	pm := comm.NewProtocolManager()
	pm.RegisterProtocol(comm.DEX, dex.DexAdapter())
}

func (p *Parser) SetParser(cfg *comm.ParserCfg) {
	for protocolTypeName, implCfg := range *cfg {
		typeParsers := p.Mgr.AllProtocols[protocolTypeName]
		for implName, selectedEvents := range implCfg {

			implParser, err := typeParsers.GetImpl(implName)
			if err != nil {
				fmt.Println("xxx")
			}

			implParser.SetEvents(selectedEvents)

			p.Mgr.CurProtocols[implName] = implParser
		}
	}
}

//func main() *comm.ProtocolManager {
//	//registry := comm.GetRegistry()
//
//	// 显式注册所有协议实现（可从配置文件动态加载）
//	//registry.RegisterProtocol(dexparser.Register())
//	// registry.RegisterProtocol(dexparser.NewSushiSwapRegistrar())
//	//pm := registry.InitProtocolManager()
//
//	return pm
//}

// 使用示例
// func main() {
// 	protocolManager := parser()
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
