package db

import (
	"context"
	"fmt"
	db2 "github.com/Crypto-ChainSentinel/test/db"
	"github.com/Crypto-ChainSentinel/types"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func Add[T any, T1 any](data *T, data1 *T1) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		return db2.AddToMysql(data) // 返回 error
	})

	//g.Go(func() error {
	//	return db2.AddToNeo4j(&data) // 返回 error
	//})

	if err := g.Wait(); err != nil {
		// 任意一个失败，你可以在这里做补偿或回滚
		return err
	}
	return nil
}

func InitMysql() *gorm.DB {
	dsn := "user:password@tcp(127.0.0.1:3306)/whale_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to conn database:", err)
		panic("数据库连接失败: " + err.Error())
	}
	if err := db.AutoMigrate(&types.ConstractInfo{}); err != nil {
		log.Fatal("migrate error:", err)
	}
	//if err := db.AutoMigrate(&types.WhaleTransaction{}); err != nil {
	//	log.Fatal("failed to migrate:", err)
	//}
	if err := db.AutoMigrate(&types.Whale{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}
	if err := db.AutoMigrate(&types.CrossPairData{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}
	return db
}

func InitNeo4j() neo4j.DriverWithContext {
	ctx := context.Background()
	// URI examples: "neo4j://", "neo4j+s://xxx.databases.neo4j.io"
	dbUri := "neo4j://localhost:7687"
	dbUser := "neo4j"
	dbPassword := "strongpass123"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	fmt.Println("Connection established.")

	return driver
}

func InitDb() {
	InitMysql()
	InitNeo4j()
}
