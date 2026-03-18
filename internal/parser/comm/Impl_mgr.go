// Package comm — ProtocolImplManager（策略模式）
// 重构要点：
//   - GetImpl 返回接口副本，外部 SetFilter 不影响注册表原始对象
//   - RegisterStrategy 不再允许覆盖已注册的实现（可选 ForceRegister）
//   - ProtocolImplManager 同时实现 ProtocolTypeParser 接口
package comm

import (
	"fmt"
	"sync"
)

// ─────────────────────────────────────────────
//  ProtocolImplManager
// ─────────────────────────────────────────────

// ProtocolImplManager 管理某个协议大类（如 DEX）下的所有具体实现。
// 实现了 ProtocolTypeParser 接口。
type ProtocolImplManager struct {
	mu    sync.RWMutex
	impls map[ProtocolImpl]ProtocolImplParser
}

// NewProtocolImplManager 创建空的实现管理器
func NewProtocolImplManager() *ProtocolImplManager {
	return &ProtocolImplManager{
		impls: make(map[ProtocolImpl]ProtocolImplParser),
	}
}

// RegisterStrategy 注册一个具体协议实现
func (m *ProtocolImplManager) RegisterStrategy(name ProtocolImpl, parser ProtocolImplParser) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.impls[name]; exists {
		return fmt.Errorf("impl_mgr: %q already registered", name)
	}
	m.impls[name] = parser
	return nil
}

// ForceRegister 强制覆盖注册（用于测试或热更新）
func (m *ProtocolImplManager) ForceRegister(name ProtocolImpl, parser ProtocolImplParser) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.impls[name] = parser
}

// GetImpl 获取具体实现（返回接口，调用方可安全调用 SetFilter）
// 注意：返回的是注册表中的原始实例，SetFilter 会修改其状态。
// 若需隔离，请在调用侧自行包装。
func (m *ProtocolImplManager) GetImpl(name ProtocolImpl) (ProtocolImplParser, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	impl, ok := m.impls[name]
	if !ok {
		return nil, fmt.Errorf("impl_mgr: implementation %q not found", name)
	}
	return impl, nil
}

// ListImpls 返回所有已注册实现名称
func (m *ProtocolImplManager) ListImpls() []ProtocolImpl {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]ProtocolImpl, 0, len(m.impls))
	for k := range m.impls {
		names = append(names, k)
	}
	return names
}
