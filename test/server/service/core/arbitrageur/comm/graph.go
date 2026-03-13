package comm

import (
	"context"
	"github.com/Crypto-ChainSentinel/types"
	"github.com/machinebox/graphql"
	"log"
)

// -----------------------------------------------------
type GraphQLDEX struct {
	Client      *graphql.Client
	AuthToken   string
	QueryString string
}

func (g GraphQLDEX) GetPairs(vars map[string]interface{}, target interface{}) error {
	req := graphql.NewRequest(g.QueryString)
	for k, v := range vars {
		req.Var(k, v)
	}
	if g.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+g.AuthToken)
	}
	return g.Client.Run(context.Background(), req, target)
}

func (g GraphQLDEX) GetAllPairs() {
	//----每个发现的套利机会 都会立即通过 Flashbots 原子提交----------
	//go func() {
	//	for pair := range arbCh {
	//		// 不直接执行链上交易，而是通过 Flashbots 提交
	//		err := SubmitFlashbotsBundle(pair)
	//		if err != nil {
	//			log.Println("Flashbots submit failed:", err)
	//		}
	//	}
	//}()
	// ----------------- 获取交易对 --------------------------------
	u_V3 := types.Uniswap_V3{
		GraphQLDEX: types.GraphQLDEX{
			Client:    graphql.NewClient("https://gateway.thegraph.com/api/subgraphs/id/5zvR82QoaXYFyDEKLZ9t6v9adgnptxYpKpSbxtgVENFV"),
			AuthToken: "df5d393ba8219b65e3eea66df2242e6b",
			QueryString: `
            query($first: Int!, $skip: Int!) {
                pools(first: $first, skip: $skip) {
                    id
                    token0 { id,symbol }
                    token1 { id,symbol }
                    feeTier
                }
            }
        `,
		},
	}
	var resp_univ3 types.Pairs
	err1 := u_V3.GetPairs(map[string]interface{}{
		"first": 100,
		"skip":  0,
	}, &resp_univ3)
	if err1 != nil {
		log.Fatal(err1)
	}

	var resp_sushi types.Pairs
	s := types.Sushiswap{
		GraphQLDEX: types.GraphQLDEX{
			Client:    graphql.NewClient("https://gateway.thegraph.com/api/subgraphs/id/A4JrrMwrEXsYNAiYw7rWwbHhQZdj6YZg1uVy5wa6g821"),
			AuthToken: "df5d393ba8219b65e3eea66df2242e6b",
			QueryString: `
            query($first: Int!, $skip: Int!) {
                pools(first: $first, skip: $skip) {
                    id
                    token0 { id,symbol }
                    token1 { id,symbol }
                }
            }
        `,
		}}
	err2 := s.GetPairs(map[string]interface{}{
		"first": 100,
		"skip":  0,
	}, &resp_sushi)
	if err2 != nil {
		log.Fatal(err2)
	}

}

// 返回两个 DEX 的交易对交集
func (g GraphQLDEX) GetCommonPairs(pairsA, pairsB types.Pairs) []types.CrossPairData {
	// 构建 DEX A 的交易对 map
	pairMapA := make(map[string]types.Pair)
	for _, p := range pairsA.Pair {
		t0, t1 := normalizePair(p)
		key := t0 + "_" + t1
		pairMapA[key] = p
	}

	// 遍历 DEX B，找交集
	commonPairs := []types.CrossPairData{}
	for _, p := range pairsB.Pair {
		t0, t1 := normalizePair(p)
		key := t0 + "_" + t1
		if pa, ok := pairMapA[key]; ok {
			commonPairs = append(commonPairs, types.CrossPairData{
				Pair_DexA: pa,
				Pair_DexB: p,
			})
		}
	}

	return commonPairs
}
