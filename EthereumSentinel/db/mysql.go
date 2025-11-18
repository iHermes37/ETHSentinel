package db

import (
	"github.com/Crypto-ChainSentinel/init"
	"github.com/ethereum/go-ethereum/common"
	"log"
)

func AddToMysql[T any](data *T) error {
	db := init.InitMysql() // 全局复用连接
	if err := db.Create(data).Error; err != nil {
		log.Println("insert error:", err)
		return err
	}
	return nil
}

type MysqlMgr struct {
}

func (m *MysqlMgr) IsWhaleInPool(addr *common.Address) bool {
	return true
}
