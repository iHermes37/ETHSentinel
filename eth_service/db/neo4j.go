package db

import (
	"context"
	"fmt"
	"github.com/CryptoQuantX/chain_monitor/initialize"
	"github.com/CryptoQuantX/chain_monitor/models"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// 更新插入节点到ne04j图数据库中的
func UpsertNode(Txnode *models.TxNode) error {
	driver := initialize.Neo4jconn()
	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]any{
			"whaleAddress":   Txnode.WhaleNode.Address.Hex(),
			"whaleFirstSeen": Txnode.WhaleNode.FirstSeen,
			"whaleLastSeen":  Txnode.WhaleNode.LastSeen,
			"whaleLabel":     Txnode.WhaleNode.Label,

			"toAddress":   Txnode.ToNode.Address.Hex(),
			"toFirstSeen": Txnode.ToNode.FirstSeen,
			"toLastSeen":  Txnode.ToNode.LastSeen,
			"toLabel":     Txnode.ToNode.Label,

			"LastTxTs": Txnode.Rp.LastTxTs,
			"RxType":   Txnode.Rp.Type,
		}

		cypher := fmt.Sprintf(
			`
			MERGE (w:$whaleLabel {address:$whaleAddress})
			  ON CREATE SET w.first_seen=$whaleFirstSeen, w.last_seen=$whaleLastSeen, w.tx_count=1
			  ON MATCH  SET w.last_seen=$whaleLastSeen   w.tx_count=w.tx_count+1
			
			MERGE (t:$toLabel {address:$toAddress})
			  ON CREATE SET t.first_seen=$toFirstSeen, t.last_seen=$toLastSeen, t.tx_count=1
			  ON MATCH  SET t.last_seen=$toLastSeen	,t.tx_count = t.tx_count + 1
			
			MERGE (w)-[r:$RxType]->(t)
			  ON CREATE SET r.tx_count=1, r.last_tx=$LastTxTs
			  ON MATCH  SET r.tx_count=r.tx_count+1, r.last_tx=$LastTxTs
			`)

		return tx.Run(ctx, cypher, params)
	})
	return err
}

func AddToNeo4j(tx *models.TxNode) {
	UpsertNode(tx)
}
