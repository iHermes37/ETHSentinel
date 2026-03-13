package uniswapv2

import (
	"github.com/Crypto-ChainSentinel/internal/parser/comm"
)

func Register() comm.ProtocolImplParser {
	event_parse_invoker := comm.NewEventParseInvoker()
	event_parse_invoker.Register(UniswapV2EventsConfig)

	var ProtocolImplParser comm.ProtocolImplParser
	ProtocolImplParser = event_parse_invoker

	return ProtocolImplParser
}
