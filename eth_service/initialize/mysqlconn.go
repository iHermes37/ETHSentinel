package initialize

import (
	"github.com/Crypto-ChainSentinel/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func InitMysql() *gorm.DB {
	dsn := "user:password@tcp(127.0.0.1:3306)/whale_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
		panic("数据库连接失败: " + err.Error())
	}
	if err := db.AutoMigrate(&models.ConstractInfo{}); err != nil {
		log.Fatal("migrate error:", err)
	}
	if err := db.AutoMigrate(&models.WhaleTransaction{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}
	if err := db.AutoMigrate(&models.Whale{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}
	if err := db.AutoMigrate(&models.CrossPairData{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}
	return db
}
