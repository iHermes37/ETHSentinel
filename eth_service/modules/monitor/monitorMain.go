package monitor

import (
	"fmt"
	"github.com/CryptoQuantX/chain_monitor/models"
	handle "github.com/CryptoQuantX/chain_monitor/modules/handler"
	parser "github.com/CryptoQuantX/chain_monitor/modules/parser"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func Monitor(recvblock *models.BlockStruct, cli *ethclient.Client) {

	threshold := new(big.Int)
	threshold.SetString("50000000000000000000", 10) // 50 ETH (单位: wei)
	BlockNumber := recvblock.Header.BlockNumber

	go WatchTokenHold()

	for _, tx := range recvblock.Transactions {
		from := parser.Parsefrom(cli, tx)
		to := tx.To()
		nonce := tx.Nonce()

		//------------新部署的合约----------------------------
		if to == nil {
			//var toType = "合约部署"
			contractaddr := parser.GetNewContractAddr(from, nonce)
			if erc, error := parser.DetectERCStandard(tx); error == nil {
				fmt.Printf("检测到标准: %s, 方法: %s, 参数: %+v\n", erc.Ercname, erc.Opmethod, erc.Params)
				handle.HandleNewContract(tx, contractaddr, erc, BlockNumber)
			}
		}
		//--------------巨鲸追踪-----------------------------
		if parser.IsWhale(from) {
			//----------------合约/Defi交互----------------------
			// 先判断是否是 DEX Router / DeFi 协议交互
			if ok, definame := parser.IsDEXRouter(to); ok {
				// 处理 DeFi 交互，例如套利或 LP Token 相关操作
				handle.HandleDeFi(tx, from, definame)
			} else {
				// 再判断是否是代币转账
				msg, err := parser.DetectERCStandard(tx)
				if err == nil {
					if msg.Opmethod == "transfer" || msg.Opmethod == "transferFrom" {
						fmt.Println("检测到普通代币转账")
						handle.HandleERCContract(&msg)
					}
				} else {
					// 其他非标准 ERC 合约交互，可选择报警/添加到监控池(重要)
					handle.HandleNormalContract(tx, from)
				}
			}
			//----------------普通地址/CEX热钱包交互----------------------
			//是热钱包和普通地址
			if ok, cex := parser.IsCEX(to); ok {
				//触发dex套利信号
				handle.HandleCex(tx, from, cex)
			} else if parser.IsNewAddress(tx) { //新的与巨鲸交互的节点插入图数据库中
				handle.HandleNewAddr(tx, from)
			}
		}
		//------------强跑机器人判断---------------------------
		if tx.GasPrice().Cmp(big.NewInt(200_000_000_000)) > 0 { // >200 Gwei
			log.Println("可能是 Bot 抢跑交易")
		}
		//---------------发现新的巨鲸地址-----------------------
		if parser.IsNewWhale(tx) {
			handle.HandleNewWhale(tx, from)
		}

	}

}
