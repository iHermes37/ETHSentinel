//package main
//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//)
//
//// 定义RPC请求结构体
//type RPCRequest struct {
//	Jsonrpc string        `json:"jsonrpc"`
//	ID      int           `json:"id"`
//	Method  string        `json:"method"`
//	Params  []interface{} `json:"params"`
//}
//
//// 定义RPC响应结构体
//type RPCResponse struct {
//	Jsonrpc string `json:"jsonrpc"`
//	ID      int    `json:"id"`
//	Result  int64  `json:"result"`
//	Error   *struct {
//		Code    int    `json:"code"`
//		Message string `json:"message"`
//	} `json:"error,omitempty"`
//}
//
//func main() {
//	// Helius RPC 地址，注意替换成你的api-key
//	url := "https://mainnet.helius-rpc.com/?api-key=41714ee1-e75b-45be-b8c3-7ffe8ae02f73"
//
//	// 构造请求体
//	reqBody := RPCRequest{
//		Jsonrpc: "2.0",
//		ID:      1,
//		Method:  "getBlockHeight",
//		Params: []interface{}{
//			map[string]interface{}{
//				"commitment":     "finalized",
//				"minContextSlot": 1,
//			},
//		},
//	}
//
//	// 编码为JSON
//	jsonData, err := json.Marshal(reqBody)
//	if err != nil {
//		panic(err)
//	}
//
//	// 发送HTTP POST请求
//	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Body.Close()
//
//	// 读取响应体
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		panic(err)
//	}
//
//	// 解析响应JSON
//	var rpcResp RPCResponse
//	err = json.Unmarshal(body, &rpcResp)
//	if err != nil {
//		panic(err)
//	}
//
//	// 判断是否有错误
//	if rpcResp.Error != nil {
//		fmt.Printf("RPC Error: Code=%d, Message=%s\n", rpcResp.Error.Code, rpcResp.Error.Message)
//		return
//	}
//
//	// 输出区块高度
//	fmt.Printf("当前区块高度: %d\n", rpcResp.Result)
//}
