package whaler

import (
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
	"github.com/ethereum/go-ethereum/common"
)

// type TrackWhaleConfig struct {
// 	StartBlock    *big.Int
// 	EndBlock      *big.Int
// 	TargetAddress common.Address
// 	IsAllWhale    bool
// 	selected      *map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl
// }

// func TrackWhale(cfg TrackWhaleConfig) {
// 	scancfg := scancommon.ScanBlocksConfig{
// 		FilterCfg: FilterConfig{
// 			Filter: Findwhale,
// 			TrackCfg: &TrackWhaleConfig{
// 				cfg.TargetAddress,
// 				cfg.IsAllWhale,
// 			},
// 		},
// 		StartBlock: cfg.StartBlock,
// 		EndBlock:   cfg.EndBlock,
// 		Selected:   cfg.selected,
// 	}

// 	var evpipline chan [][]ParserEngineCommon.UnifiedEvent
// 	evpipline = scancommon.ScanBlocks(scancfg)

// 	for evlists := range evpipline {
// 		fmt.Println("Received event batch:", evlists)
// 		for _, evlist := range evlists {
// 			//处理巨鲸交互
// 			HandleWhaleEvents(evlist)
// 		}
// 	}
// }

// func HandleWhaleEvents(evlist []ParserEngineCommon.UnifiedEvent) {
// 	for _, ev := range evlist {
// 		switch ev.ProtocolType() {
// 		//ERC代币交易
// 		case ERC:
// 			break

// 		// 其他非标准 ERC 合约交互
// 		case Unknow:
// 			break

// 		}
// 	}
// }

// ==========================================

type HarpoonSettings struct {
	TrackETHTransfer     bool
	TrackTokenTransfer   bool
	TrackDefiInteraction bool
	IsAllWhale           bool
	DefiSettings         *map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl
}

type Harpoon struct {
	TargetAddress common.Address
}

type Track interface {
	TrackETHTransfer()
	TrackTokenTransfer()
	TrackDefiInteraction()
}
