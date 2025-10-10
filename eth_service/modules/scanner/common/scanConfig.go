package common

import (
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/ParserEngine/common"
	"github.com/ethereum/go-ethereum/common"
)

type ScanTransConfig struct {
	WhaleAddr         *common.Address
	SelectedProtocols *map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl
}
