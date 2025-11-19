package Conn

import (
	"golang.org/x/net/proxy"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
)

func RPCConn(rpc_url string, proxy_url string) *ethclient.Client {
	//rpcURL := "https://mainnet.infura.io/v3/0d79a9c32c814e1da6133850f6fa1128"
	//proxyStr := "http://192.168.182.215:7890"
	proxyURL, err := url.Parse(proxy_url)
	if err != nil {
		log.Fatalf("代理地址解析失败: %v", err)

		panic("数据库连接失败: " + err.Error())

	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	httpClient := &http.Client{
		Transport: transport,
	}
	// 使用自定义 httpClient 建立 rpc 连接
	rpcClient, err := rpc.DialHTTPWithClient(rpc_url, httpClient)
	if err != nil {
		log.Fatalf("连接以太坊节点失败: %v", err)
	}
	client := ethclient.NewClient(rpcClient)
	//defer client.Close()

	return client
}

func WSConn(ws_url string, proxy_url string) *ethclient.Client {
	socksDialer, _ := proxy.SOCKS5("tcp", proxy_url, nil, proxy.Direct)

	// 创建 websocket.Dialer 并设置 NetDial
	wsDialer := websocket.Dialer{
		NetDial: func(network, addr string) (net.Conn, error) {
			return socksDialer.Dial(network, addr)
		},
	}

	// 直接用 rpc.DialWebsocketWithDialer + *websocket.Dialer
	rpcClient, err := rpc.DialWebsocketWithDialer(context.Background(), ws_url, "", wsDialer)
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	client := ethclient.NewClient(rpcClient)
	log.Println("连接成功:", client)

	return client
}

// Infura 配置结构体
type Infura struct {
	RPCURL string `yaml:"RPC"`
	WSURL  string `yaml:"WS"`
	Proxy  string
}

// Config 整个配置文件的根结构
type Config struct {
	Infura Infura `yaml:"Infura"`
	Proxy  string `yaml:"proxy"`
}

func InitConnConfig(file_path string) *Infura {
	// 读取文件内容
	data, err := ioutil.ReadFile(file_path)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析 YAML
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("解析 YAML 失败: %v", err)
	}

	// 创建 Infura 实例并设置 ProxyURL
	infura := &Infura{
		RPCURL: config.Infura.RPCURL,
		WSURL:  config.Infura.WSURL,
		Proxy:  config.Proxy,
	}

	return infura
}

// ==================================
type ConnMethod string

const (
	RPC ConnMethod = "RPC"
	WS  ConnMethod = "WS"
)

type Conn interface {
	Connect(method ConnMethod) *ethclient.Client
}

func NewInfura() *Infura {
	conn_config_path := "EthereumSentinel/config/conn.yaml"
	return InitConnConfig(conn_config_path)
}

func (infura *Infura) Connect(method ConnMethod) *ethclient.Client {
	if method == RPC {
		return RPCConn(infura.RPCURL, infura.Proxy)
	} else {
		return WSConn(infura.WSURL, infura.Proxy)
	}
}

//===========================

type ConnMgr struct{}

func (f *ConnMgr) SelectNode(node string) *Infura {
	switch node {
	case "Infura":
		return NewInfura()
	default:
		return nil
	}
}
