package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

// ProtocolInfo 定义协议信息
type ProtocolInfo struct {
	Protocol string `json:"protocol"`
	Type     string `json:"type"`
}

// ProtocolMap 多链协议映射表
type ProtocolMap map[string]map[string]ProtocolInfo

// LoadProtocolMap 从 JSON 文件加载映射表
func LoadProtocolMap(filePath string) (ProtocolMap, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var pm ProtocolMap
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&pm); err != nil {
		return nil, err
	}

	// 将所有合约地址转为小写，便于查询
	for _, addrMap := range pm {
		for addr, info := range addrMap {
			lower := strings.ToLower(addr)
			if lower != addr {
				addrMap[lower] = info
				delete(addrMap, addr)
			}
		}
	}
	return pm, nil
}

// GetProtocol 根据 chainID 和 to 查询协议信息
func GetProtocol(pm ProtocolMap, chainID string, to string) (ProtocolInfo, bool) {
	to = strings.ToLower(to)
	if chainProtocols, ok := pm[chainID]; ok {
		if info, exists := chainProtocols[to]; exists {
			return info, true
		}
	}
	return ProtocolInfo{}, false
}
