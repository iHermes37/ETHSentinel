package db

import (
	"github.com/Crypto-ChainSentinel/types"
	"github.com/ethereum/go-ethereum/common"
)

func AddContract(c *types.ConstractInfo) {

}

func AddWhale(w *types.Whale) {

}

type Curd interface {
	StoreWhaleToMonitorPool()
}

type RedisMgr struct {
}

func NewRedisMgr() *RedisMgr {
	return &RedisMgr{}
}

func (r *RedisMgr) StoreWhaleToMonitorPool(whale *common.Address) {

}

func (r *RedisMgr) IsWhaleInMonitorPool(whale *common.Address) bool {
	return true
}

func (r *RedisMgr) UpdateWhaleTransactionCount(whale *common.Address) error {
	return nil
}

func (r *RedisMgr) UpdateWhaleDeFiTransactionCount(whale *common.Address) error {
	return nil
}

// ============================================================================
func (r *RedisMgr) StoreContractToMonitorPool(contract *common.Address) {

}

func (r *RedisMgr) IsContractInMonitorPool(contract *common.Address) bool {

}

func (r *RedisMgr) UpdateContractTransactionCount(whale *common.Address) error {
	return nil
}

func (r *RedisMgr) UpdateContractDeFiTranCount(whale *common.Address) error {
	return nil
}
