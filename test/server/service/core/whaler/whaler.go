package whaler

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// func HandleNewWhale(evlists chan [][]ParserEngineCommon.UnifiedEventData, Threshold tools.Decimal) {
// 	for evlist := range evlists { // 第一层：事件批次
// 		for _, events := range evlist { // 第二层：事件切片
// 			for _, ev := range events { // 第三层：单个事件
// 				processWhaleEvent(ev, Threshold)
// 			}
// 		}
// 	}
// }

// type FindwhaleConfig struct {
// 	fcfg      Filter.FilterConfig
// 	selected  *map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl
// 	Threshold tools.Decimal

// 	TokenHold  bool
// 	TokenName  *string
// 	ScanWhale  bool
// 	StartBlock *big.Int
// 	EndBlock   *big.Int
// }

// func Findwhale(cfg FindwhaleConfig) {

// 	if cfg.TokenHold {
// 		analyzerType, _ := tokenHold.GetAnalyzer(*cfg.TokenName)
// 		ctx := tokenHold.AnalyzerContext{Strategy: analyzerType}
// 		res := ctx.Analyze()

// 		for i, row := range res {

// 			whaleAddr, addrOk := row["holder"].(tools.Address)
// 			whaleAmount, amountOk := row["balance"].(tools.Decimal)

// 			if !addrOk || !amountOk {
// 				log.Printf("Skipping row %d: invalid types - addr: %T, amount: %T",
// 					i, row["holder"], row["balance"])
// 				continue
// 			}

// 			NewWhale := types.Whale{
// 				Address: whaleAddr,
// 				Amount:  whaleAmount,
// 			}

// 			StoreWhale(NewWhale)

// 		}
// 	}

// 	if cfg.ScanWhale {
// 		scancfg := scanner.ScanBlocksConfig{
// 			FilterCfg:  cfg.fcfg,
// 			StartBlock: cfg.StartBlock,
// 			EndBlock:   cfg.EndBlock,
// 			Selected:   cfg.selected,
// 		}

// 		whalePipline := scanner.ScanBlocks(scancfg)
// 		HandleNewWhale(whalePipline, cfg.Threshold)
// 	}

// }

type Interval struct {
	StartBlock *big.Int
	EndBlock   *big.Int
}

type Whaler struct {
	Interval
	Sonar *Sonar
}

func (whaler *Whaler) ProcessEthTransaction(tx *types.Transaction) {

}

func (whaler *Whaler) ProcessTokenTransaction(rec *types.Receipt) {
}

func (whaler *Whaler) ProcessDefiTransaction(rec *types.Receipt) {
}

func (whaler *Whaler) ProcessNewContract(rec *types.Receipt) {

}
