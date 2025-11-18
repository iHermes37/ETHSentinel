package RealtimeMonitor

import (
	"context"
	"github.com/machinebox/graphql"
)


// 对应 GraphQL 返回的每个池
type Pool struct {
	ID      string `json:"id"`
	Token0  Token  `json:"token0"`
	Token1  Token  `json:"token1"`
	FeeTier int    `json:"feeTier,string"`    // GraphQL 返回的是 string，需要用 string 标签解析为 int
	Liquidity string `json:"liquidity"`         // 数量较大，用 string 保存
	VolumeUSD string `json:"volumeUSD"`         // 流动性金额，用 string 保存
	CreatedAt int64  `json:"createdAtTimestamp,string"` // GraphQL 返回 string，需要转换为 int64
}

// 对应 token 对象
type Token struct {
	ID       string `json:"id"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals,string"` // GraphQL 返回 string
}

uniswapv3_query :=`
{
	pools(first: 100, orderBy: createdAtTimestamp, orderDirection: asc) {
		id
		token0 { 
		id 
		symbol 
		decimals 
		}
		token1 { 
		id 
		symbol 
		decimals 
		}
		feeTier
		liquidity
		volumeUSD
		createdAtTimestamp
	}
}
`
uniswapv2_query:=`
{
	pairs(first: 5, orderBy: createdAtTimestamp, orderDirection: desc) {
		id             # 池地址
		token0 {
		id
		symbol
		decimals
		}
		token1 {
		id
		symbol
		decimals
		}
		reserveUSD
		volumeUSD
		createdAtTimestamp
	}
}
`

sushiswap_query:=`
{
	pairs(first: 5, orderBy: createdAtTimestamp, orderDirection: desc) {
		id              # 交易对地址
		token0 {
		id
		symbol
		decimals
		}
		token1 {
		id
		symbol
		decimals
		}
		reserve0        # token0 当前储备
		reserve1        # token1 当前储备
		reserveUSD      # 总流动性（美元计）
		totalSupply     # LP 代币总量
		volumeUSD       # 交易量
		createdAtTimestamp       # 创建时间戳
	}
}
`

// curve_query:=`
// {
// 	liquidityPools(first: 5, orderBy: createdTimestamp, orderDirection: desc) {
// 		id
// 		inputTokens {
// 		id
// 		symbol
// 		decimals
// 		}
// 		outputToken {
// 		id
// 		symbol
// 		decimals
// 		}

// 	}
// }
// `
// 泛型函数
func FetchGraphQL[T any](client *graphql.Client, query string) (*T, error) {
	req := graphql.NewRequest(query)
	var respData T
	if err := client.Run(context.Background(), req, &respData); err != nil {
		return nil, err
	}
	return &respData, nil
}



type Dex string

const (
	UniswapV2 Dex = "UniswapV2"
	UniswapV3 Dex = "UniswapV3"
	SushiSwap Dex = "SushiSwap"
	// Curve     Dex = "Curve"
)

type UnifiedToken struct {
	ID       string
	Symbol   string
	Decimals int
}

type UnifiedPool struct {
	PoolID       string
	DEX          Dex
	Tokens       []UnifiedToken
	LiquidityUSD float64
	VolumeUSD    float64
	CreatedAt    int64
	FeeTier      int    // 可选，V3 有 feeTier
	LPTokenSupply float64 // 可选，SushiSwap 有
}

func MapUniswapV2PairToUnified(pair Pair) UnifiedPool {
	return UnifiedPool{
		PoolID: pair.ID,
		DEX:    UniswapV2,
		Tokens: []UnifiedToken{
			{pair.Token0.ID, pair.Token0.Symbol, pair.Token0.Decimals},
			{pair.Token1.ID, pair.Token1.Symbol, pair.Token1.Decimals},
		},
		LiquidityUSD: parseStringToFloat(pair.ReserveUSD),
		VolumeUSD:    parseStringToFloat(pair.VolumeUSD),
		CreatedAt:    pair.CreatedAt,
	}
}


type PoolEvent struct {
	Pool UnifiedPool
}

var poolChan = make(chan PoolEvent, 1000)

func MonitorDex[T any](client *graphql.Client, query string, dex Dex, mapFunc func(T) []UnifiedPool) {
	for {
		resp, err := FetchGraphQL[T](client, query)
		if err != nil {
			fmt.Println("GraphQL error:", err)
			continue
		}
		pools := mapFunc(*resp)
		for _, p := range pools {
			poolChan <- PoolEvent{Pool: p}
		}
		time.Sleep(time.Second * 10) // 定时查询
	}
}

func PoolConsumer() {
	for event := range poolChan {
		// TODO: 套利策略处理
		fmt.Println("New pool:", event.Pool.PoolID, "DEX:", event.Pool.DEX)
	}
}