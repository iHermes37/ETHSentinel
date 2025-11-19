package Conn

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func Neo4jconn() neo4j.DriverWithContext {
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
