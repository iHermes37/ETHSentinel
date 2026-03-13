package handler

import (
	"github.com/Crypto-ChainSentinel/server/db"
	db2 "github.com/Crypto-ChainSentinel/test/db"
	"github.com/Crypto-ChainSentinel/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func HandleNewWhale(tx *types.Transaction, from common.Address) {
	newWhale := types.Whale{
		Address:   from,
		FirstSeen: tx.Time(),
		Note:      "Whale",
	}
	db2.AddWhale(&newWhale)
}

func HandleNewAddr(tx *types.Transaction, from common.Address) {
	var txnode types.TxNode
	txnode.BuildTxnode(tx, from, "Whale", "Normal", "TRANSFER")

	//插入mysql数据库的数据构造
	var whaletranc types.WhaleTransaction
	whaletranc.BuildWhaleTransaction()

	db.Add(&txnode, &whaletranc)

}

func HandleCex(tx *types.Transaction, from common.Address, cex string) {
	var txnode types.TxNode
	txnode.BuildTxnode(tx, from, "Whale", cex, "TRANSFER")

	//插入mysql数据库的数据构造

	db.Add(&txnode)

}
