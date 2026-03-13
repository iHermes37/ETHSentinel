package core

import "github.com/ethereum/go-ethereum/core/types"

func Distribute(processor Process) {
	EthCh := make(chan *types.Transaction, 100)
	TokenCh := make(chan *types.Receipt, 100)
	DefiCh := make(chan *types.Receipt, 100)
	NewContractCh := make(chan *types.Receipt, 100)

	select {
	case tx := <-EthCh:
		processor.ProcessEthTransaction(tx)
	case receipt := <-TokenCh:
		processor.ProcessTokenTransaction(receipt)
	case receipt := <-DefiCh:
		processor.ProcessDefiTransaction(receipt)
	case receipt := <-NewContractCh:
		processor.ProcessNewContract(receipt)
	}
}

type Process interface {
	ProcessEthTransaction(tx *types.Transaction)
	ProcessTokenTransaction(rec *types.Receipt)
	ProcessDefiTransaction(rec *types.Receipt)
	ProcessNewContract(rec *types.Receipt)
}
