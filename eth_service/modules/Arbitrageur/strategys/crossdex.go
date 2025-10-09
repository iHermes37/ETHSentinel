package arbitragur

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Crypto-ChainSentinel/initialize"
	"github.com/Crypto-ChainSentinel/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/machinebox/graphql"
)

type CrossDEXStrategy struct{}

func (c *CrossDEXStrategy) Name() string {
	return "CrossDEX"
}

// getReserves 获取对应交易对的实时储备
func GetReserves(p *models.CrossPairData) error {
	client := initialize.InfuraConn_ws() // 返回 *ethclient.Client

	pairAddress := p.Pair_DexA.PoolAddr
	//初始化交易对实例
	pairContract, err := uniswapv2.NewUniswapV2Pair(pairAddress, client)
	if err != nil {
		return fmt.Errorf("failed to create pair contract: %w", err)
	}

	reserves, err := pairContract.GetReserves(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		return fmt.Errorf("failed to get reserves: %w", err)
	}

	// 确认 token0/token1 对应
	//var reserveA, reserveB *big.Int
	//if p.Pair_DexA.Token0 == p.Token0 {
	//	reserveA = reserves.Reserve0
	//	reserveB = reserves.Reserve1
	//} else {
	//	reserveA = reserves.Reserve1
	//	reserveB = reserves.Reserve0
	//}
	p.Pair_DexA.PairReserve.ReserveA0 = reserves.Reserve0
	p.Pair_DexA.PairReserve.ReserveA1 = reserves.Reserve1

}

func (c *CrossDEXStrategy) Run(stop <-chan struct{}) {
	// 套利机会通道
	arbCh := make(chan models.CrossPairData, 100)
	// ----------------- 套利执行 worker -----------------
	go func() {
		for pair := range arbCh {
			Executearbitrage(pair) // 立即执行套利
		}
	}()
	//----每个发现的套利机会 都会立即通过 Flashbots 原子提交----------
	//go func() {
	//	for pair := range arbCh {
	//		// 不直接执行链上交易，而是通过 Flashbots 提交
	//		err := SubmitFlashbotsBundle(pair)
	//		if err != nil {
	//			log.Println("Flashbots submit failed:", err)
	//		}
	//	}
	//}()
	// ----------------- 获取交易对 --------------------------------
	u_V3 := models.Uniswap_V3{
		GraphQLDEX: models.GraphQLDEX{
			Client:    graphql.NewClient("https://gateway.thegraph.com/api/subgraphs/id/5zvR82QoaXYFyDEKLZ9t6v9adgnptxYpKpSbxtgVENFV"),
			AuthToken: "df5d393ba8219b65e3eea66df2242e6b",
			QueryString: `
            query($first: Int!, $skip: Int!) {
                pools(first: $first, skip: $skip) {
                    id
                    token0 { id,symbol }
                    token1 { id,symbol }
                    feeTier
                }
            }
        `,
		},
	}
	var resp_univ3 models.Pairs
	err1 := u_V3.GetPairs(map[string]interface{}{
		"first": 100,
		"skip":  0,
	}, &resp_univ3)
	if err1 != nil {
		log.Fatal(err1)
	}

	var resp_sushi models.Pairs
	s := models.Sushiswap{
		GraphQLDEX: models.GraphQLDEX{
			Client:    graphql.NewClient("https://gateway.thegraph.com/api/subgraphs/id/A4JrrMwrEXsYNAiYw7rWwbHhQZdj6YZg1uVy5wa6g821"),
			AuthToken: "df5d393ba8219b65e3eea66df2242e6b",
			QueryString: `
            query($first: Int!, $skip: Int!) {
                pools(first: $first, skip: $skip) {
                    id
                    token0 { id,symbol }
                    token1 { id,symbol }
                }
            }
        `,
		}}
	err2 := s.GetPairs(map[string]interface{}{
		"first": 100,
		"skip":  0,
	}, &resp_sushi)
	if err2 != nil {
		log.Fatal(err2)
	}

	//----------------------------------------------------------------
	commonpairs := GetCommonPairs(resp_univ3, resp_sushi)
	for _, pair := range commonpairs {
		go func(p models.CrossPairData) {
			ticker := time.NewTicker(200 * time.Millisecond) // 可改为链上事件推送
			defer ticker.Stop()
			//循环等待事件到来
			for {
				select {
				case <-stop:
					fmt.Println("CrossDEX stopped")
					return
				case <-ticker.C:
					fmt.Println("CrossDEX running tick")
					// -------核心套利逻辑------------

					// ---------- 核心套利逻辑 ----------
					// 更新储备信息
					err := GetReserves(&p)
					if err != nil {
						fmt.Println("xx")
						continue
					}
					// 决定套利方向
					DecidePullToken(&p)
					// 计算套利机会
					p.Opportunity = CalBN(&p)
					// 校验并记录结果
					ValidateAndRecord(&p)
					// 发现套利机会 → 立即发送 channel 执行
					if p.Opportunity.Profit > thresholdProfit {
						arbCh <- p
					}
				}
			}
		}(pair)
	}
	// 阻塞等待 stop 信号
	<-stop
	close(arbCh)
}

