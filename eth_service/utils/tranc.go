package utils

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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
