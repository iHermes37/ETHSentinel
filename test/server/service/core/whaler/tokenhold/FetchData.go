package tokenhold

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 配置常量
const (
	DUNE_API = "https://api.dune.com/api/v1"
	// QUERY_ID = "3237025" // 在 Dune 上保存的查询ID
)

// 获取查询状态
func GetQueryStatus(apiKey, queryID string) (string, error) {
	url := fmt.Sprintf("%s/query/%s/status", DUNE_API, queryID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-dune-api-key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var statusResp struct {
		State string `json:"state"`
	}
	if err := json.Unmarshal(body, &statusResp); err != nil {
		return "", err
	}

	return statusResp.State, nil
}

// 获取查询结果
func GetQueryResult(apiKey, queryID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/query/%s/results", DUNE_API, queryID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-dune-api-key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var resultResp struct {
		Result struct {
			Rows []map[string]interface{} `json:"rows"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &resultResp); err != nil {
		return nil, err
	}

	return resultResp.Result.Rows, nil
}

// 获取数据结构
func FetchData(apiKey, QUERYID string) ([]map[string]interface{}, error) {
	// 轮询查询状态
	for {
		status, err := GetQueryStatus(apiKey, QUERYID)
		if err != nil {
			fmt.Println("获取查询状态失败:", err)
			return nil, err
		}
		if status == "QUERY_STATE_COMPLETED" {
			break
		}
		fmt.Println("查询中，等待中...")
		time.Sleep(3 * time.Second)
	}

	// 获取查询结果
	rows, err := GetQueryResult(apiKey, QUERYID)
	if err != nil {
		fmt.Println("获取查询结果失败:", err)
		return nil, err
	}

	// 打印结果
	fmt.Println("USDT 前100持仓地址：")
	for i, row := range rows {
		fmt.Printf("%d. 地址: %s, 余额: %v\n", i+1, row["holder"], row["balance"])
	}

	return rows, nil
}

// func ComputeHoldings(rows []map[string]interface{}) (map[string]float64, error) {

// 	var results map[string]float64
// 	for i, row := range rows {
// 		fmt.Printf("%d. 地址: %s, 余额: %v\n", i+1, row["holder"], row["balance"])
// 		holder := row["holder"]
// 		results[string(holder)] = row["balance"]
// 	}
// }
