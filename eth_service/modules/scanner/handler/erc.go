package handler

import (
	"github.com/Crypto-ChainSentinel/db"
	"github.com/Crypto-ChainSentinel/models"
)

func HandleERCContract(msg *models.ERCStandard) {

	ERCTxD := models.ERC20TxDetail{}
	var wt = models.WhaleTransaction{}
	wt.BuildWhaleTransaction()

	var txnode models.TxNode
	txnode.BuildTxnode()
	// -----------双数据存储处理------------------------
	db.Add(&wt, &txnode)

}
