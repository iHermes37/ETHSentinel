package utils

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"io/ioutil"
	"os"
)

// 定义泛型函数
func ReadJSONFile[T any](filePath string) (*T, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var obj T
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

func ReadABIFile(filePath string) (*abi.ABI, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	parsedAbi, err := abi.JSON(f)
	if err != nil {
		return nil, err
	}

	return &parsedAbi, nil
}
