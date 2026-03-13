package scanner

import (
	Conn "github.com/Crypto-ChainSentinel/internal"
	"github.com/Crypto-ChainSentinel/internal/parser"
	"github.com/Crypto-ChainSentinel/internal/parser/comm"
	ParserEngineCommon "github.com/Crypto-ChainSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ScanEngine interface {
	Init()
	ParseTranByLog(
		tranreceipt *types.Receipt,
		selectedProtocols map[comm.ProtocolTypeName][]comm.ProtocolImplName,

	) []comm.UnifiedEvent

	ScanBlock(
		block *types.Block,

		EthCh chan<- *types.Transaction,
		TokenCh chan<- *types.Receipt,
		DefiCh chan<- *types.Receipt,
		NewContractCh chan<- *types.Receipt,

	) [][]ParserEngineCommon.UnifiedEvent

	ScanBlocks(
		cfg Interval,
	) chan [][]ParserEngineCommon.UnifiedEventData
}

type Scanner struct {
	Parser parser.Parser
	Client *ethclient.Client
}

func NewScanner() *Scanner {
	connMgr := Conn.ConnMgr{}
	Infura := connMgr.SelectNode("Infura")
	client := Infura.Connect(Conn.WS)
	return &Scanner{
		Client: client,
	}
}

func (s *Scanner) Init(scfg *ScannCfg) {
	s.Parser.RegisterAllParser()
	s.Parser.SetParser(&scfg.Selected)
}
