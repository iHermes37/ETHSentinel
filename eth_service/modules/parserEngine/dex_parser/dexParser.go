package dex_parser

import (
	"encoding/json"
	"fmt"
	dexcommon "github.com/Crypto-ChainSentinel/modules/parserEngine/dex_parser/common"
	"github.com/Crypto-ChainSentinel/modules/parserEngine/dex_parser/protocols"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type MyEventParser struct {
	ContractAddresses []common.Address
	EventConfigs      map[dexcommon.Protocol]protocols.ProtocolParser // topic0 -> configs
	//AddressCache      sync.Map                               // map[txHash]string -> []commonParser.Address
	//mu                sync.Mutex
}

func (myparser MyEventParser) NewDexEventParser(needDexs []dexcommon.Protocol, needEvent dexcommon.EventFilter) {
	for _, dex := range needDexs {
		var dexparser = protocols.DEXParseConfigManager[dex]

		filteredConfigs := map[dexcommon.EventSig]dexcommon.EventParserFunc{}
		for eventsig, eventfunc := range dexparser.Configs {
			match := false
			if len(needEvent.FilterEvent) > 0 {
				for _, sig := range needEvent.FilterEvent {
					if sig == eventsig {
						match = true
					}
					break
				}
			}
			if !match {
				continue // 跳过不匹配的 EventType
			}
			filteredConfigs[eventsig] = eventfunc

		}

		if len(filteredConfigs) > 0 {
			myparser.EventConfigs[dex] = protocols.ProtocolParser{
				ContractAddrs: dexparser.ContractAddrs,
				Configs:       filteredConfigs,
			}
		}

		b, _ := json.MarshalIndent(myparser, "", "  ")
		fmt.Println(string(b))
	}
}

func (myparser MyEventParser) ParseTran(tranreceipt *types.Receipt) {
	metadata := dexcommon.EventMetadata{}
	metadata.BlockNumberVal = tranreceipt.BlockNumber
	metadata.TransactionIndexVal = tranreceipt.TransactionIndex
	metadata.TxHashVal = tranreceipt.TxHash

	for _, log := range tranreceipt.Logs {
		myparser.ParseLog(log, &metadata)
	}

}

func (myparser MyEventParser) ParseLog(log *types.Log, metadata *dexcommon.EventMetadata) dexcommon.UnifiedEvent {
	metadata.OuterIndexVal = log.Index
	myparser.ContractAddresses = append(myparser.ContractAddresses, log.Address)

	// 1. 遍历所有协议
	for _, parser := range myparser.EventConfigs {
		// 2. 检查日志地址是否属于该协议
		if _, ok := parser.ContractAddrs[log.Address]; !ok {
			continue
		}
		// 3. 遍历该协议的事件配置
		for sig, parserFunc := range parser.Configs {
			if common.Hash(sig) == log.Topics[0] {
				//通过 Topic[0] 判断事件类型 ,假设 Parser 内部可以处理 topic 校验
				unifiedEvent, err := parserFunc(*log, *metadata, nil)
				if err != nil {
					continue
				}
				// 这里可以把解析结果存储到 metadata 或其他字段
				//_ = unifiedEvent
				return unifiedEvent
			}
		}

	}

	return nil

}
