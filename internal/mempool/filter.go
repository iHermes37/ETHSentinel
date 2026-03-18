package mempool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ─────────────────────────────────────────────
//  内置过滤器工厂函数
// ─────────────────────────────────────────────

// FilterByToAddress 只保留发往指定地址的交易（合约调用过滤）
func FilterByToAddress(addrs ...common.Address) FilterFunc {
	set := make(map[common.Address]struct{}, len(addrs))
	for _, a := range addrs {
		set[a] = struct{}{}
	}
	return func(tx *types.Transaction) bool {
		if tx.To() == nil {
			return false
		}
		_, ok := set[*tx.To()]
		return ok
	}
}

// FilterByMinValue 只保留 value >= minWei 的 ETH 转账
func FilterByMinValue(minWei *big.Int) FilterFunc {
	return func(tx *types.Transaction) bool {
		return tx.Value().Cmp(minWei) >= 0
	}
}

// FilterByMethodSig 只保留 input data 前4字节匹配的交易
// sig 例如: "0x38ed1739"（swapExactTokensForTokens）
func FilterByMethodSig(sigs ...string) FilterFunc {
	set := make(map[string]struct{}, len(sigs))
	for _, s := range sigs {
		if len(s) >= 10 {
			set[s[:10]] = struct{}{}
		}
	}
	return func(tx *types.Transaction) bool {
		data := tx.Data()
		if len(data) < 4 {
			return false
		}
		sig := common.Bytes2Hex(data[:4])
		_, ok := set["0x"+sig]
		return ok
	}
}

// FilterByMinGasPrice 只保留 gasPrice >= min 的交易（用于抢跑检测）
func FilterByMinGasPrice(minGwei int64) FilterFunc {
	minWei := new(big.Int).Mul(big.NewInt(minGwei), big.NewInt(1e9))
	return func(tx *types.Transaction) bool {
		return tx.GasPrice().Cmp(minWei) >= 0
	}
}

// FilterContractCreation 只保留合约部署交易
func FilterContractCreation() FilterFunc {
	return func(tx *types.Transaction) bool {
		return tx.To() == nil
	}
}

// FilterHasData 只保留有 input data 的交易（合约调用）
func FilterHasData() FilterFunc {
	return func(tx *types.Transaction) bool {
		return len(tx.Data()) > 0
	}
}
