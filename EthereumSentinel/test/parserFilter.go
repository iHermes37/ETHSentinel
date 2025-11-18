package Analysis

import (
	"errors"
	"fmt"
	"github.com/Crypto-ChainSentinel/utils"
	"math"
	"math/big"

	"github.com/Crypto-ChainSentinel/modules/ConnManager"
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/parse_engine/comm"

	filter "github.com/Crypto-ChainSentinel/modules/scanner/Filter"
	"github.com/ethereum/go-ethereum/core/types"
)

// 策略接口
type FilterStrategy interface {
	Exec(block *types.Block) []ParserEngineCommon.UnifiedEvent
	Name() string
}

// 策略工厂
type FilterFactory struct {
	strategies map[string]FilterStrategy
}

func NewFilterFactory() *FilterFactory {
	factory := &FilterFactory{
		strategies: make(map[string]FilterStrategy),
	}

	// 注册所有可用策略
	factory.strategies["find_whale"] = &FindWhaleByChainScan{}
	factory.strategies["track_whale"] = &TrackWhaleByChainScan{}
	factory.strategies["monitor_contract"] = &FindMintoredContractByChainScan{}

	return factory
}

func (ff *FilterFactory) Create(filterType string) FilterStrategy {
	return ff.strategies[filterType]
}

func (ff *FilterFactory) GetAvailableStrategies() []string {
	var names []string
	for name := range ff.strategies {
		names = append(names, name)
	}
	return names
}

// 策略管理器
type FilterManager struct {
	factory         *FilterFactory
	currentStrategy FilterStrategy
}

func NewFilterManager() *FilterManager {
	return &FilterManager{
		factory: NewFilterFactory(),
	}
}

func (fm *FilterManager) SetCurrentFilter(filterType string) error {
	strategy := fm.factory.Create(filterType)
	if strategy == nil {
		return fmt.Errorf("未知的策略类型: %s", filterType)
	}
	fm.currentStrategy = strategy
	return nil
}

func (fm *FilterManager) Exec(block *types.Block) ([]ParserEngineCommon.UnifiedEvent, error) {
	if fm.currentStrategy == nil {
		return nil, errors.New("未设置当前策略")
	}
	return fm.currentStrategy.Exec(block), nil
}

func (fm *FilterManager) ExecuteMultiple(block *types.Block, filterTypes []string) (map[string][]ParserEngineCommon.UnifiedEvent, error) {
	results := make(map[string][]ParserEngineCommon.UnifiedEvent)

	for _, filterType := range filterTypes {
		strategy := fm.factory.Create(filterType)
		if strategy == nil {
			return nil, fmt.Errorf("未知策略: %s", filterType)
		}

		events := strategy.Exec(block)
		results[filterType] = events
	}

	return results, nil
}

// ======================================================
type FindWhaleByChainScan struct {
}

func (fw *FindWhaleByChainScan) Exec() {

}

type TrackWhaleByChainScan struct {
}

func (tw *TrackWhaleByChainScan) Exec() {

}

// =========================================================
type FindMintoredContractByChainScan struct {
}

func (fmc *FindMintoredContractByChainScan) Exec() {

}

// =========================================================

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
	}
	return false
}
