// Package scanner — 扫描器配置类型
package scanner

import (
	"math/big"

	"github.com/ETHSentinel/internal/parser/comm"
)

// ScanBlockCfg 单块扫描配置
type ScanBlockCfg struct {
	// ActiveParsers 由 parser.Engine.BuildActive(ParserCfg) 提供
	// 若为 nil，则跳过事件解析，仅做交易分类
	Active *comm.ActiveParsers
}

// ScanBlocksCfg 区间扫描配置
type ScanBlocksCfg struct {
	StartBlock *big.Int
	EndBlock   *big.Int
	// ParserCfg 解析配置，Scanner 内部会调用 Engine.BuildActive
	ParserCfg      comm.ParserCfg
	WorkerPoolSize int // 协程池大小，默认 100
}

// TxCategory 交易分类
type TxCategory string

const (
	TxCategoryETH         TxCategory = "ETH"         // 普通以太转账
	TxCategoryToken       TxCategory = "Token"       // ERC20/721 代币交互
	TxCategoryDeFi        TxCategory = "DeFi"        // DeFi 协议调用
	TxCategoryNewContract TxCategory = "NewContract" // 合约部署
)