// 执行套利合约调用
func Executearbitrage(pair models.CrossPairData) {
	// 1. 初始化客户端
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_INFURA_KEY")
	if err != nil {
		log.Fatal(err)
	}

	// 2. 合约地址
	contractAddr := common.HexToAddress("0xYourCrocssdexAddress")
	instance, err := crocssdex.NewCrocssdex(contractAddr, client)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 创建交易授权
	auth, err := bind.NewTransactorWithChainID(yourKeyReader, yourPrivateKey, big.NewInt(1)) // 1 = Ethereum mainnet
	if err != nil {
		log.Fatal(err)
	}

	// 4. 确定借贷资产和数量
	var asset common.Address
	var amount big.Int

	// 根据 pullToken 决定基准代币
	if pair.Direction == 1 || pair.Direction == 2 {
		asset = pair.PullToken.Id // 基准代币
		amount = pair.Opportunity.X
	} else {
		fmt.Println("No Arbitrageur opportunity")
		return
	}

	// 5. 调用闪电贷 + 套利
	tx, err := instance.RequestFlashLoan(auth, asset, amount)
	if err != nil {
		log.Fatal("Flashloan execution failed:", err)
	}

	fmt.Println("Arbitrage tx sent:", tx.Hash().Hex())
}

func SubmitFlashbotsBundle(pair models.CrossPairData) error {
	// 1. 构建交易数据，调用套利合约
	tx := &flashbots.Transaction{
		To:    common.HexToAddress(ARBITRAGE_CONTRACT), //指定了交易 目标合约地址，就是你部署的 闪电贷 + 跨 DEX 套利合约 地址
		Data:  EncodeExecuteArbitrage(pair),            // 合约函数的 ABI 编码
		Gas:   500000,
		Value: big.NewInt(0),
	}

	// 2. 获取当前区块号
	blockNumber, err := ethClient.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	// 3. 构建 bundle
	bundle := []*flashbots.Transaction{tx}

	// 4. 发送 bundle 给 Flashbots relay
	resp, err := fbProvider.SendBundle(context.Background(), bundle, blockNumber+1)
	if err != nil {
		return err
	}

	fmt.Println("Bundle submitted, hash:", resp.BundleHash.Hex())
	return nil
}

// DecidePullToken 决定套利方向和基准代币
func DecidePullToken(p *models.CrossPairData) {
	// 计算两池价格：price = token1 / token0
	priceA := new(big.Float).Quo(
		new(big.Float).SetInt(p.Pair_DexA.PairReserve.Reserve1),
		new(big.Float).SetInt(p.Pair_DexA.PairReserve.Reserve0),
	)
	priceB := new(big.Float).Quo(
		new(big.Float).SetInt(p.Pair_DexB.PairReserve.Reserve1),
		new(big.Float).SetInt(p.Pair_DexB.PairReserve.Reserve0),
	)

	cmp := priceA.Cmp(priceB)

	switch {
	case cmp < 0:
		// priceA < priceB → ETH 在 A 更便宜
		p.Direction = 1
		p.PullToken = p.Pair_DexA.Token0 // 买 Token0（基准代币）
	case cmp > 0:
		// priceA > priceB → ETH 在 B 更便宜，等价于 USDC 在 A 更便宜
		p.Direction = 2
		p.PullToken = p.Pair_DexA.Token1 // 买 Token1（基准代币）
	default:
		p.Direction = 0
		p.PullToken = models.Token{}
	}
}

// 计算套利利润
func CalBN(p *models.CrossPairData) models.ArbitrageOpportunity {
	var A0, A1, B0, B1 *big.Int
	// 根据 PullToken 选择基准代币
	if p.PullToken == p.Pair_DexA.Token0 {
		// 基准代币是 Token0
		A0, A1 = p.Pair_DexA.PairReserve.Reserve0, p.Pair_DexA.PairReserve.Reserve1
		B0, B1 = p.Pair_DexB.PairReserve.Reserve0, p.Pair_DexB.PairReserve.Reserve1
	} else if p.PullToken == p.Pair_DexA.Token1 {
		// 基准代币是 Token1
		A0, A1 = p.Pair_DexA.PairReserve.Reserve1, p.Pair_DexA.PairReserve.Reserve0
		B0, B1 = p.Pair_DexB.PairReserve.Reserve1, p.Pair_DexB.PairReserve.Reserve0
	} else {
		return models.ArbitrageOpportunity{Profit: big.NewFloat(0)}
	}

	// 将 big.Int 转 big.Float 做浮点运算
	fA0 := new(big.Float).SetInt(A0)
	fA1 := new(big.Float).SetInt(A1)
	fB0 := new(big.Float).SetInt(B0)
	fB1 := new(big.Float).SetInt(B1)

	// 假设投入量 x = min(A1, B1)/2 （这里可以改为优化公式）
	fX := new(big.Float).Quo(fA1.Add(fA1, fB1), big.NewFloat(2))

	// dy = x * 0.997 * A0 / (A1 + x * 0.997)
	tmp := new(big.Float).Mul(fX, big.NewFloat(0.997))
	dy := new(big.Float).Quo(new(big.Float).Mul(tmp, fA0), new(big.Float).Add(fA1, tmp))

	// dz = x * B0 / (B1 - x)
	dz := new(big.Float).Quo(new(big.Float).Mul(fX, fB0), new(big.Float).Sub(fB1, fX))

	profit := new(big.Float).Sub(dy, dz)

	// 转回 big.Int
	xInt, _ := fX.Int(nil)
	yInt, _ := dy.Int(nil)
	zInt, _ := dz.Int(nil)
	profitInt, _ := profit.Int(nil)

	return models.ArbitrageOpportunity{
		X:      xInt,
		Y:      yInt,
		Z:      zInt,
		Profit: profitInt,
	}
}
