package db

import (
	"context"
	"fmt"
	"github.com/Crypto-ChainSentinel/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

type Node struct {
	// 节点唯一标识（通常是地址）
	Address common.Address `json:"address"`
	// 节点类型（Neo4j 标签）
	Label string `json:"label"`
	// 创建时间/首次发现时间
	FirstSeen time.Time `json:"first_seen,omitempty"`
	// 最近交易时间/最后活动时间
	LastSeen time.Time `json:"last_seen,omitempty"`
	// 累计交易次数（可选）
	TxCount int `json:"tx_count,omitempty"`
	// 其他可扩展属性
	Extra map[string]any `json:"extra,omitempty"`
}

type Relationship struct {
	Type     string         `json:"type"`     // INTERACTS_WITH
	TxCount  int            `json:"tx_count"` // 累计交易次数
	LastTxTs time.Time      `json:"last_tx"`  // 最近交易时间
	Extra    map[string]any `json:"extra,omitempty"`
}

type TxNode struct {
	WhaleNode Node         `json:"whaleNode"`
	ToNode    Node         `json:"toNode"`
	Rp        Relationship `json:"relationship"`
}

func (txnode *TxNode) BuildTxnode(tx *types.Transaction,
	from common.Address, fromLabel string,
	toLabel string,
	RpType string) {

	txnode.WhaleNode = Node{
		Address:   from,
		Label:     fromLabel,
		FirstSeen: tx.Time(),
		LastSeen:  tx.Time(),
	}

	txnode.ToNode = Node{
		Address:   *tx.To(),
		Label:     toLabel,
		FirstSeen: tx.Time(),
		LastSeen:  tx.Time(),
	}

	txnode.Rp = Relationship{
		LastTxTs: tx.Time(),
		Type:     RpType,
	}
}

// 更新插入节点到ne04j图数据库中的
func UpsertNode(Txnode *types.TxNode) error {
	driver := Conn.Conn.Neo4jconn()
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

func AddToNeo4j(tx *types.TxNode) {
	UpsertNode(tx)
}
