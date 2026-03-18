package scanner

import "github.com/ETHSentinel/internal/parser"

// Engine 返回注入的解析引擎（供 sentinel 层调用 BuildActive）
func (s *Scanner) Engine() *parser.Engine {
	return s.engine
}
