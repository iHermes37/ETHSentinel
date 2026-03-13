package arbitrageur

import (
	"fmt"
	Conn "github.com/Crypto-ChainSentinel/internal"
	"github.com/Crypto-ChainSentinel/internal/core/arbitrageur/mempool"
	"github.com/Crypto-ChainSentinel/internal/lib/dex"
	"github.com/Crypto-ChainSentinel/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"strings"
)

func normalizePair(p types.Pair) (string, string) {
	if strings.ToLower(p.Token0.Symbol) < strings.ToLower(p.Token1.Symbol) {
		return strings.ToLower(p.Token0.Symbol), strings.ToLower(p.Token1.Symbol)
	} else {
		return strings.ToLower(p.Token1.Symbol), strings.ToLower(p.Token0.Symbol)
	}
}

type Monitor struct {
	*mempool.Radar
	Client    *ethclient.Client
	DexAapter *dex.DexAdapter
}

func NewMonitor() *Monitor {
	connMgr := Conn.ConnMgr{}
	Infura := connMgr.SelectNode("Infura")
	client := Infura.Connect(Conn.WS)
	return &Monitor{
		Client: client,
	}
}

// 获取交易对的实时储备
func (m *Monitor) GetPairReserves(
	d *types.DEXProtocol,
	pairAddr common.Address,
) (types.DexPairReserver, error) {

	// 初始化交易对实例
	pairContract, err := m.DexAapter.SelectDexPair(d, pairAddr)
	if err != nil {
		return types.DexPairReserver{}, err
	}

	// 获取储备量
	res, err := pairContract.GetReserves(&bind.CallOpts{})
	if err != nil {
		return types.DexPairReserver{}, err
	}

	// 返回结果封装为你的类型
	return types.DexPairReserver{
		Reserve0: res.Reserve0,
		Reserve1: res.Reserve1,
	}, nil
}

// 监控指定 DEX 新创建的交易对
func (m *Monitor) MonitoringNewPair(d *types.DEXProtocol) {
	// 选择工厂地址
	factoryAddr := m.DexAapter.GetFactoryAddress(d)

	// 创建 Factory 事件 Filterer
	factory, err := uniswap_v2.NewUniswapv2FactoryFilterer(factoryAddr, m.Client)
	if err != nil {
		log.Fatalf("factory binding error: %v", err)
	}

	// 订阅 PairCreated 事件
	logs := make(chan *uniswap_v2.Uniswapv2FactoryPairCreated)
	sub, err := factory.WatchPairCreated(
		&bind.WatchOpts{Context: context.Background()},
		logs,
		nil, // token0 (indexed)
		nil, // token1 (indexed)
	)
	if err != nil {
		log.Fatalf("subscribe error: %v", err)
	}

	log.Println("Start monitoring new pairs on", d.String())

	// 事件循环
	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Printf("subscription error: %v", err)
				// TODO: 你可以在此处自动重连，确保系统长期稳定运行
				return

			case evt := <-logs:
				m.handleNewPairEvent(d, evt)
			}
		}
	}()
}

// 处理新创建的代币池
func (m *Monitor) handleNewPairEvent(d *types.DEXProtocol, evt *uniswap_v2.Uniswapv2FactoryPairCreated) {
	fmt.Println("=== New Pair Found ===")
	fmt.Println("Token0:", evt.Token0.Hex())
	fmt.Println("Token1:", evt.Token1.Hex())
	fmt.Println("Pair  :", evt.Pair.Hex())

	// 查询储备量（一般刚创建是 0）
	reserves, err := m.GetPairReserves(d, evt.Pair)
	if err != nil {
		log.Println("GetReserves error:", err)
		return
	}

	fmt.Println("Reserves:", reserves.Reserve0, reserves.Reserve1)

	// TODO：这里你可以加入自己的逻辑，例如：
	// 1. 判断是否为 honeypot（结合 router 校验是否可卖）
	// 2. 自动监听该池子 swap/sync
	// 3. 自动加入监控池列表
	// 4. 推送 Telegram / Discord 报警
	// 5. 与你的套利引擎对接
}
