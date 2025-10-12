package common

import (
	connectionManager "github.com/Crypto-ChainSentinel/modules/ConnectionManager"
	handle "github.com/Crypto-ChainSentinel/modules/Scanner/Handle"
	"github.com/Crypto-ChainSentinel/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// //单笔交易阈值
// Threshold int
// // 需要跟踪的地址
// TrackAddr common.Address
type FilterSetting bool
// 过滤设置
var(
	Findwhale       FilterSetting=true //发现巨鲸
	Trackwhale      FilterSetting=true //跟踪巨鲸
	NewContract     FilterSetting=true // 是否需要获取新部署的合约
	FindArbitargBot FilterSetting=true // 套利机器人发现
)


func KnownAddress() utils.Set[common.Address] {

	s := utils.NewSet("apple")
	s.Add("cherry")

	return s
}

func ParserFilter(trans *types.Transaction, cfg FilterConfig,) bool{
	cli := connectionManager.InfuraConn()
	from := utils.Parsefrom(cli, trans)

	switch cfg.Filter{
	case Findwhale:
		if trans.To() != nil && len(trans.Data()) == 0 {
			// 普通ETH转账
			handle.HandleNewWhale(trans, from)
			return  true
		};
		break;
	case Trackwhale:
		if (from!= && to!=){
			return  true
		};
		if trans.To() != nil && len(trans.Data()) == 0 {
			if(IsCEX(to)){
				//触发dex套利信号
//				handle.HandleCex(tx, from, cex)
			}else{
				//新的与巨鲸交互的节点,普通ETH转账
//				handle.HandleNewAddr(tx, from)
			}
			return  true
		};
		break;
	case NewContract:
		if trans.To() == nil {
			//新增潜在defi项目/套利/抢跑机器人发现/新增代币合约
			contractaddr := utils.GetNewContractAddr(from, trans.Nonce())
			handle.HandleNewContract(trans, contractaddr)
			return true
		}
	case FindArbitargBot:break
	}
	return false
}

