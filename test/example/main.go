// test/example — SDK 使用示例
//
// 演示两种使用姿势：
//  1. 嵌入模式（直接连接以太坊节点，进程内解析）
//  2. 远程 gRPC 模式（连接独立部署的 Sentinel Server）
package main

import (
	"context"
	"fmt"
	"log"

	sentinel "github.com/ETHSentinel/client"
)

func main() {
	// ──────────────────────────────────────────────
	//  示例 1：嵌入模式
	// ──────────────────────────────────────────────
	embeddedExample()

	// ──────────────────────────────────────────────
	//  示例 2：远程 gRPC 模式
	// ──────────────────────────────────────────────
	// remoteExample()
}

// embeddedExample 直接连接节点，在进程内解析
func embeddedExample() {
	client, err := sentinel.New(
		sentinel.WithRPCURL("xxx"),
		sentinel.WithWSURL("xxx"),
		sentinel.WithProxy("http://127.0.0.1:7890"), // 可选代理
		sentinel.WithWorkerPoolSize(3),
	)
	if err != nil {
		log.Fatal("init sentinel:", err)
	}
	defer client.Close()

	ctx := context.Background()

	// ── 方式 A：扫描单个区块 ──────────────────────
	fmt.Println("=== 扫描单块 ===")
	result, err := sentinel.NewScanBuilder(client).
		WithDEX(sentinel.ProtocolImplUniswapV2, sentinel.EventMethodSwap).
		WithToken(sentinel.ProtocolImplERC20, sentinel.EventMethodTransfer).
		ScanOne(ctx, 22000000)
	if err != nil {
		log.Fatal("scan block:", err)
	}
	printResult(result)

	// ── 方式 B：流式扫描区间 ──────────────────────
	fmt.Println("\n=== 流式扫描区间 [22000000, 22000010) ===")
	ch, err := sentinel.NewScanBuilder(client).
		FromBlock(22000000).
		ToBlock(22000010).
		WithDEX(sentinel.ProtocolImplUniswapV2). // 不传 methods = 全部事件
		WithToken(sentinel.ProtocolImplERC20).
		Stream(ctx)
	if err != nil {
		log.Fatal("scan blocks:", err)
	}

	totalEvents := 0
	for res := range ch {
		totalEvents += len(res.Events)
		fmt.Printf("  block %s: %d txs, %d events\n",
			res.BlockNumber, res.TxCount, len(res.Events))
	}
	fmt.Printf("总计解析事件：%d\n", totalEvents)

	// ── 方式 C：实时订阅（生产使用） ──────────────
	// subCtx, cancel := context.WithCancel(ctx)
	// defer cancel()
	// sub, err := sentinel.NewScanBuilder(client).
	//     WithDEX(sentinel.ProtocolImplUniswapV2, sentinel.EventMethodSwap).
	//     Subscribe(subCtx)
	// for res := range sub {
	//     for _, ev := range res.Events {
	//         handleEvent(ev)
	//     }
	// }
}

// remoteExample 通过 gRPC 连接独立部署的 Sentinel Server
func remoteExample() {
	client, err := sentinel.NewRemote("localhost:50051",
		sentinel.WithWorkerPoolSize(20),
	)
	if err != nil {
		log.Fatal("init remote sentinel:", err)
	}
	defer client.Close()

	ctx := context.Background()

	ch, err := sentinel.NewScanBuilder(client).
		FromBlock(22000000).
		ToBlock(22001000).
		WithDEX(sentinel.ProtocolImplUniswapV2, sentinel.EventMethodSwap).
		Stream(ctx)
	if err != nil {
		log.Fatal("scan blocks:", err)
	}

	for res := range ch {
		printResult(res)
	}
}

func printResult(res *sentinel.ScanResult) {
	fmt.Printf("Block #%s | txs=%d | events=%d\n",
		res.BlockNumber, res.TxCount, len(res.Events))
	for _, ev := range res.Events {
		fmt.Printf("  [%s/%s] %s tx=%s\n",
			ev.GetProtocolType(),
			ev.GetProtocolImpl(),
			ev.GetEventType(),
			ev.GetTxHash().Hex()[:12]+"…",
		)
		// 访问 Swap 详情
		if swap, ok := ev.GetDetail().(*sentinel.SwapData); ok {
			fmt.Printf("    Swap: %s → %s  amount=%s\n",
				swap.FromToken.Hex()[:8]+"…",
				swap.ToToken.Hex()[:8]+"…",
				swap.FromAmount.String(),
			)
		}
	}
}
