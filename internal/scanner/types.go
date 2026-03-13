package scanner

import (
	"github.com/Crypto-ChainSentinel/internal/parser/comm"
	"math/big"
)

type Interval struct {
	StartBlock *big.Int
	EndBlock   *big.Int
}

type ScannCfg struct {
	Interval
	Selected comm.ParserCfg
}
