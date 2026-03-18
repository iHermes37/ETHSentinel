// Package token 统一注册所有 Token 协议实现（ERC20 / ERC721）
package token

import (
	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ETHSentinel/internal/parser/token/erc20"
	"github.com/ETHSentinel/internal/parser/token/erc721"
)

// RegisterAll 将所有 Token 实现注册到 mgr 中
func RegisterAll(mgr *comm.ProtocolImplManager) error {
	if err := mgr.RegisterStrategy(comm.ProtocolImplERC20, erc20.NewParser()); err != nil {
		return err
	}
	if err := mgr.RegisterStrategy(comm.ProtocolImplERC721, erc721.NewParser()); err != nil {
		return err
	}
	return nil
}
