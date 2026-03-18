// Package erc20 实现 ERC20 标准事件解析（Transfer / Approval）
package erc20

import (
	"fmt"
	"math/big"

	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// erc20ABI 精简 ABI，只包含 Transfer 事件
const transferABIJSON = `[{
	"anonymous": false,
	"inputs": [
		{"indexed": true,  "name": "from",  "type": "address"},
		{"indexed": true,  "name": "to",    "type": "address"},
		{"indexed": false, "name": "value", "type": "uint256"}
	],
	"name": "Transfer",
	"type": "event"
}]`

// Parser ERC20 事件解析器
type Parser struct {
	invoker    *comm.EventParseInvoker
	transferAB abi.ABI
}

// NewParser 创建 ERC20 解析器
func NewParser() *Parser {
	parsedABI, err := abi.JSON(stringReader(transferABIJSON))
	if err != nil {
		panic(fmt.Sprintf("erc20: parse ABI: %v", err)) // ABI 是常量，panic 合理
	}
	p := &Parser{
		invoker:    comm.NewEventParseInvoker(comm.ProtocolImplERC20),
		transferAB: parsedABI,
	}
	p.invoker.RegisterOne(comm.SigERC20Transfer, comm.EventMethodTransfer, p.parseTransfer)
	return p
}

// ── comm.ProtocolImplParser 接口 ──────────────

func (p *Parser) HandleEvent(sig comm.EventSig, log types.Log, meta comm.EventMetadata) (comm.UnifiedEvent, error) {
	return p.invoker.HandleEvent(sig, log, meta)
}

func (p *Parser) ListEventSigs() []comm.EventSig       { return p.invoker.ListEventSigs() }
func (p *Parser) SetFilter(methods []comm.EventMethod) { p.invoker.SetFilter(methods) }

// ─────────────────────────────────────────────
//  内部解析实现
// ─────────────────────────────────────────────

type transferLog struct {
	Value *big.Int
}

func (p *Parser) parseTransfer(log types.Log, meta comm.EventMetadata) (comm.UnifiedEvent, error) {
	if len(log.Topics) < 3 {
		return nil, fmt.Errorf("erc20: transfer log: expected 3 topics, got %d", len(log.Topics))
	}

	from := common.BytesToAddress(log.Topics[1].Bytes())
	to := common.BytesToAddress(log.Topics[2].Bytes())

	var tl transferLog
	if err := p.transferAB.UnpackIntoInterface(&tl, "Transfer", log.Data); err != nil {
		return nil, fmt.Errorf("erc20: unpack transfer: %w", err)
	}

	return &comm.UnifiedEventData{
		Metadata: comm.EventMetadata{
			TxHash:           meta.TxHash,
			ProtocolTypeVal:  comm.ProtocolTypeToken,
			ProtocolImplVal:  comm.ProtocolImplERC20,
			Age:              meta.Age,
			To:               log.Address,
			BlockNumber:      meta.BlockNumber,
			OuterIndex:       meta.OuterIndex,
			TransactionIndex: meta.TransactionIndex,
		},
		Base: comm.BaseEvent{
			EventType: comm.EventMethodTransfer,
			From:      from,
			RefTokens: []comm.RefToken{
				{Name: log.Address.Hex(), Amount: tl.Value},
			},
		},
		DetailVal: &comm.TransferData{
			Token:  log.Address,
			From:   from,
			To:     to,
			Amount: tl.Value,
		},
	}, nil
}

// stringReader 将字符串包装为 io.Reader（abi.JSON 需要）
type stringReaderType struct {
	s   string
	pos int
}

func stringReader(s string) *stringReaderType { return &stringReaderType{s: s} }
func (r *stringReaderType) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.s) {
		return 0, fmt.Errorf("EOF")
	}
	n = copy(p, r.s[r.pos:])
	r.pos += n
	return n, nil
}
