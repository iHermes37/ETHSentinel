package whaler

import (
	"github.com/Crypto-ChainSentinel/db"
	"github.com/Crypto-ChainSentinel/modules/parse_engine/comm"
	"github.com/Crypto-ChainSentinel/modules/scanner"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// type TrackWhaleConfig struct {
// 	StartBlock    *big.Int
// 	EndBlock      *big.Int
// 	TargetAddress tools.Address
// 	IsAllWhale    bool
// 	selected      *map[comm.ProtocolType][]comm.ProtocolImpl
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

// 	var evpipline chan [][]comm.UnifiedEvent
// 	evpipline = scancommon.ScanBlocks(scancfg)

// 	for evlists := range evpipline {
// 		fmt.Println("Received event batch:", evlists)
// 		for _, evlist := range evlists {
// 			//处理巨鲸交互
// 			HandleWhaleEvents(evlist)
// 		}
// 	}
// }

// func HandleWhaleEvents(evlist []comm.UnifiedEvent) {
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
	ParseEngineSettings  *map[comm.ProtocolType][]comm.ProtocolImpl
}

type Harpoon struct {
	*HarpoonSettings
	TargetAddress common.Address
	MysqlMgr      *db.MysqlMgr
}

type Track interface {
	TrackETHTransfer()
	TrackTokenTransfer()
	TrackDefiInteraction()
}

//==================================================

func (hp *Harpoon) TrackETHTransfer() {

}

func (hp *Harpoon) TrackTokenTransfer() {

}

func (hp *Harpoon) TrackDefiInteraction(receipt *types.Receipt) {
	tranEvlist := scanner.ParseTranByLog(receipt, *hp.HarpoonSettings.ParseEngineSettings)
	for _, tranEv := range tranEvlist {
		base := tranEv.CoreEvent()
		addr := base.From
		if hp.MysqlMgr.IsWhaleInPool(&addr) {
			//生成全局概述图

			//将tranEv转化成表结构进行日志存储
		}
	}
}
