package handler

import (
	"github.com/Crypto-ChainSentinel/db"
	"github.com/Crypto-ChainSentinel/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func HandleNewWhale(tx *types.Transaction, from common.Address) {
	newWhale := models.Whale{
		Address:   from,
		FirstSeen: tx.Time(),
		Note:      "Whale",
	}
	db.AddWhale(&newWhale)
}

func HandleNewAddr(tx *types.Transaction, from common.Address) {
	var txnode models.TxNode
	txnode.BuildTxnode(tx, from, "Whale", "Normal", "TRANSFER")

	//插入mysql数据库的数据构造
	var whaletranc models.WhaleTransaction
	whaletranc.BuildWhaleTransaction()

	db.Add(&txnode, &whaletranc)

}

func HandleCex(tx *types.Transaction, from common.Address, cex string) {
	var txnode models.TxNode
	txnode.BuildTxnode(tx, from, "Whale", cex, "TRANSFER")

	//插入mysql数据库的数据构造

	db.Add(&txnode)

}
