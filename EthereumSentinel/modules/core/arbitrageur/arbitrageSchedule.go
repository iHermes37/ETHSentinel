package arbitrageur

var strategies = map[string]ArbitrageStrategy{
	"CrossDEX":   &models.CrossDEXStrategy{},
	"Triangular": &models.TriangularStrategy{},
	"CrossChain": &models.CrossChainStrategy{},
}

func arbitrageSchedule(selected []string) {
	stopChans := make(map[string]chan struct{})
	for _, name := range selected {
		if s, ok := strategies[name]; ok {
			stop := make(chan struct{})
			stopChans[name] = stop
			go s.Run(stop)
		}
	}
	// 阻塞等待退出
	select {}
}

// 事件触发：监听任意一个池子的 Swap
// 更新缓存：只更新这个池子的储备
// 拉取/缓存：获取其他两条交易路径价格（缓存或实时）
// 边抓边算：计算三角套利利润
// 安全判断：利润 > 手续费 + 滑点阈值
// 调用原子合约：一次交易执行买卖
// 如果链上实际利润不足 → revert
// 完成或失败：不会亏钱，等待下一次事件

// 监控多链 DEX 价格
// 在链 A、链 B 分别监听目标交易对价格（用事件驱动或定时轮询）。
// 计算跨链差价
// 假设链 A 的 USDT/WETH 价格低于链 B → 可以低买高卖。
// 资产跨链
// 将买入的资产通过跨链桥（如 Wormhole, LayerZero, Celer 等）发送到目标链。
// 在目标链卖出
// 收回资产，实现盈利。
// 风险控制
// 跨链桥延迟 + 波动 → 盈利需要覆盖桥费和滑点。
// 适合 中低频套利，不像单链套利那样毫秒级。

// 目标：借用大量资金在同一区块内完成套利，不需要自有资金。
// 特点：原子交易，一次交易完成借贷 + 买入 + 卖出 + 还贷。
// 前提：必须在单一区块内完成所有操作，否则贷款会回滚。
// 实现步骤
// 借贷
// 使用 Aave、dYdX、Uniswap V3 等平台的闪电贷接口，一次性借入大量资产。
// 套利操作
// 使用借来的资金在不同 DEX 或不同交易对进行套利（跨DEX、三角等）。
// 归还贷款
// 将套利获得的资产归还闪电贷 + 支付手续费。
// 剩余利润
// 剩余就是净收益。
// 安全保证
// 如果在交易执行过程中利润不足，整个交易会 revert → 不亏本金。
// 小结
// 关键是 原子交易 + 足够的套利差价覆盖手续费。
// 不需要自有资金，非常适合高频套利和抢先交易。
// 和你之前讲的原子合约套利逻辑很类似，但加入了 借贷步骤。
