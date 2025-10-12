package common

import (
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
	"github.com/ethereum/go-ethereum/common"
)

type TrackWhaleConfig struct {
	TargetAddress common.Address
	IsAllWhale    bool
}

type FilterConfig struct {
	Filter   FilterSetting
	TrackCfg *TrackWhaleConfig
}

type ScanTransConfig struct {
	BeforFilter       FilterConfig
	SelectedProtocols *map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl
}
