package tokenHold

import "fmt"

type HoldingAnalyzer interface {
	FetchData() ([]map[string]interface{}, error)
	// ParseData() error
	// ComputeHoldings(rows []map[string]interface{}) (map[string]float64, error)
}

func GetAnalyzer(token string) (HoldingAnalyzer, error) {
	switch token {
	case "USDC":
		return &USDCAnalyzer{}, nil
	case "USDT":
		return &USDTAnalyzer{}, nil
	default:
		return nil, fmt.Errorf("unsupported token: %s", token)
	}
}
