package mempool

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

type Mempool struct {
	client *ethclient.Client
}

func NewMempool() *Mempool {
	return &Mempool{}
}
