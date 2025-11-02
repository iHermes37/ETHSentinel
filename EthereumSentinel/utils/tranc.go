package utils

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func Parsefrom(cli *ethclient.Client, tx *types.Transaction) common.Address {
	chainID, _ := cli.NetworkID(context.Background())
	signer := types.LatestSignerForChainID(chainID)
	from, err := types.Sender(signer, tx)

	if err != nil {
		fmt.Println("Error decoding sender:", err)
	}

	return from
}

func GetNewContractAddr(sender common.Address, nonce uint64) common.Address {

	// RLP 编码 [sender, nonce]
	rlpStream, err := rlp.EncodeToBytes([]interface{}{sender, nonce})
	if err != nil {
		log.Fatal("")
	}

	// keccak256 哈希
	hash := crypto.Keccak256(rlpStream)

	// 取最后 20 字节
	return common.BytesToAddress(hash[12:])

}
