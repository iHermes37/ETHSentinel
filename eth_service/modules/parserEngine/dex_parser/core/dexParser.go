package dexcore

import (
	dexcommon "github.com/CryptoQuantX/chain_monitor/modules/parserEngine/dex_parser/common"
	"github.com/CryptoQuantX/chain_monitor/modules/parserEngine/dex_parser/protocols"
)

type MyEventParser struct {
	//ContractAddresses []common.Address
	EventConfigs map[dexcommon.Protocol]protocols.ProtocolParsers // topic0 -> configs
	//AddressCache      sync.Map                               // map[txHash]string -> []common.Address
	//mu                sync.Mutex
}

func (myparser MyEventParser) NewDexEventParser(needDexs []dexcommon.Protocol, needEvent dexcommon.EventFilter) {
	for _, dex := range needDexs {
		var dexparser = protocols.DEXParseConfigManager[dexcommon.Protocol(dex)]

		filteredConfigs := []dexcommon.EventParseConfig{}
		for _, config := range dexparser.Configs {
			// 如果需要过滤 EventType，则进行过滤
			match := false
			if len(needEvent.FilterEvent) > 0 {
				for _, evt := range needEvent.FilterEvent {
					if evt == config.EventType {
						match = true
					}
					break
				}
			}
			if !match {
				continue // 跳过不匹配的 EventType
			}
			filteredConfigs = append(filteredConfigs, config)

		}

		if len(filteredConfigs) > 0 {
			myparser.EventConfigs[dex] = protocols.ProtocolParsers{
				Addr:    dexparser.Addr,
				Configs: filteredConfigs,
			}
		}

		print(myparser)
	}
}
