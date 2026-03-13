package Filter

import (
	"github.com/Crypto-ChainSentinel/test/db"
	"time"

	"github.com/Crypto-ChainSentinel/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func HandleNewContract(tx *types.Transaction, contractaddr common.Address) {
	now := time.Now()
	txTime := tx.Time()
	ContractAge := now.Sub(txTime)
	ConstractInfo := types.ConstractInfo{
		Address:     contractaddr,
		ContractAge: ContractAge,
		TxHash:      tx.Hash(),
		DeployTime:  tx.Time(),
	}
	//添加到合约监控池
	db.AddContract(&ConstractInfo)
}

func HandleNewWhale(tx *types.Transaction, from common.Address) {

}

func HandleCex() {

}

func HandleNewAddr() {

}
