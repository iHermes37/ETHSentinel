// Package dex 统一注册所有 DEX 协议实现
package dex

import (
	"github.com/ETHSentinel/internal/parser/comm"
	uniswapv2 "github.com/ETHSentinel/internal/parser/dex/uniswap_v2"
	"github.com/ethereum/go-ethereum/ethclient"
)

// RegisterAll 将所有 DEX 实现注册到 mgr 中
func RegisterAll(mgr *comm.ProtocolImplManager, client *ethclient.Client) error {
	if err := mgr.RegisterStrategy(comm.ProtocolImplUniswapV2, uniswapv2.NewParser(client)); err != nil {
		return err
	}
	// 未来扩展：SushiSwap、Curve、Balancer 等
	// mgr.RegisterStrategy(comm.ProtocolImplSushiSwap, sushiswap.NewParser())
	return nil
}
