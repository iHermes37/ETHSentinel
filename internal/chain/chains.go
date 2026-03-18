package chain

// Ethereum 返回以太坊主网配置
func Ethereum() Chain {
	return NewBaseChain(
		ChainETH,
		"ethereum",
		"ETH",
		"https://mainnet.infura.io/v3/YOUR_KEY",
		"wss://mainnet.infura.io/ws/v3/YOUR_KEY",
	)
}

// BSC 返回 BNB Smart Chain 配置
func BSC() Chain {
	return NewBaseChain(
		ChainBSC,
		"bsc",
		"BNB",
		"https://bsc-dataseed.binance.org",
		"wss://bsc-ws-node.nariox.org:443",
	)
}

// Polygon 返回 Polygon 主网配置
func Polygon() Chain {
	return NewBaseChain(
		ChainPolygon,
		"polygon",
		"MATIC",
		"https://polygon-rpc.com",
		"wss://polygon-bor.publicnode.com",
	)
}

// Arbitrum 返回 Arbitrum One 配置
func Arbitrum() Chain {
	return NewBaseChain(
		ChainArbitrum,
		"arbitrum",
		"ETH",
		"https://arb1.arbitrum.io/rpc",
		"wss://arb1.arbitrum.io/ws",
	)
}
