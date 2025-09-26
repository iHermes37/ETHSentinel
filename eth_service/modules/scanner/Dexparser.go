package scanner

import (
	"context"
	"fmt"
	"github.com/CryptoQuantX/chain_monitor/initialize"
	"github.com/CryptoQuantX/chain_monitor/models"
	"github.com/CryptoQuantX/chain_monitor/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
	"strings"
)

func ParseDexInput(tx *types.Transaction, dexname string) models.DeFiTxDetail {
	cli := initialize.InfuraConn(rpcURL, proxyStr)
	if dexname == "" {
		return ParseUniswapV2(*tx, cli)
	}

	return models.DeFiTxDetail{}
}

// 获取 pair 的 token0 / token1
func getPairTokens(pair common.Address, cli *ethclient.Client) (common.Address, common.Address) {
	pairAbi, _ := abi.JSON(strings.NewReader(pairABI))
	call := ethereum.CallMsg{To: &pair}
	res0, _ := cli.CallContract(context.Background(), call, nil)
	out0, _ := pairAbi.Unpack("token0", res0)

	res1, _ := cli.CallContract(context.Background(), call, nil)
	out1, _ := pairAbi.Unpack("token1", res1)

	return out0[0].(common.Address), out1[0].(common.Address)
}

func ParseUniswapV2(tx types.Transaction, cli *ethclient.Client) models.DeFiTxDetail {
	var d = models.DeFiTxDetail{}
	data := tx.Data()
	if len(data) < 4 {
		log.Fatalf("	1")
	}

	routerAbi, err := utils.ReadABIFile("chainTxTracker/eth_service/config/abi/UniswapV2_Router.json")
	if err != nil {
		log.Fatalf("加载 ABI 失败: %v", err)
	}

	method, err := routerAbi.MethodById(data[:4])
	if err != nil {
		log.Fatalf("")
	}

	// 解析参数
	params := make(map[string]interface{})
	err = method.Inputs.UnpackIntoMap(params, data[4:])
	if err != nil {
		fmt.Println("参数解析失败:", err)
	}

	if method.Name == "swapExactTokensForTokens" {

		if path, ok := params["path"].([]common.Address); ok {
			fmt.Println("兑换路径:")
			for _, token := range path {
				symbol, decimals := getERC20Meta(token, cli)
				fmt.Printf(" - %s (%s), decimals=%d\n", token.Hex(), symbol, decimals)
			}
		}
		// 拿到交易回执，解析 Swap 日志
		receipt, err := cli.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			fmt.Println("获取交易回执失败:", err)
		}

		pairAbi, _ := abi.JSON(strings.NewReader(pairABI))
		for _, vLog := range receipt.Logs {
			// 遍历日志，解码 Swap 事件
			event := struct {
				Sender     common.Address
				Amount0In  *big.Int
				Amount1In  *big.Int
				Amount0Out *big.Int
				Amount1Out *big.Int
				To         common.Address
			}{}

			// 解码 Swap 事件
			err := pairAbi.UnpackIntoInterface(&event, "Swap", vLog.Data)
			if err == nil {
				// Pair 合约地址就是 token0/token1 的来源
				pairAddress := vLog.Address
				token0, token1 := getPairTokens(pairAddress, cli)

				// 获取 token 元信息
				symbol0, decimals0 := getERC20Meta(token0, cli)
				symbol1, decimals1 := getERC20Meta(token1, cli)

				// 计算 In / Out
				if event.Amount0In.Cmp(big.NewInt(0)) > 0 {
					d.TokensInsymbol = append(d.TokensInsymbol, symbol0)
					amt := new(big.Float).Quo(new(big.Float).SetInt(event.Amount0In), big.NewFloat(math.Pow10(int(decimals0))))
					val, _ := amt.Float64()
					d.AmountsIn = append(d.AmountsIn, val)
				}
				if event.Amount1In.Cmp(big.NewInt(0)) > 0 {
					d.TokensInsymbol = append(d.TokensInsymbol, symbol1)
					amt := new(big.Float).Quo(new(big.Float).SetInt(event.Amount1In), big.NewFloat(math.Pow10(int(decimals1))))
					val, _ := amt.Float64()
					d.AmountsIn = append(d.AmountsIn, val)
				}
				if event.Amount0Out.Cmp(big.NewInt(0)) > 0 {
					d.TokensOutsymbol = append(d.TokensOutsymbol, symbol0)
					amt := new(big.Float).Quo(new(big.Float).SetInt(event.Amount0Out), big.NewFloat(math.Pow10(int(decimals0))))
					val, _ := amt.Float64()
					d.AmountsOut = append(d.AmountsOut, val)
				}
				if event.Amount1Out.Cmp(big.NewInt(0)) > 0 {
					d.TokensOutsymbol = append(d.TokensOutsymbol, symbol1)
					amt := new(big.Float).Quo(new(big.Float).SetInt(event.Amount1Out), big.NewFloat(math.Pow10(int(decimals1))))
					val, _ := amt.Float64()
					d.AmountsOut = append(d.AmountsOut, val)
				}
			}
		}

		d.Exchange = "UniswapV2Protocol"
		d.Direction = models.TxDirectionOut // 假设是用户主动 swap
	}
	return d
}
