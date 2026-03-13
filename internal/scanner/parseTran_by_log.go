package scanner

import (
	"fmt"

	"github.com/Crypto-ChainSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/core/types"
)

func (s *Scanner) ParseTranByLog(
	tranreceipt *types.Receipt,
	selected map[comm.ProtocolTypeName][]comm.ProtocolImplName) []comm.UnifiedEvent {

	metadata := comm.EventMetadata{}
	metadata.BlockNumberVal = tranreceipt.BlockNumber
	metadata.TransactionIndexVal = tranreceipt.TransactionIndex
	metadata.TxHashVal = tranreceipt.TxHash
	//metadata.ProtocolType=
	//metadata.ProtocolImpl=

	var protocolManager *comm.ProtocolManager
	//protocolManager = parser.Parser.Mgr

	chain := BuildParserChain(protocolManager, selected)

	var tran_evlist []comm.UnifiedEvent

	for _, log := range tranreceipt.Logs {
		if ev, ok := chain.Handle(*log, metadata); ok {
			fmt.Println("解析成功:", ev)
			tran_evlist = append(tran_evlist, ev)
		} else {
			fmt.Println("未匹配事件", log.Topics[0])
		}
	}

	return tran_evlist
}
