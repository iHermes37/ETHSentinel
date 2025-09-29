package ERC

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var erc721ABIJSon = `[
	{
	"anonymous": false,
	"inputs": [
	{
	"indexed": true,
	"internalType": "address",
	"name": "owner",
	"type": "address"
	},
	{
	"indexed": true,
	"internalType": "address",
	"name": "approved",
	"type": "address"
	},
	{
	"indexed": true,
	"internalType": "uint256",
	"name": "tokenId",
	"type": "uint256"
	}
	],
	"name": "Approval",
	"type": "event"
	},
	{
	"anonymous": false,
	"inputs": [
	{
	"indexed": true,
	"internalType": "address",
	"name": "owner",
	"type": "address"
	},
	{
	"indexed": true,
	"internalType": "address",
	"name": "operator",
	"type": "address"
	},
	{
	"indexed": false,
	"internalType": "bool",
	"name": "approved",
	"type": "bool"
	}
	],
	"name": "ApprovalForAll",
	"type": "event"
	},
	{
	"anonymous": false,
	"inputs": [
	{
	"indexed": true,
	"internalType": "address",
	"name": "from",
	"type": "address"
	},
	{
	"indexed": true,
	"internalType": "address",
	"name": "to",
	"type": "address"
	},
	{
	"indexed": true,
	"internalType": "uint256",
	"name": "tokenId",
	"type": "uint256"
	}
	],
	"name": "Transfer",
	"type": "event"
	},
	{
	"inputs": [
	{
	"internalType": "address",
	"name": "to",
	"type": "address"
	},
	{
	"internalType": "uint256",
	"name": "tokenId",
	"type": "uint256"
	}
	],
	"name": "approve",
	"outputs": [],
	"stateMutability": "nonpayable",
	"type": "function"
	},
	{
	"constant": true,
	"inputs": [],
	"name": "totalSupply",
	"outputs": [
	{
	"name": "",
	"type": "uint256"
	}
	],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
	},
	{
	"inputs": [
	{
	"internalType": "address",
	"name": "owner",
	"type": "address"
	}
	],
	"name": "balanceOf",
	"outputs": [
	{
	"internalType": "uint256",
	"name": "balance",
	"type": "uint256"
	}
	],
	"stateMutability": "view",
	"type": "function"
	},
	{
	"inputs": [
	{
	"internalType": "uint256",
	"name": "tokenId",
	"type": "uint256"
	}
	],
	"name": "getApproved",
	"outputs": [
	{
	"internalType": "address",
	"name": "operator",
	"type": "address"
	}
	],
	"stateMutability": "view",
	"type": "function"
	},
	{
	"inputs": [
	{
	"internalType": "address",
	"name": "owner",
	"type": "address"
	},
	{
	"internalType": "address",
	"name": "operator",
	"type": "address"
	}
	],
	"name": "isApprovedForAll",
	"outputs": [
	{
	"internalType": "bool",
	"name": "",
	"type": "bool"
	}
	],
	"stateMutability": "view",
	"type": "function"
	},
	{
	"inputs": [],
	"name": "name",
	"outputs": [
	{
	"internalType": "string",
	"name": "",
	"type": "string"
	}
	],
	"stateMutability": "view",
	"type": "function"
	},
	{
	"inputs": [
	{
	"internalType": "uint256",
	"name": "tokenId",
	"type": "uint256"
	}
	],
	"name": "ownerOf",
	"outputs": [
	{
	"internalType": "address",
	"name": "owner",
	"type": "address"
	}
	],
	"stateMutability": "view",
	"type": "function"
	},
	{
	"inputs": [
	{
	"internalType": "address",
	"name": "from",
	"type": "address"
	},
	{
	"internalType": "address",
	"name": "to",
	"type": "address"
	},
	{
	"internalType": "uint256",
	"name": "tokenId",
	"type": "uint256"
	}
	],
	"name": "safeTransferFrom",
	"outputs": [],
	"stateMutability": "nonpayable",
	"type": "function"
	},
	{
	"inputs": [
	{
	"internalType": "address",
	"name": "from",
	"type": "address"
	},
	{
	"internalType": "address",
	"name": "to",
	"type": "address"
	},
	{
	"internalType": "uint256",
	"name": "tokenId",
	"type": "uint256"
	},
	{
	"internalType": "bytes",
	"name": "data",
	"type": "bytes"
	}
	],
	"name": "safeTransferFrom",
	"outputs": [],
	"stateMutability": "nonpayable",
	"type": "function"
	},
	{
	"inputs": [
	{
	"internalType": "address",
	"name": "operator",
	"type": "address"
	},
	{
	"internalType": "bool",
	"name": "_approved",
	"type": "bool"
	}
	],
	"name": "setApprovalForAll",
	"outputs": [],
	"stateMutability": "nonpayable",
	"type": "function"
	},
	{
	"inputs": [
	{
	"internalType": "bytes4",
	"name": "interfaceId",
	"type": "bytes4"
	}
	],
	"name": "supportsInterface",
	"outputs": [
	{
	"internalType": "bool",
	"name": "",
	"type": "bool"
	}
	],
	"stateMutability": "view",
	"type": "function"
	},
	{
	"inputs": [],
	"name": "symbol",
	"outputs": [
	{
	"internalType": "string",
	"name": "",
	"type": "string"
	}
	],
	"stateMutability": "view",
	"type": "function"
	},
	{
	"inputs": [
	{
	"internalType": "uint256",
	"name": "tokenId",
	"type": "uint256"
	}
	],
	"name": "tokenURI",
	"outputs": [
	{
	"internalType": "string",
	"name": "",
	"type": "string"
	}
	],
	"stateMutability": "view",
	"type": "function"
	},
	{
	"inputs": [
	{
	"internalType": "address",
	"name": "from",
	"type": "address"
	},
	{
	"internalType": "address",
	"name": "to",
	"type": "address"
	},
	{
	"internalType": "uint256",
	"name": "tokenId",
	"type": "uint256"
	}
	],
	"name": "transferFrom",
	"outputs": [],
	"stateMutability": "nonpayable",
	"type": "function"
	}
]`

