package sentinel

import "github.com/ETHSentinel/internal/parser/comm"

// SwapData UniswapV2 Swap 详情（透出给 sentinel 调用方）
type SwapData = comm.SwapData

// TransferData ERC20/ETH 转账详情（透出给 sentinel 调用方）
type TransferData = comm.TransferData
