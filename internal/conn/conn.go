// Package conn 管理以太坊节点的 RPC / WebSocket 连接。
// 重构要点：
//   - 去掉全局状态，改用依赖注入
//   - 增加连接池支持（多路复用）
//   - 错误明确返回，不 panic
package conn

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
)

// ─────────────────────────────────────────────
//  配置
// ─────────────────────────────────────────────

// NodeConfig 节点连接配置
type NodeConfig struct {
	Name     string // 节点名称，例如 "Infura"
	RPCURL   string `yaml:"RPC"`
	WSURL    string `yaml:"WS"`
	ProxyURL string `yaml:"proxy"` // 可选：SOCKS5 / HTTP 代理
}

// Method 连接方式
type Method string

const (
	MethodRPC Method = "RPC"
	MethodWS  Method = "WS"
)

// ─────────────────────────────────────────────
//  Manager（核心：替换原 ConnMgr + 全局函数）
// ─────────────────────────────────────────────

// Manager 管理多个节点配置，提供线程安全的连接获取
type Manager struct {
	mu      sync.RWMutex
	nodes   map[string]*NodeConfig
	pool    map[string]*ethclient.Client // 简单单例池，可替换为真正连接池
	logger  *zap.Logger
}

// NewManager 创建连接管理器
func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		nodes:  make(map[string]*NodeConfig),
		pool:   make(map[string]*ethclient.Client),
		logger: logger,
	}
}

// Register 注册节点配置
func (m *Manager) Register(cfg *NodeConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nodes[cfg.Name] = cfg
}

// Get 获取（或创建）一个以太坊客户端连接
// key 格式："{NodeName}:{Method}"，例如 "Infura:WS"
func (m *Manager) Get(ctx context.Context, name string, method Method) (*ethclient.Client, error) {
	key := fmt.Sprintf("%s:%s", name, method)

	m.mu.RLock()
	if client, ok := m.pool[key]; ok {
		m.mu.RUnlock()
		return client, nil
	}
	m.mu.RUnlock()

	// 创建新连接
	m.mu.Lock()
	defer m.mu.Unlock()

	// double-check
	if client, ok := m.pool[key]; ok {
		return client, nil
	}

	cfg, ok := m.nodes[name]
	if !ok {
		return nil, fmt.Errorf("conn: node %q not registered", name)
	}

	client, err := dial(ctx, cfg, method)
	if err != nil {
		return nil, fmt.Errorf("conn: dial %s/%s: %w", name, method, err)
	}

	m.pool[key] = client
	m.logger.Info("eth node connected", zap.String("node", name), zap.String("method", string(method)))
	return client, nil
}

// Close 关闭所有连接
func (m *Manager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, c := range m.pool {
		c.Close()
		delete(m.pool, k)
	}
}

// ─────────────────────────────────────────────
//  底层 Dial 实现
// ─────────────────────────────────────────────

func dial(ctx context.Context, cfg *NodeConfig, method Method) (*ethclient.Client, error) {
	switch method {
	case MethodRPC:
		return dialRPC(cfg)
	case MethodWS:
		return dialWS(ctx, cfg)
	default:
		return nil, fmt.Errorf("conn: unknown method %q", method)
	}
}

func dialRPC(cfg *NodeConfig) (*ethclient.Client, error) {
	if cfg.ProxyURL == "" {
		c, err := ethclient.Dial(cfg.RPCURL)
		if err != nil {
			return nil, err
		}
		return c, nil
	}

	proxyURL, err := url.Parse(cfg.ProxyURL)
	if err != nil {
		return nil, fmt.Errorf("conn: parse proxy url: %w", err)
	}
	transport := &http.Transport{
		Proxy:               http.ProxyURL(proxyURL),
		DialContext:         (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	httpClient := &http.Client{Transport: transport, Timeout: 30 * time.Second}

	rpcClient, err := rpc.DialHTTPWithClient(cfg.RPCURL, httpClient)
	if err != nil {
		return nil, err
	}
	return ethclient.NewClient(rpcClient), nil
}

func dialWS(ctx context.Context, cfg *NodeConfig) (*ethclient.Client, error) {
	wsDialer := websocket.DefaultDialer

	if cfg.ProxyURL != "" {
		socksDialer, err := proxy.SOCKS5("tcp", cfg.ProxyURL, nil, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("conn: create socks5 dialer: %w", err)
		}
		wsDialer = &websocket.Dialer{
			NetDial: func(network, addr string) (net.Conn, error) {
				return socksDialer.Dial(network, addr)
			},
			HandshakeTimeout: 10 * time.Second,
		}
	}

	rpcClient, err := rpc.DialWebsocketWithDialer(ctx, cfg.WSURL, "", *wsDialer)
	if err != nil {
		return nil, err
	}
	return ethclient.NewClient(rpcClient), nil
}
