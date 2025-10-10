package tokenHold

import (
	"fmt"
	"log"
)

// 2️⃣ 不同代币实现不同策略
type USDCAnalyzer struct {
	apikey  string
	qureyId string
}

func (a *USDCAnalyzer) FetchData() ([]map[string]interface{}, error) {
	fmt.Println("执行 USDC 的 Dune SQL 脚本")
	rows, _ := FetchData(a.apikey, a.qureyId)

	return rows, nil
}

type USDTAnalyzer struct {
	apikey  string
	qureyId string
}

func (a *USDTAnalyzer) FetchData() ([]map[string]interface{}, error) {
	fmt.Println("执行 USDT 的Dune SQL 脚本")
	rows, _ := FetchData(a.apikey, a.qureyId)
	return rows, nil
}

// 3️⃣ 上下文：持有策略接口
type AnalyzerContext struct {
	Strategy HoldingAnalyzer
}

func (ctx *AnalyzerContext) Analyze() []map[string]interface{} {
	res, err := ctx.Strategy.FetchData()
	if err != nil {
		log.Fatal("")
	}
	// ctx.Strategy.ParseData()
	// res, _ := ctx.Strategy.ComputeHoldings()
	// fmt.Println("分析结果：", res)
	return res
}

// func main() {
// 	analyzerType, _ := analyzer.GetAnalyzer("USDC")
// 	ctx := analyzer.AnalyzerContext{Strategy: analyzerType}
// 	ctx.Analyze()
// }
