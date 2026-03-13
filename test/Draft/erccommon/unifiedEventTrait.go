package erccommon

import "math/big"

type TokenEvent struct {
	Protocol     string   // ERC20 or ERC721
	EventType    string   // Transfer / Approval / ApprovalForAll
	TokenAddress string   // 合约地址
	Operator     string   // 操作者地址 (spender/operator)
	From         string   // 转出地址
	To           string   // 转入地址
	TokenId      *big.Int // NFT ID (ERC721 专用，ERC20 为 nil)
	Amount       *big.Int // 转账数量 (ERC20 专用，ERC721 固定为 1)
	Approved     *bool    // 是否授权 (ERC721 的 setApprovalForAll 专用)
	TxHash       string   // 交易哈希
	BlockNumber  uint64   // 区块高度
}

type MethodName string

const (
	Transfer     MethodName = "transfer"
	TransferFrom MethodName = "transferFrom"
	Approve      MethodName = "approve"
)

type Protocol string

const (
	ERC20 Protocol = "ERC20"
)
