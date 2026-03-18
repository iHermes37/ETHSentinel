// cmd/example — ETH Sentinel SDK 完整使用示例
// 演示：多链扫描 + Mempool 监控 + 钱包功能
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	sentinel "github.com/ETHSentinel/client"
)

func main() {
	// ── 示例 1：多链扫描 ──────────────────────────
	fmt.Println("=== 示例 1：ETH 主网扫描 ===")
	scanExample()

	// ── 示例 2：Mempool 监控 ──────────────────────
	fmt.Println("\n=== 示例 2：Mempool 监控（5秒） ===")
	mempoolExample()

	// ── 示例 3：钱包功能 ──────────────────────────
	fmt.Println("\n=== 示例 3：钱包 ===")
	walletExample()

	// ── 示例 4：BSC 链扫描 ───────────────────────
	fmt.Println("\n=== 示例 4：BSC 扫描 ===")
	bscExample()
}

// scanExample ETH 主网区块扫描
func scanExample() {
	client, err := sentinel.New(
		sentinel.WithChainID(uint64(sentinel.ChainETH)),
		sentinel.WithRPCURL("https://mainnet.infura.io/v3/YOUR_KEY"),
		sentinel.WithWorkerPoolSize(5),
	)
	if err != nil {
		log.Fatal("init sdk:", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 扫描单块
	result, err := sentinel.NewScanBuilder(client).
		WithDEX(sentinel.ProtocolImplUniswapV2, sentinel.EventMethodSwap).
		WithToken(sentinel.ProtocolImplERC20, sentinel.EventMethodTransfer).
		ScanOne(ctx, 22000000)
	if err != nil {
		log.Println("scan error:", err)
		return
	}

	fmt.Printf("Block #%s | chain=%d | txs=%d | events=%d\n",
		result.BlockNumber, result.ChainID, result.TxCount, len(result.Events))

	for _, ev := range result.Events {
		fmt.Printf("  [%s/%s] %s\n",
			ev.GetProtocolType(), ev.GetProtocolImpl(), ev.GetEventType())
		if swap, ok := ev.GetDetail().(*sentinel.SwapData); ok {
			fmt.Printf("    Swap: %s → %s amount=%s\n",
				swap.FromToken.Hex()[:10], swap.ToToken.Hex()[:10], swap.FromAmount)
		}
	}
}

// mempoolExample 监控 Mempool 中的 pending 交易
func mempoolExample() {
	client, err := sentinel.New(
		sentinel.WithChainID(uint64(sentinel.ChainETH)),
		sentinel.WithWSURL("wss://mainnet.infura.io/ws/v3/YOUR_KEY"),
	)
	if err != nil {
		log.Fatal("init sdk:", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mempoolClient := client.Mempool()

	// 监控大额 UniswapV2 swap（>1 ETH value，GasPrice > 10 Gwei）
	pendingCh, err := mempoolClient.Subscribe(ctx,
		sentinel.FilterByMinValueETH(1.0),
		sentinel.FilterByMinGas(10),
		sentinel.FilterByMethod("0x38ed1739"), // swapExactTokensForTokens
	)
	if err != nil {
		log.Println("mempool subscribe error:", err)
		return
	}

	count := 0
	for ptx := range pendingCh {
		fmt.Printf("  pending: %s from=%s value=%s\n",
			ptx.Tx.Hash().Hex()[:12],
			ptx.From.Hex()[:10],
			ptx.Tx.Value().String(),
		)
		count++
		if count >= 5 {
			break
		}
	}
	fmt.Printf("  监控到 %d 笔 pending 交易\n", count)
}

// walletExample 钱包功能演示
func walletExample() {
	// 生成新助记词
	mnemonic, err := sentinel.GenerateMnemonic()
	if err != nil {
		log.Fatal("generate mnemonic:", err)
	}
	fmt.Printf("  新助记词: %s\n", mnemonic[:20]+"...")

	client, err := sentinel.New(
		sentinel.WithChainID(uint64(sentinel.ChainETH)),
		sentinel.WithRPCURL("https://mainnet.infura.io/v3/YOUR_KEY"),
		sentinel.WithMnemonic(mnemonic),
	)
	if err != nil {
		log.Fatal("init sdk:", err)
	}
	defer client.Close()

	walletClient, err := client.Wallet()
	if err != nil {
		log.Fatal("get wallet:", err)
	}

	// 派生前 3 个账户地址
	for i := uint32(0); i < 3; i++ {
		addr, err := walletClient.Address(i)
		if err != nil {
			log.Println("derive address error:", err)
			continue
		}
		fmt.Printf("  账户 #%d: %s\n", i, addr.Hex())
	}

	// 签名消息
	sig, err := walletClient.SignMessage(0, []byte("Hello ETH Sentinel"))
	if err != nil {
		log.Println("sign error:", err)
		return
	}
	fmt.Printf("  签名: 0x%x...\n", sig[:8])
}

// bscExample BSC 链扫描演示
func bscExample() {
	client, err := sentinel.New(
		sentinel.WithChainID(uint64(sentinel.ChainBSC)),
		// BSC 使用公共节点，无需 API Key
		sentinel.WithRPCURL("https://bsc-dataseed.binance.org"),
		sentinel.WithWorkerPoolSize(3),
	)
	if err != nil {
		log.Fatal("init bsc sdk:", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 流式扫描 BSC 区间
	ch, err := sentinel.NewScanBuilder(client).
		FromBlockBig(big.NewInt(40000000)).
		ToBlockBig(big.NewInt(40000003)).
		WithToken(sentinel.ProtocolImplERC20, sentinel.EventMethodTransfer).
		Stream(ctx)
	if err != nil {
		log.Println("scan error:", err)
		return
	}

	for res := range ch {
		fmt.Printf("  BSC Block #%s | chain=%d | txs=%d | events=%d\n",
			res.BlockNumber, res.ChainID, res.TxCount, len(res.Events))
	}
}
