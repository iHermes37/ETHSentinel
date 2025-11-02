package capturedWhale

// func HandleNewWhale(evlists chan [][]ParserEngineCommon.UnifiedEventData, Threshold common.Decimal) {
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
// 	Threshold common.Decimal

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

// 			whaleAddr, addrOk := row["holder"].(common.Address)
// 			whaleAmount, amountOk := row["balance"].(common.Decimal)

// 			if !addrOk || !amountOk {
// 				log.Printf("Skipping row %d: invalid types - addr: %T, amount: %T",
// 					i, row["holder"], row["balance"])
// 				continue
// 			}

// 			NewWhale := models.Whale{
// 				Address: whaleAddr,
// 				Amount:  whaleAmount,
// 			}

// 			StoreWhale(NewWhale)

// 		}
// 	}

// 	if cfg.ScanWhale {
// 		scancfg := Scanner.ScanBlocksConfig{
// 			FilterCfg:  cfg.fcfg,
// 			StartBlock: cfg.StartBlock,
// 			EndBlock:   cfg.EndBlock,
// 			Selected:   cfg.selected,
// 		}

// 		whalePipline := Scanner.ScanBlocks(scancfg)
// 		HandleNewWhale(whalePipline, cfg.Threshold)
// 	}

// }
