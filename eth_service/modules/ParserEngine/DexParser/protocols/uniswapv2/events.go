package uniswapv2

import (
	abligens "github.com/Crypto-ChainSentinel/modules/ParserEngine/DexParser/abigens"
	"github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
)

type UniswapV2_SwapEvent struct {
	common.EventMetadata // 嵌入 BaseEvent 自动实现 UnifiedEvent
	Event                abligens.UniswappairSwap
}

type UniswapV2_MintEvent struct {
	event abligens.UniswappairMint
}

type UniswapV2_BurnEvent struct {
	event abligens.UniswappairBurn
}

type UniswapV2_SyncEvent struct {
	event abligens.UniswappairSync
}

type UniswapV2_TransferEvent struct {
	event abligens.UniswappairTransfer
}

type UniswapV2_PairCreatedEvent struct {
	event abligens.UniswapPairCreated
}
