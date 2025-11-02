package models

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/machinebox/graphql"
	"math/big"
	"time"
)

type DeFi struct {
	Dex Dex `json:"Dex"`
}

type Dex struct {
	Uni   Uniswap   `json:"Uniswap"`
	Sushi SushiSwap `json:"SushiSwap"`
}

//-------------------uniswap--------------------------

// 定义最内层的 V2/V3 结构
type UniswapV2 struct {
	Router  string `json:"Router"`
	Factory string `json:"Factory"`
	Pair    string `json:"Pair"`
}

type UniswapV3 struct {
	Router  string `json:"Router"`
	Factory string `json:"Factory"`
	Pool    string `json:"Pool"`
}

// 定义 Uniswap 结构，里面嵌套 V2 和 V3
type Uniswap struct {
	V2 UniswapV2 `json:"V2"`
	V3 UniswapV3 `json:"V3"`
}

//-----------------SushiSwap---------------------------------------

type SushiSwap struct {
	Router  string `json:"Router"`
	Factory string `json:"Factory"`
	Pair    string `json:"Pair"`
}

// --------------------------------------------------
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

type Uniswap_V3 struct {
	GraphQLDEX
}

type Sushiswap struct {
	GraphQLDEX
}

type DexPairLiquidity struct {
	PairAddress common.Address `json:"pairAddress"` // 交易对合约地址
	Token0      Token          `json:"token0"`
	Token1      Token          `json:"token1"`
	Reserve0    *big.Int       `json:"reserve0"`    // token0 储备量
	Reserve1    *big.Int       `json:"reserve1"`    // token1 储备量
	TotalSupply *big.Int       `json:"totalSupply"` // LP 总量
	Timestamp   time.Time      `json:"timestamp"`   // 事件时间
}

type DexPairLiquidityChange struct {
	PairAddress common.Address `json:"pairAddress"`
	Token0      Token          `json:"token0"`
	Token1      Token          `json:"token1"`

	DeltaReserve0 *big.Int `json:"deltaReserve0"` // token0 变化量
	DeltaReserve1 *big.Int `json:"deltaReserve1"` // token1 变化量
	DeltaLP       *big.Int `json:"deltaLP"`       // LP 总量变化
	Reserve0After *big.Int `json:"reserve0After"` // 变化后的储备量
	Reserve1After *big.Int `json:"reserve1After"`
	TotalLPAfter  *big.Int `json:"totalLPAfter"`

	Timestamp time.Time `json:"timestamp"` // 变化发生时间
	EventType string    `json:"eventType"` // "Mint" / "Burn" / "Swap" / "ManualSnapshot"
}
