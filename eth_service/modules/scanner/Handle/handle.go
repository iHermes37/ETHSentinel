package Handle

import (
	"fmt"
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
	"time"

	"github.com/Crypto-ChainSentinel/db"
	"github.com/Crypto-ChainSentinel/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func HandleNewContract(tx *types.Transaction, contractaddr common.Address) {
	now := time.Now()
	txTime := tx.Time()
	ContractAge := now.Sub(txTime)
	ConstractInfo := models.ConstractInfo{
		Address:     contractaddr,
		ContractAge: ContractAge,
		TxHash:      tx.Hash(),
		DeployTime:  tx.Time(),
	}
	//添加到合约监控池
	db.AddContract(&ConstractInfo)
}

func HandleNewWhale(tx *types.Transaction, from common.Address) {
	newWhale := models.Whale{
		Address:   from,
		FirstSeen: tx.Time(),
		Note:      "Whale",
	}
	db.AddWhale(&newWhale)
}

func HandleEvents(ch chan []ParserEngineCommon.UnifiedEvent) {
	for events := range ch {
		fmt.Println("处理事件数量：", len(events))
	}
}
