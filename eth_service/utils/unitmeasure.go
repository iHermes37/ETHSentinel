package utils

import (
	"fmt"
	"math/big"
)

func WeiToETH(val *big.Int) string {
	// 1 ETH = 10^18 wei
	eth := new(big.Float).Quo(new(big.Float).SetInt(val), big.NewFloat(1e18))
	return fmt.Sprintf("%.3f ETH", eth) // 保留 3 位小数
}
