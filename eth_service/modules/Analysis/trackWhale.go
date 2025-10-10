package Analysis

import (
	"fmt"
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
	scancommon "github.com/Crypto-ChainSentinel/modules/Scanner/Common"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type TrackWhaleConfig struct {
	WhaleAddr  *common.Address
	StartBlock *big.Int
	EndBlock   *big.Int
	selected   *map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl
}

func TrackWhale(cfg TrackWhaleConfig) {
	scancfg := scancommon.ScanBlocksConfig{
		WhaleAddr:  cfg.WhaleAddr,
		StartBlock: cfg.StartBlock,
		EndBlock:   cfg.EndBlock,
	}

	var evpipline chan [][]ParserEngineCommon.UnifiedEvent

	evpipline = scancommon.ScanBlocks(scancfg)

	for evlists := range evpipline {
		fmt.Println("Received event batch:", evlists)
		for _, evlist := range evlists {
			//处理巨鲸交互
			HandleWhaleEvents(evlist)
		}

	}

}

func HandleWhaleEvents(evlist []ParserEngineCommon.UnifiedEvent) {

}
