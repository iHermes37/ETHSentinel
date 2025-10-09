package common

import (
	"fmt"

	"github.com/Crypto-ChainSentinel/modules/ParserEngine"
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func ParseTranByLog(tranreceipt *types.Receipt,
	selectedProtocols map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl) ParserEngineCommon.UnifiedEvent {

	metadata := ParserEngineCommon.EventMetadata{}
	metadata.BlockNumberVal = tranreceipt.BlockNumber
	metadata.TransactionIndexVal = tranreceipt.TransactionIndex
	metadata.TxHashVal = tranreceipt.TxHash

	protocolManager := ParserEngine.ParserEngine()
	chain := BuildParserChain(protocolManager, selectedProtocols)

	var ev ParserEngineCommon.UnifiedEvent

	for _, log := range tranreceipt.Logs {
		if ev, ok := chain.Handle(*log, metadata); ok {
			fmt.Println("解析成功:", ev)
		} else {
			fmt.Println("未匹配事件", log.Topics[0])
		}
	}

	return ev
}
