package db

import (
	"github.com/CryptoQuantX/chain_monitor/initialize"
	"log"
)

func AddToMysql[T any](data *T) error {
	db := initialize.InitMysql() // 全局复用连接
	if err := db.Create(data).Error; err != nil {
		log.Println("insert error:", err)
		return err
	}
	return nil
}
