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

// 过滤设置
type FilterSetting struct {
	Findwhale       bool //发现巨鲸
	Trackwhale      bool //跟踪巨鲸
	NewContract     bool // 是否需要获取新部署的合约
	FindArbitargBot bool // 套利机器人发现
}

func KnownAddress() utils.Set[common.Address] {

	s := utils.NewSet("apple")
	s.Add("cherry")

	return s
}

func ParserFilter(trans *types.Transaction, filtersetting FilterSetting) bool{

	cli := connectionManager.InfuraConn()
	from := utils.Parsefrom(cli, trans)

	if filtersetting.NewContract {
		if trans.To() == nil {
			//新增潜在defi项目/套利/抢跑机器人发现/新增代币合约
			contractaddr := utils.GetNewContractAddr(from, trans.Nonce())
			handle.HandleNewContract(trans, contractaddr)
			return true
		}
	}

	if filtersetting.Findwhale {
		if trans.To() != nil && len(trans.Data()) == 0 {
			// 普通ETH转账
			handle.HandleNewWhale(trans, from)
			return  true
		}
	}

	if filtersetting.Trackwhale {
		if (from!= && to!=){
			return  true
		}
	}

	return false
}
