//package handler
//
//import (
//	"github.com/Crypto-ETHSentinel/db"
//	"github.com/Crypto-ETHSentinel/types"
//	"github.com/Crypto-ETHSentinel/internal/parser"
//	"github.com/ethereum/go-ethereum/commonParser"
//	"github.com/ethereum/go-ethereum/core/types"
//)
//
//func HandleDex(tx *types.Transaction, from commonParser.Address, dexname string) {
//
//	// 解析交易输入数据，确定代币和数量
//	dextxinfo := parser.ParseDexInput(tx, dexname)
//	var wt = types.WhaleTransaction{}
//	var txnode = types.TxNode{}
//	wt.BuildWhaleTransaction(tx, from, &dextxinfo)
//	txnode.BuildTxnode(tx, from, "whale", "DEX", "")
//
//	// 更新 MySQL
//	db.Add(&wt, &txnode)
//
//}
//
//func HandleDeFi(tx *types.Transaction, from commonParser.Address, defiName string) {
//	if defiName == "dex" {
//		HandleDex(tx, from, defiName)
//	}
//}
