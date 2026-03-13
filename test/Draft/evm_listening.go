//package main
//
//import (
//	"context"
//	"fmt"
//	"github.com/ethereum/go-ethereum/ethclient"
//	"github.com/ethereum/go-ethereum/rpc"
//	"log"
//	"net/http"
//	"net/url"
//)
//
//func main() {
//	// Infura URL（可加密钥）
//	rpcUrl := "https://mainnet.infura.io/v3/0d79a9c32c814e1da6133850f6fa1128"
//
//	// 代理地址，比如你本地代理端口7890
//	proxyStr := "http://127.0.0.1:7890"
//	proxyURL, err := url.Parse(proxyStr)
//	if err != nil {
//		log.Fatalf("代理地址解析失败: %v", err)
//	}
//	transport := &http.Transport{
//		Proxy: http.ProxyURL(proxyURL),
//	}
//	httpClient := &http.Client{
//		Transport: transport,
//	}
//	// 使用自定义 httpClient 建立 rpc 连接
//	rpcClient, err := rpc.DialHTTPWithClient(rpcUrl, httpClient)
//	if err != nil {
//		log.Fatalf("连接以太坊节点失败: %v", err)
//	}
//	client := ethclient.NewClient(rpcClient)
//	defer client.Close()
//
//	fmt.Println("连接成功，开始获取最新区块高度...")
//
//	header, err := client.HeaderByNumber(context.Background(), nil)
//
//	if err != nil {
//		log.Fatalf("获取最新区块失败: %v", err)
//	}
//
//	fmt.Printf("当前最新区块高度: %d, Hash: %s\n", header.Number.Uint64(), header.Hash().Hex())
//}
