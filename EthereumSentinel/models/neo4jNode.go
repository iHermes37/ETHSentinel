package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
