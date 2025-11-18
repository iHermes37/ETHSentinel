package Filter

import (
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/parse_engine/comm"
	"github.com/ethereum/go-ethereum/common"
)

type FilterSetting bool

// 过滤设置
var (
	FindWhale       FilterSetting = true //发现巨鲸
	TrackWhale      FilterSetting = true //跟踪巨鲸
	NewContract     FilterSetting = true // 是否需要获取新部署的合约
	FindArbitargBot FilterSetting = true // 套利机器人发现
)

type TrackWhaleConfig struct {
	TargetAddress common.Address
	IsAllWhale    bool
}

type FindWhaleConfig struct {
}

type FilterConfig struct {
	Filter    FilterSetting
	TrackCfg  *TrackWhaleConfig
	FindWhale *FindWhaleConfig
}

type ScanTransConfig struct {
	BeforFilter       FilterConfig
	SelectedProtocols *map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl
}
