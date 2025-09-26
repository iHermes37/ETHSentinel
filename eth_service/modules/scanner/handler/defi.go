package handler

import (
	"github.com/CryptoQuantX/chain_monitor/db"
	"github.com/CryptoQuantX/chain_monitor/models"
	"github.com/CryptoQuantX/chain_monitor/modules/parser"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func HandleDex(tx *types.Transaction, from common.Address, dexname string) {

	// 解析交易输入数据，确定代币和数量
	dextxinfo := parser.ParseDexInput(tx, dexname)
	var wt = models.WhaleTransaction{}
	var txnode = models.TxNode{}
	wt.BuildWhaleTransaction(tx, from, &dextxinfo)
	txnode.BuildTxnode(tx, from, "whale", "DEX", "")

	// 更新 MySQL
	db.Add(&wt, &txnode)

}

func HandleDeFi(tx *types.Transaction, from common.Address, defiName string) {
	if defiName == "dex" {
		HandleDex(tx, from, defiName)
	}
}
