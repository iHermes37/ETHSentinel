// Package comm — ProtocolManager（顶层注册表）
// 重构要点：
//   - 去掉 CurProtocols（职责不清），改为 BuildChain() 提供责任链构建
//   - 所有方法线程安全
//   - 提供 MustGet 快捷方法（测试/初始化阶段使用）
package comm

import (
	"fmt"
	"sync"
)

// ─────────────────────────────────────────────
//  ProtocolManager
// ─────────────────────────────────────────────

// ProtocolManager 全局协议注册表：ProtocolType → ProtocolTypeParser。
// 是整个解析引擎的入口，由 parser.Engine 持有。
type ProtocolManager struct {
	mu       sync.RWMutex
	registry map[ProtocolType]ProtocolTypeParser
}

// NewProtocolManager 创建空的协议管理器
func NewProtocolManager() *ProtocolManager {
	return &ProtocolManager{
		registry: make(map[ProtocolType]ProtocolTypeParser),
	}
}

// Register 注册一个协议大类的解析器
func (m *ProtocolManager) Register(protocolType ProtocolType, parser ProtocolTypeParser) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.registry[protocolType] = parser
}

// Get 获取协议大类解析器
func (m *ProtocolManager) Get(protocolType ProtocolType) (ProtocolTypeParser, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.registry[protocolType]
	if !ok {
		return nil, fmt.Errorf("protocol_mgr: protocol type %q not registered", protocolType)
	}
	return p, nil
}

// MustGet 获取协议大类解析器，不存在时 panic（仅用于初始化阶段）
func (m *ProtocolManager) MustGet(protocolType ProtocolType) ProtocolTypeParser {
	p, err := m.Get(protocolType)
	if err != nil {
		panic(err)
	}
	return p
}

// ListTypes 返回所有已注册的协议大类
func (m *ProtocolManager) ListTypes() []ProtocolType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	types := make([]ProtocolType, 0, len(m.registry))
	for k := range m.registry {
		types = append(types, k)
	}
	return types
}

// ─────────────────────────────────────────────
//  ActiveParsers — 运行时选中的解析器集合
// ─────────────────────────────────────────────

// ActiveParsers 由 ParserCfg 实例化，持有运行时真正参与解析的实现列表。
// 责任链由 scanner 包根据此结构构建。
type ActiveParsers struct {
	// Impls 按插入顺序排列，决定责任链的优先级
	Impls []ProtocolImplParser
}

// BuildActive 根据 ParserCfg 从 ProtocolManager 中选出所有激活的实现，
// 同时应用事件过滤器，返回 ActiveParsers。
func (m *ProtocolManager) BuildActive(cfg ParserCfg) (*ActiveParsers, error) {
	active := &ActiveParsers{}

	for protoType, implCfg := range cfg {
		typeParser, err := m.Get(protoType)
		if err != nil {
			return nil, fmt.Errorf("protocol_mgr: build active: %w", err)
		}

		for implName, selectedEvents := range implCfg {
			implParser, err := typeParser.GetImpl(implName)
			if err != nil {
				return nil, fmt.Errorf("protocol_mgr: build active: %w", err)
			}
			// 应用事件过滤
			implParser.SetFilter(selectedEvents)
			active.Impls = append(active.Impls, implParser)
		}
	}

	return active, nil
}
