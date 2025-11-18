package Filter

import (
	ParserEngineCommon "github.com/Crypto-ChainSentinel/modules/parse_engine/comm"
	"github.com/ethereum/go-ethereum/common"
)

func JudgeIsWhale(to common.Address, from common.Address, cfg TrackWhaleConfig) bool {
	return true
}

func JudgeIsCex(to common.Address) bool {
	return true
}

func JudgeIsTargetProtocol(to common.Address, selected map[ParserEngineCommon.ProtocolType][]ParserEngineCommon.ProtocolImpl) bool {
	return true
}
