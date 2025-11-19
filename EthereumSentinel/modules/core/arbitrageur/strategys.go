package arbitrageur

import "github.com/Crypto-ChainSentinel/modules/core/arbitrageur/strategy"

var Strategies = map[string]Strategy{
	"CrossDEX": &strategy.CrossDEXStrategy{},
}

type Strategy interface {
	Name() string
	Run(stop <-chan struct{})
}
