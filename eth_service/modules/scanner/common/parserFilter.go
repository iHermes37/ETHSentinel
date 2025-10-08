package common

import (
	"github.com/Crypto-ChainSentinel/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// 过滤设置
type FilterConfig struct {
	//单笔交易阈值
	Threshold int
	// 需要跟踪的地址
	TrackAddr common.Address

	// -------------------------------
	findwhale       bool
	trackwhale      bool
	NewContract     bool // 是否需要获取新部署的合约
	FindArbitargBot bool // 套利机器人发现
	MonitorMev      bool //交易池监控
}

func KnownAddress() utils.Set[common.Address] {

	s := utils.NewSet("apple")
	s.Add("cherry")

	return s
}

func ParserFilter(trans *types.Transaction, filtersetting FilterConfig) bool {

	if trans.To() != nil {
		//合约创建---------------------
		//新增代币合约
		//新增潜在defi项目/套利/抢跑机器人发现

		return false
	}

	if trans.To() != nil && len(trans.Data()) == 0 {
		// 普通ETH转账

		return false
	}

	if trans.To() != nil && len(trans.Data()) > 0 {
		// fmt.Println("合约调用交易")

		// 已知合约项目调用
		s := KnownAddress()
		if s.Contains(*trans.To()) {
			return true
		} else {
			// 未知合约项目调用

			return false
		}

	}

	return false

}
