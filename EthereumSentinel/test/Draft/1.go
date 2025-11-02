//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"log"
//	"net/http
//	"net/url"
//	"time"
//)
//
//// RPC 请求结构体
//type RPCRequest struct {
//	Jsonrpc string        `json:"jsonrpc"`
//	ID      int           `json:"id"`
//	Method  string        `json:"method"`
//	Params  []interface{} `json:"params"`
//}
//
//// RPC 响应结构体
//type RPCResponse struct {
//	Jsonrpc string `json:"jsonrpc"`
//	ID      int    `json:"id"`
//	Result  uint64 `json:"result"`
//	Error   *struct {
//		Code    int    `json:"code"`
//		Message string `json:"message"`
//	} `json:"error,omitempty"`
//}
//
//func main() {
//	// 设置代理地址（HTTP代理或HTTPS代理都可）
//	proxyStr := "http://127.0.0.1:7890" // 这里改成你的代理地址和端口
//
//	proxyURL, err := url.Parse(proxyStr)
//	if err != nil {
//		log.Fatalf("代理地址解析失败: %v", err)
//	}
//
//	// 创建带代理的HTTP客户端
//	client := &http.Client{
//		Transport: &http.Transport{
//			Proxy: http.ProxyURL(proxyURL),
//		},
//		Timeout: 10 * time.Second,
//	}
//
//	// 构造RPC请求体
//	rpcReq := RPCRequest{
//		Jsonrpc: "2.0",
//		ID:      1,
//		Method:  "getBlockHeight",
//		Params:  []interface{}{map[string]interface{}{"commitment": "finalized"}},
//	}
//
//	reqBody, err := json.Marshal(rpcReq)
//	if err != nil {
//		log.Fatalf("请求序列化失败: %v", err)
//	}
//
//	// 目标RPC地址（Solana官方主网RPC）
//	rpcURL := "https://api.mainnet-beta.solana.com"
//
//	// 发起POST请求
//	req, err := http.NewRequestWithContext(context.Background(), "POST", rpcURL, bytes.NewBuffer(reqBody))
//	if err != nil {
//		log.Fatalf("创建请求失败: %v", err)
//	}
//	req.Header.Set("Content-Type", "application/json")
//
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Fatalf("请求失败: %v", err)
//	}
//	defer resp.Body.Close()
//
//	// 读取响应
//	respBody, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatalf("读取响应失败: %v", err)
//	}
//
//	// 解析JSON响应
//	var rpcResp RPCResponse
//	err = json.Unmarshal(respBody, &rpcResp)
//	if err != nil {
//		log.Fatalf("响应解析失败: %v", err)
//	}
//
//	// 检查RPC错误
//	if rpcResp.Error != nil {
//		log.Fatalf("RPC错误: %d - %s", rpcResp.Error.Code, rpcResp.Error.Message)
//	}
//
//	fmt.Printf("当前区块高度: %d\n", rpcResp.Result)
//}
