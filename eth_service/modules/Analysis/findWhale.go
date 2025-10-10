package Analysis

import (
	"math/big"

	"github.com/Crypto-ChainSentinel/models"
	"github.com/Crypto-ChainSentinel/modules/RealtimeMonitor/tokenHold"
)

type FindwhaleConfig struct {
	TokenHold bool
	TokenName *string
	ScanWhale bool
	StartBlock big.Int	
	EndBlock big.Int
}

func Findwhale(cfg FindwhaleConfig) {

	if cfg.TokenHold {
		analyzerType, _ := tokenHold.GetAnalyzer(*cfg.TokenName)
		ctx := tokenHold.AnalyzerContext{Strategy: analyzerType}
		res := ctx.Analyze()

		for _, row := range res {
			whaleAddr:=row["holder"]
			whaleAmount:=row["balance"]

			whaleTable:=models.Whale{
				Address:whaleAddr,
				// Amount:whaleAmount
			}
			
		}
	}

	if cfg.ScanWhale{
		scancfg:=ScanBlocksConfig{
			cfg.StartBlock
			cfg.EndBlock
		}

		whalePipline:=ScanBlocks(scancfg)
		HandleNewWhale(whalePipline)

	}
}


func HandleNewWhale(evlists chan [][]UnifiedEvent) {
	for evlist := range evlists {
		for _, ev := range evlist {
			switch ev.EventType {
				
			case "Transfer", "TransferFrom":
				if ev.Amount > Threshold {
					// 记录巨鲸地址、金额
				}
			case "Swap":
				if ev.Amount0In > Threshold || ev.Amount1In > Threshold ||
					ev.Amount0Out > Threshold || ev.Amount1Out > Threshold {
					// 记录巨鲸 Swap 行为
				}
			case "Mint", "Burn":
				if ev.Amount0 > Threshold || ev.Amount1 > Threshold {
					// 记录巨鲸流动性操作
				}
			case "SafeTransferFrom", "SetApprovalForAll":
				// ERC721 巨鲸操作
			}
		}
	}
}
