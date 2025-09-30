package parser

import (
	"context"
	"fmt"
	"github.com/Crypto-ChainSentinel/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"strings"
)

func IsContract(client *ethclient.Client, addr common.Address) bool {
	bytecode, err := client.CodeAt(context.Background(), addr, nil)
	if err != nil {
		return false
	}
	return len(bytecode) > 0
}

func DetectERCStandard(tx *types.Transaction) (models.ERCStandard, error) {
	data := tx.Data()

	// ERC20
	erc20ABI, _ := abi.JSON(strings.NewReader(Erc.erc20ABIJSon))
	if _, err := erc20ABI.MethodById(data[:4]); err == nil {
		erc := models.ERCStandard{}
		erc.Ercname = models.ContractType("ERC20")
		erc.Opmethod, erc.Params = Erc.ParseERC20Tx(data, erc20ABI)
		return erc, nil
	}

	// ERC721
	erc721ABI, _ := abi.JSON(strings.NewReader(Erc.erc721ABIJSon))
	if _, err := erc721ABI.MethodById(data[:4]); err == nil {
		erc := models.ERCStandard{}
		erc.Ercname = models.ContractType("ERC721")                   // ✅ 修改
		erc.Opmethod, erc.Params = Erc.ParseERC721Tx(data, erc721ABI) // ✅ 修改 ABI
		return erc, nil
	}

	return models.ERCStandard{}, fmt.Errorf("未识别 ERC 标准")
}
