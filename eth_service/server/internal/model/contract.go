package model

import (
	"github.com/CryptoQuantX/chain_monitor/models"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type ContractQueryParams struct {
	Address      *common.Address      `json:"address"`      // 可选
	ContractType *models.ContractType `json:"contractType"` // 对应 JSON 驼峰命名
	ContractAge  *time.Duration       `json:"contractAge"`
	DeployTime   *time.Time           `json:"deployTime"`
	TxHash       *common.Hash         `json:"txHash"`
}

type ContractQueryResponse struct {
	Address      common.Address      `json:"address"`
	ContractType models.ContractType `json:"contractType"`
	ContractAge  time.Duration       `json:"contractAge"`
	DeployTime   time.Time           `json:"deployTime"`
	TxHash       common.Hash         `json:"txHash"`
}

type NewTokenPair struct {
	PairAddress   string    `json:"pairAddress"`   // 代币对合约地址
	Token0        string    `json:"token0"`        // 代币0的合约地址
	Token1        string    `json:"token1"`        // 代币1的合约地址
	Token0Symbol  string    `json:"token0Symbol"`  // 代币0的符号
	Token1Symbol  string    `json:"token1Symbol"`  // 代币1的符号
	PairCreatedAt time.Time `json:"pairCreatedAt"` // 创建时间
	Factory       string    `json:"factory"`       // 创建该交易对的工厂合约地址
	TxHash        string    `json:"txHash"`        // 创建交易哈希
}

type LiquidityChangeType struct {
	ADD    int
	REMOVE int
}

type LiquidityChange struct {
	PoolAddress  string              `json:"poolAddress"`  // 资金池合约地址
	Token0       string              `json:"token0"`       // 代币0地址
	Token1       string              `json:"token1"`       // 代币1地址
	Token0Symbol string              `json:"token0Symbol"` // 代币0符号
	Token1Symbol string              `json:"token1Symbol"` // 代币1符号
	Amount0      float64             `json:"amount0"`      // 代币0变化数量
	Amount1      float64             `json:"amount1"`      // 代币1变化数量
	ChangeType   LiquidityChangeType `json:"changeType"`   // ADD 或 REMOVE
	Time         time.Time           `json:"time"`         // 发生时间
	TxHash       string              `json:"txHash"`       // 对应交易哈希
	Sender       string              `json:"sender"`       // 操作者地址
}
