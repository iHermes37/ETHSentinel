package ParserEngine

import "github.com/Crypto-ChainSentinel/modules/ParserEngine/dex_parser"

type GenParser string

const (
	ERCParser GenParser = "ERCParser"
	DEXParser GenParser = "DEXParser"
)

type chooseparams struct {
	parser *GenParser
	token  *string
	dex    *string
	event  *string
}

type ParserStrategy struct {
	DexParser *dex_parser.MyEventParser
}

type ParserFactory struct{}

func CreateParser(params chooseparams) ParserStrategy {
	parser := *params.parser
	Parser := ParserStrategy{}
	switch parser {
	case DEXParser:
		dexparser := dex_parser.MyEventParser{}
		dexparser.NewParser(dex, event)
		Parser.DexParser = &dexparser
		break
	default:
		return ParserStrategy{}
	}

	return ParserStrategy{}
}
