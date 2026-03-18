// Package erc721 实现 ERC721 NFT 标准事件解析
package erc721

import (
	"fmt"

	"github.com/ETHSentinel/internal/parser/comm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Parser ERC721 事件解析器
type Parser struct {
	invoker *comm.EventParseInvoker
}

// NewParser 创建 ERC721 解析器
func NewParser() *Parser {
	p := &Parser{
		invoker: comm.NewEventParseInvoker(comm.ProtocolImplERC721),
	}
	// ERC721 的 Transfer 签名与 ERC20 相同（indexed tokenId 而非 value）
	p.invoker.RegisterOne(comm.SigERC721Transfer, comm.EventMethodTransfer, p.parseTransfer)
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

func (p *Parser) parseTransfer(log types.Log, meta comm.EventMetadata) (comm.UnifiedEvent, error) {
	if len(log.Topics) < 4 {
		return nil, fmt.Errorf("erc721: transfer log: expected 4 topics, got %d", len(log.Topics))
	}

	from := common.BytesToAddress(log.Topics[1].Bytes())
	to := common.BytesToAddress(log.Topics[2].Bytes())
	// tokenId 在 Topics[3]（indexed）
	tokenIdHash := log.Topics[3]

	return &comm.UnifiedEventData{
		Metadata: comm.EventMetadata{
			TxHash:           meta.TxHash,
			ProtocolTypeVal:  comm.ProtocolTypeToken,
			ProtocolImplVal:  comm.ProtocolImplERC721,
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
				{Name: log.Address.Hex(), Amount: tokenIdHash.Big()},
			},
		},
		DetailVal: map[string]any{
			"contract": log.Address.Hex(),
			"from":     from.Hex(),
			"to":       to.Hex(),
			"tokenId":  tokenIdHash.Big().String(),
		},
	}, nil
}
