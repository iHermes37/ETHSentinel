// Package parser 提供解析引擎（Engine），负责：
//  1. 持有 ProtocolManager 全局注册表
//  2. 提供统一的协议注册入口
//  3. 根据 ParserCfg 构建运行时 ActiveParsers
//
// 重构要点：
//   - 去掉原来 Parser.RegisterAllParser() 的隐式全局副作用
//   - 使用 functional options 支持灵活配置
//   - Engine 是可复用的无状态对象（线程安全）
package parser

import (
	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ETHSentinel/internal/parser/dex"
	"github.com/ETHSentinel/internal/parser/token"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Engine 解析引擎：持有 ProtocolManager，对外提供构建 ActiveParsers 的能力。
type Engine struct {
	mgr *comm.ProtocolManager
}

// // NewEngine 创建并初始化解析引擎，默认注册所有内置协议。
//
//	func NewEngine() (*Engine, error) {
//		mgr := comm.NewProtocolManager()
//
//		// ── DEX 协议 ──────────────────────────────
//		dexMgr := comm.NewProtocolImplManager()
//		if err := dex.RegisterAll(dexMgr); err != nil {
//			return nil, err
//		}
//		mgr.Register(comm.ProtocolTypeDEX, dexMgr)
//
//		// ── Token 协议 ────────────────────────────
//		tokenMgr := comm.NewProtocolImplManager()
//		if err := token.RegisterAll(tokenMgr); err != nil {
//			return nil, err
//		}
//		mgr.Register(comm.ProtocolTypeToken, tokenMgr)
//
//		return &Engine{mgr: mgr}, nil
//	}
//
// 改为接收 client
func NewEngine(client *ethclient.Client) (*Engine, error) {
	mgr := comm.NewProtocolManager()

	dexMgr := comm.NewProtocolImplManager()
	if err := dex.RegisterAll(dexMgr, client); err != nil { // 传入 client
		return nil, err
	}
	mgr.Register(comm.ProtocolTypeDEX, dexMgr)

	tokenMgr := comm.NewProtocolImplManager()
	if err := token.RegisterAll(tokenMgr); err != nil {
		return nil, err
	}
	mgr.Register(comm.ProtocolTypeToken, tokenMgr)

	return &Engine{mgr: mgr}, nil
}

// Manager 返回底层 ProtocolManager（用于高级扩展）
func (e *Engine) Manager() *comm.ProtocolManager {
	return e.mgr
}

// BuildActive 根据配置构建运行时解析器集合
func (e *Engine) BuildActive(cfg comm.ParserCfg) (*comm.ActiveParsers, error) {
	return e.mgr.BuildActive(cfg)
}
