package arbitragur

import (
	"github.com/Crypto-ChainSentinel/db"
	"github.com/Crypto-ChainSentinel/models"
	"strings"
)

type ArbitrageStrategy interface {
	Name() string
	run(stop <-chan struct{})
}

func ValidateAndRecord(p *models.CrossPairData) {
	db.AddToMysql(p)
}

func normalizePair(p models.Pair) (string, string) {
	if strings.ToLower(p.Token0.Symbol) < strings.ToLower(p.Token1.Symbol) {
		return strings.ToLower(p.Token0.Symbol), strings.ToLower(p.Token1.Symbol)
	} else {
		return strings.ToLower(p.Token1.Symbol), strings.ToLower(p.Token0.Symbol)
	}
}

// GetCommonPairs 返回两个 DEX 的交易对交集
func GetCommonPairs(pairsA, pairsB models.Pairs) []models.CrossPairData {
	// 构建 DEX A 的交易对 map
	pairMapA := make(map[string]models.Pair)
	for _, p := range pairsA.Pair {
		t0, t1 := normalizePair(p)
		key := t0 + "_" + t1
		pairMapA[key] = p
	}

	// 遍历 DEX B，找交集
	commonPairs := []models.CrossPairData{}
	for _, p := range pairsB.Pair {
		t0, t1 := normalizePair(p)
		key := t0 + "_" + t1
		if pa, ok := pairMapA[key]; ok {
			commonPairs = append(commonPairs, models.CrossPairData{
				Pair_DexA: pa,
				Pair_DexB: p,
			})
		}
	}

	return commonPairs
}
