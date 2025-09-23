package handler

import (
	"github.com/CryptoQuantX/chain_monitor/db"
	"github.com/CryptoQuantX/chain_monitor/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

func HandleNewContract(tx *types.Transaction, contractaddr common.Address, erc models.ERCStandard, BlockNumber *big.Int) {
	now := time.Now()
	txTime := tx.Time()
	ContractAge := now.Sub(txTime)
	ConstractInfo := models.ConstractInfo{
		Address:      contractaddr,
		ContractType: erc.Ercname,
		ContractAge:  ContractAge,
		TxHash:       tx.Hash(),
		BlockNumber:  BlockNumber,
		DeployTime:   tx.Time(),
	}

	//添加到合约监控池
	db.AddContract(&ConstractInfo)
}

func HandleNormalContract(tx *types.Transaction, from common.Address) {

	//判断合约监控池是否存在，存在则去除相关字段

	//构造ne04j和mysql字段插入
	var txnode models.TxNode
	var whaletranc models.WhaleTransaction

	txnode.BuildTxnode(tx, from, "whale", "normalcontract", "")
	whaletranc.BuildWhaleTransaction()

	db.Add(&txnode, &whaletranc)
}
