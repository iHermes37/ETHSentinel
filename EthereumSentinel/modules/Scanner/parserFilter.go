package Scanner

import (
	"fmt"
	"github.com/Crypto-ChainSentinel/modules/ConnManager"
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
	"math"
	"math/big"

	filter "github.com/Crypto-ChainSentinel/modules/Scanner/Filter"
	"github.com/Crypto-ChainSentinel/utils"
	"github.com/ethereum/go-ethereum/core/types"
)

// 初步筛除交易，筛出 ETH交易/合约部署/非巨鲸相关
func ParserFilter(trans *types.Transaction, cfg filter.FilterConfig, selected map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl) (ParserEngineCommon.ProtocolType, bool) {
	cli := ConnManager.InfuraConn()
	from := utils.Parsefrom(cli, trans)
	to := trans.To()

	switch cfg.Filter {
	case filter.FindWhale:
		if to != nil && len(trans.Data()) == 0 {
			ethValue := new(big.Float).Quo(
				new(big.Float).SetInt(trans.Value()),
				big.NewFloat(math.Pow10(18)), // wei -> ETH
			)
			threshold := big.NewFloat(100.0) // 巨鲸阈值，比如 100 ETH
			if ethValue.Cmp(threshold) >= 0 {
				// 普通ETH转账,进行拦截
				filter.HandleNewWhale(trans, from)
				return _, true
			}
		}
		break
	case filter.TrackWhale:
		// 判断是否巨鲸的ETH相关操作
		if filter.JudgeIsWhale(from, *to, *cfg.TrackCfg) && filter.JudgeIsTargetProtocol(*to, selected) {
			// 巨鲸ETH转账
			if trans.To() != nil && len(trans.Data()) == 0 {
				if filter.JudgeIsCex(*to) {
					//触发dex套利信号
					filter.HandleCex()
					return true
				} else {
					//新的与巨鲸交互的节点,普通ETH转账
					filter.HandleNewAddr()
				}
				return true
			}
		}
		break
	case filter.NewContract:
		if trans.To() == nil {
			//新增潜在defi项目/抢跑机器人发现/新增代币合约
			contractaddr := utils.GetNewContractAddr(from, trans.Nonce())
			filter.HandleNewContract(trans, contractaddr)
			return true
		}
		//寻找潜在潜力合约项目
	case filter.FindArbitargBot:
		fmt.Print("等待完成")
		break
	}
	return false
}
