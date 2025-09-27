package connectionManager

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"net/url"
)

func InfuraConn() *ethclient.Client {
	rpcURL := "https://mainnet.infura.io/v3/0d79a9c32c814e1da6133850f6fa1128"
	proxyStr := "http://192.168.150.215:7890"
	proxyURL, err := url.Parse(proxyStr)
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
	rpcClient, err := rpc.DialHTTPWithClient(rpcURL, httpClient)
	if err != nil {
		log.Fatalf("连接以太坊节点失败: %v", err)
	}
	client := ethclient.NewClient(rpcClient)
	defer client.Close()

	return client
}

func InfuraConn_ws() *ethclient.Client {
	socksDialer, _ := proxy.SOCKS5("tcp", "192.168.248.215:7890", nil, proxy.Direct)

	// 创建 websocket.Dialer 并设置 NetDial
	wsDialer := websocket.Dialer{
		NetDial: func(network, addr string) (net.Conn, error) {
			return socksDialer.Dial(network, addr)
		},
	}

	wsURL := "wss://mainnet.infura.io/ws/v3/0d79a9c32c814e1da6133850f6fa1128"

	// 直接用 rpc.DialWebsocketWithDialer + *websocket.Dialer
	rpcClient, err := rpc.DialWebsocketWithDialer(context.Background(), wsURL, "", wsDialer)
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	client := ethclient.NewClient(rpcClient)
	log.Println("连接成功:", client)

	return client
}
