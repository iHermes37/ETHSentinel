package ERC

import (
	erccommon "github.com/Crypto-ChainSentinel/modules/parserEngine/erc_parser/erccommon"
)

var erc20ABIJSon = `[
						{
							"constant": true,
							"inputs": [],
							"name": "name",
							"outputs": [
								{
									"name": "",
									"type": "string"
								}
							],
							"payable": false,
							"stateMutability": "view",
							"type": "function"
						},
						{
							"constant": false,
							"inputs": [
								{
									"name": "_spender",
									"type": "address"
								},
								{
									"name": "_value",
									"type": "uint256"
								}
							],
							"name": "approve",
							"outputs": [
								{
									"name": "",
									"type": "bool"
								}
							],
							"payable": false,
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
							"constant": false,
							"inputs": [
								{
									"name": "_from",
									"type": "address"
								},
								{
									"name": "_to",
									"type": "address"
								},
								{
									"name": "_value",
									"type": "uint256"
								}
							],
							"name": "transferFrom",
							"outputs": [
								{
									"name": "",
									"type": "bool"
								}
							],
							"payable": false,
							"stateMutability": "nonpayable",
							"type": "function"
						},
						{
							"constant": true,
							"inputs": [],
							"name": "decimals",
							"outputs": [
								{
									"name": "",
									"type": "uint8"
								}
							],
							"payable": false,
							"stateMutability": "view",
							"type": "function"
						},
						{
							"constant": true,
							"inputs": [
								{
									"name": "_owner",
									"type": "address"
								}
							],
							"name": "balanceOf",
							"outputs": [
								{
									"name": "balance",
									"type": "uint256"
								}
							],
							"payable": false,
							"stateMutability": "view",
							"type": "function"
						},
						{
							"constant": true,
							"inputs": [],
							"name": "symbol",
							"outputs": [
								{
									"name": "",
									"type": "string"
								}
							],
							"payable": false,
							"stateMutability": "view",
							"type": "function"
						},
						{
							"constant": false,
							"inputs": [
								{
									"name": "_to",
									"type": "address"
								},
								{
									"name": "_value",
									"type": "uint256"
								}
							],
							"name": "transfer",
							"outputs": [
								{
									"name": "",
									"type": "bool"
								}
							],
							"payable": false,
							"stateMutability": "nonpayable",
							"type": "function"
						},
						{
							"constant": true,
							"inputs": [
								{
									"name": "_owner",
									"type": "address"
								},
								{
									"name": "_spender",
									"type": "address"
								}
							],
							"name": "allowance",
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
							"payable": true,
							"stateMutability": "payable",
							"type": "fallback"
						},
						{
							"anonymous": false,
							"inputs": [
								{
									"indexed": true,
									"name": "owner",
									"type": "address"
								},
								{
									"indexed": true,
									"name": "spender",
									"type": "address"
								},
								{
									"indexed": false,
									"name": "value",
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
									"name": "from",
									"type": "address"
								},
								{
									"indexed": true,
									"name": "to",
									"type": "address"
								},
								{
									"indexed": false,
									"name": "value",
									"type": "uint256"
								}
							],
							"name": "Transfer",
							"type": "event"
						}
					]`

//func getERC20Meta(addr common.Address, client *ethclient.Client) (string, int) {
//	contract, err := abi.JSON(strings.NewReader(erc20ABI))
//	if err != nil {
//		return "", 18
//	}
//	call := bind.NewBoundContract(addr, contract, client, client, client)
//
//	var symbol string
//	var decimals uint8
//	err = call.Call(nil, &symbol, "symbol")
//	if err != nil {
//		symbol = "UNKNOWN"
//	}
//	err = call.Call(nil, &decimals, "decimals")
//	if err != nil {
//		decimals = 18
//	}
//	return symbol, int(decimals)
//}

//func ParseERC20Tx(data []byte, ercABI abi.ABI) erccommon.TokenEvent {
//
//	token := common.TokenEvent{}
//	if len(data) < 4 {
//		fmt.Println("非合约调用或普通 ETH 转账，tx.Data() 太短")
//		return common2.TokenEvent{}
//	}
//
//	method, err := ercABI.MethodById(data[:4])
//	if err != nil {
//		fmt.Println("未知操作或非 ERC20 调用")
//		return common2.TokenEvent{}
//	}
//	fmt.Println("调用函数:", method.Name)
//	// 解析参数
//	params := make(map[string]interface{})
//	err = method.Inputs.UnpackIntoMap(params, data[4:])
//	if err != nil {
//		fmt.Println("参数解析失败:", err)
//		return common2.TokenEvent{}
//	}
//
//	// 打印关键参数
//	if method.Name == "transfer" {
//		token.Protocol = method.Name
//		token.To = common.Address(params["_to"])
//		return method.Name, params
//	} else if method.Name == "approve" {
//		fmt.Println("授权目标地址 (_to):", params["_spender"])
//		fmt.Println("授权数量 (_value):", params["_value"])
//		return method.Name, params
//	} else if method.Name == "transferFrom" {
//		fmt.Println("转移源地址 (_from):", params["_from"])
//		fmt.Println("转移目标地址 (_to):", params["_to"])
//		fmt.Println("转移目标数量 (_value):", params["_value"])
//		return method.Name, params
//	}
//	return "", make(map[string]interface{})
//}

var ERC20EventsConfig = map[erccommon.MethodName]erccommon.EventParserFunc{
	erccommon.Transfer: ParserTransferEvent,
}

var TokenAddr={

}

func ParserTransferEvent(data []byte) (erccommon.TokenEvent, error) {

}
