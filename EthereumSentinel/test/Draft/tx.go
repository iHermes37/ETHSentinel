package scanner

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"github.com/ethereum/go-ethereum/erccommon"
//	"github.com/ethereum/go-ethereum/core/types"
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/ethereum/go-ethereum/ethclient"
//	"github.com/ethereum/go-ethereum/rlp"
//	"log"
//	"math/big"
//	"os"
//	"strings"
//)
//
//func IsWhale(addr erccommon.Address) bool {
//	_,data:=os.ReadFile("")
//
//	var raw map[srting]interface{}
//	json.Unmarshal(data,&raw)
//
//	_,whalelist:=raw["whale_addresses"]
//
//	for _,whale := range whalelist{
//			if whale["address"]==addr{
//				return true
//			}
//	}
//}
//
//func IsCEX(to *erccommon.Address) (bool, string) {
//
//	_, data := os.ReadFile("")
//
//	var raw map[string]interface{}
//
//	json.Unmarshal(data, &raw)
//
//	cexsRaw, ok := raw["CEX"].(map[string]interface{})
//
//	if !ok {
//		return false, ""
//	}
//
//	for cexName, addresses := range cexsRaw {
//
//		addrSlice, ok := addresses.([]interface{})
//		if !ok {
//			continue
//		}
//
//		for _, addr := range addrSlice {
//			if addrStr, ok := addr.(string); ok {
//				if strings.ToLower(addrStr) == strings.ToLower(to) {
//					return true, cexName, nil
//				}
//			}
//		}
//
//		return (false,nil)
//
//	}
//
//}
//// 判断是否 DEX Router
//func IsNewAddress(tx *types.Transaction) bool {
//	return false
//}
//
//func IsNewWhale(tx *types.Transaction) bool {
//	return false
//}
//
//func AnalyzeWhale(tx *types.Transaction) {
//	to := tx.To()
//	data := tx.Data()
//	value := tx.Value()
//
//	// 1. 判断接收方
//	if IsCEX(to) {
//		log.Println("鲸鱼入金/出金 -> CEX", to.Hex(), "金额", value)
//	} else if IsDEXRouter(to) {
//		log.Println("鲸鱼操作 DEX ->", to.Hex())
//	} else if IsNewAddress(to) {
//		log.Println("鲸鱼转账到新地址:", to.Hex())
//	} else {
//		log.Println("链上普通交易:", to.Hex(), "金额", value)
//	}
//
//	// 2. 判断行为
//	if len(data) == 0 {
//		log.Println("普通转账 ETH")
//	} else if bytes.HasPrefix(data, []byte{0xa9, 0x05, 0x9c, 0xbb}) {
//		log.Println("ERC20 转账")
//	} else if bytes.HasPrefix(data, []byte{0x7f, 0xf3, 0x6a, 0xb5}) {
//		log.Println("Uniswap SwapExactETHForTokens")
//	} else {
//		log.Println("其他合约调用")
//	}
//
//	// 3. Gas 分析
//	if tx.GasPrice().Cmp(big.NewInt(200_000_000_000)) > 0 { // >200 Gwei
//		log.Println("可能是 Bot 抢跑交易")
//	}
//}
//
//func Parsefrom(cli *ethclient.Client, tx *types.Transaction) erccommon.Address {
//	chainID, _ := cli.NetworkID(context.Background())
//	signer := types.LatestSignerForChainID(chainID)
//	from, err := types.Sender(signer, tx)
//
//	if err != nil {
//		fmt.Println("Error decoding sender:", err)
//		continue
//	}
//
//	return from
//}
//
//func GetNewContractAddr(sender erccommon.Address, nonce uint64) erccommon.Address {
//
//	// RLP 编码 [sender, nonce]
//	rlpStream, err := rlp.EncodeToBytes([]interface{}{sender, nonce})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// keccak256 哈希
//	hash := crypto.Keccak256(rlpStream)
//
//	// 取最后 20 字节
//	return erccommon.BytesToAddress(hash[12:])
//
//}
