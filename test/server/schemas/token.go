package schemas

// 单个地址持仓信息
type TokenHolder struct {
	Address string  `json:"address"` // 持币地址
	Balance float64 `json:"balance"` // 持币数量
	Share   float64 `json:"share"`   // 占代币总量比例（百分比）
}

// 代币分布统计
type TokenDistribution struct {
	TokenAddress    string        `json:"tokenAddress"`    // 代币合约地址
	TokenSymbol     string        `json:"tokenSymbol"`     // 代币符号
	TotalSupply     float64       `json:"totalSupply"`     // 代币总量
	HolderCount     int           `json:"holderCount"`     // 持币地址总数
	TopHolders      []TokenHolder `json:"topHolders"`      // 前 N 大持币地址
	GiniCoefficient float64       `json:"giniCoefficient"` // 基尼系数，衡量持仓集中度
	SnapshotTime    int64         `json:"snapshotTime"`    // 数据快照时间戳
}
