package handler

import (
	"github.com/Crypto-ChainSentinel/server/db"
	"github.com/Crypto-ChainSentinel/types"
)

func HandleERCContract(msg *types.ERCStandard) {

	ERCTxD := types.ERC20TxDetail{}
	var wt = types.WhaleTransaction{}
	wt.BuildWhaleTransaction()

	var txnode types.TxNode
	txnode.BuildTxnode()
	// -----------双数据存储处理------------------------
	db.Add(&wt, &txnode)

}
