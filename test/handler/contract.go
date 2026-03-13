package handler

import (
	"github.com/Crypto-ChainSentinel/server/db"
	db2 "github.com/Crypto-ChainSentinel/test/db"
	"github.com/Crypto-ChainSentinel/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

func HandleNewContract(tx *types.Transaction, contractaddr common.Address, erc types.ERCStandard, BlockNumber *big.Int) {
	now := time.Now()
	txTime := tx.Time()
	ContractAge := now.Sub(txTime)
	ConstractInfo := types.ConstractInfo{
		Address:      contractaddr,
		ContractType: erc.Ercname,
		ContractAge:  ContractAge,
		TxHash:       tx.Hash(),
		BlockNumber:  BlockNumber,
		DeployTime:   tx.Time(),
	}

	//添加到合约监控池
	db2.AddContract(&ConstractInfo)
}

func HandleNormalContract(tx *types.Transaction, from common.Address) {

	//判断合约监控池是否存在，存在则去除相关字段

	//构造ne04j和mysql字段插入
	var txnode types.TxNode
	var whaletranc types.WhaleTransaction

	txnode.BuildTxnode(tx, from, "whale", "normalcontract", "")
	whaletranc.BuildWhaleTransaction()

	db.Add(&txnode, &whaletranc)
}
