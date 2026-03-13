package common

import (
	"log"
	"math/big"

	ParserEngineCommon "github.com/Crypto-ChainSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/common"
)

type Analysis interface {
	ThresholdAlarm(TransEvents []ParserEngineCommon.UnifiedEventData, Threshold common.Decimal)
}

type Analyst struct{}

func (A Analyst) ThresholdAlarm(TransEvents []ParserEngineCommon.UnifiedEventData, Threshold common.Decimal) {
	for _, ev := range TransEvents {
		processWhaleEvent(ev, Threshold)
	}
}

func processWhaleEvent(ev ParserEngineCommon.UnifiedEventData, Threshold common.Decimal) {
	switch ev.EventTypeVal {
	case "Transfer", "TransferFrom":
		// 检查 AmountVal 列表中是否有超过阈值的值
		if hasAmountAboveThreshold(ev.AmountVal, Threshold) {
			log.Printf("Whale transfer: %s, Amounts: %v", ev.From, ev.AmountVal)
		}
	case "Swap":
		// 检查 Swap 相关的金额列表
		if isWhaleSwap(ev) {
			log.Printf("Whale swap detected: %s", ev.From)
		}
		// ... 其他事件类型
	}
}

// 检查金额列表中是否有超过阈值的值
func hasAmountAboveThreshold(amounts []*big.Int, Threshold common.Decimal) bool {
	for _, amount := range amounts {
		if amount != nil && amount.Cmp(Threshold) > 0 {
			return true
		}
	}
	return false
}

// 获取超过阈值的金额列表
func getAmountsAboveThreshold(amounts []*big.Int, Threshold common.Decimal) []*big.Int {
	var result []*big.Int
	for _, amount := range amounts {
		if amount != nil && amount.Cmp(Threshold) > 0 {
			result = append(result, amount)
		}
	}
	return result
}
