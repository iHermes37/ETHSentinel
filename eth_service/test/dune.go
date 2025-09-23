//package main
//
//import (
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"log"
//	"net/http"
//)
//
//func main() {
//	apiKey := "c1QaFjpS6LaXXkDob7vG3dDdtbkdaL9F"
//	url := "https://api.dune.com/api/v1/query/5270617/results?limit=1000"
//
//	// 创建请求
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 添加 API Key Header
//	req.Header.Add("X-Dune-API-Key", apiKey)
//
//	// 发起请求
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer resp.Body.Close()
//
//	// 读取响应
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 定义结构体解析 JSON
//	type Result struct {
//		Rows []map[string]interface{} `json:"rows"`
//	}
//
//	type DuneResponse struct {
//		Result Result `json:"result"`
//	}
//
//	var duneResp DuneResponse
//	if err := json.Unmarshal(body, &duneResp); err != nil {
//		log.Fatal(err)
//	}
//
//	// 打印前 100 条
//	for i, row := range duneResp.Result.Rows {
//		if i >= 100 {
//			break
//		}
//		fmt.Printf("%d: %v\n", i+1, row)
//	}
//}
