//package main
//
//import (
//	"context"
//	"fmt"
//	"log"
//	"os"
//
//	"github.com/gagliardetto/solana-go/rpc"
//)
//
//func main() {
//
//	// 设置环境变量（可选，如果代理已经全局配置）
//	os.Setenv("HTTP_PROXY", "http://127.0.0.1:7890")
//	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:7890")
//
//	client := rpc.New(rpc.MainNetBeta_RPC)
//
//	height, err := client.GetBlockHeight(
//		context.Background(),
//		rpc.CommitmentFinalized,
//	)
//	if err != nil {
//		log.Fatalf("GetBlockHeight failed: %v", err)
//	}
//
//	fmt.Printf("当前区块高度: %d\n", height)
//}