func ParseERC721Tx(data []byte, ercABI abi.ABI) (string, map[string]interface{}) {
	// 输入长度小于 4，无法解析函数选择器
	if len(data) < 4 {
		fmt.Println("非合约调用或普通 ETH 转账，tx.Data() 太短")
		return "", make(map[string]interface{})
	}

	method, err := ercABI.MethodById(data[:4])
	if err != nil {
		fmt.Println("未知操作或非 ERC721 调用:", err)
		return "", make(map[string]interface{})
	}

	// 解析参数
	params := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(params, data[4:]); err != nil {
		fmt.Println("参数解析失败:", err)
		return method.Name, make(map[string]interface{})
	}

	// 根据方法名打印关键参数（可选）
	switch method.Name {
	case "approve":
		fmt.Println("授权目标地址 (to):", params["to"])
		fmt.Println("授权 tokenId:", params["tokenId"])
	case "setApprovalForAll":
		fmt.Println("操作员地址 (operator):", params["operator"])
		fmt.Println("是否授权 (approved):", params["_approved"])
	case "transferFrom":
		fmt.Println("转移源地址 (from):", params["from"])
		fmt.Println("转移目标地址 (to):", params["to"])
		fmt.Println("转移 tokenId:", params["tokenId"])
	case "safeTransferFrom":
		fmt.Println("安全转移源地址 (from):", params["from"])
		fmt.Println("安全转移目标地址 (to):", params["to"])
		fmt.Println("安全转移 tokenId:", params["tokenId"])
		if _, ok := params["data"]; ok {
			fmt.Println("附加数据 (data):", params["data"])
		}
	default:
		fmt.Println("未处理的方法:", method.Name)
	}

	return method.Name, params
}
