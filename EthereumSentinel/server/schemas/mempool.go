package model

import "time"

// 严格来说，在以太坊网络中 Mempool（内存池）不是单一“全局池”，
// 而是 每个节点都有自己的交易池。所以从概念上讲，可以有多个 Mempool，但它们是 分布式、节点级别的：

// 单条 Mempool 交易
type MempoolTx struct {
	TxHash    string    `json:"txHash"`              // 交易哈希
	From      string    `json:"from"`                // 发起地址
	To        *string   `json:"to,omitempty"`        // 接收地址，可为空
	Value     float64   `json:"value"`               // 交易金额（ETH 或代币数量）
	GasPrice  float64   `json:"gasPrice"`            // Gas 单价（Gwei）
	GasLimit  uint64    `json:"gasLimit"`            // Gas 限额
	Nonce     uint64    `json:"nonce"`               // 交易序号
	InputData *string   `json:"inputData,omitempty"` // 调用数据，可解析方法
	Type      string    `json:"type"`                // 交易类型：普通 / 合约调用 / 内部转账等
	Timestamp time.Time `json:"timestamp"`           // 交易被检测到的时间
	PoolTime  time.Time `json:"poolTime"`            // 入池时间
}

// Mempool 监控统计
type MempoolStats struct {
	TotalTx      int            `json:"totalTx"`      // 当前池中交易总数
	AvgGasPrice  float64        `json:"avgGasPrice"`  // 平均 Gas 价格
	MaxGasPrice  float64        `json:"maxGasPrice"`  // 最大 Gas 价格
	MinGasPrice  float64        `json:"minGasPrice"`  // 最小 Gas 价格
	TxByType     map[string]int `json:"txByType"`     // 各类型交易数量统计
	TopSenders   []string       `json:"topSenders"`   // 持续入池交易量最大的地址
	DetectSpikes bool           `json:"detectSpikes"` // 是否检测到交易量突增
	SnapshotTime time.Time      `json:"snapshotTime"` // 数据快照时间
}
