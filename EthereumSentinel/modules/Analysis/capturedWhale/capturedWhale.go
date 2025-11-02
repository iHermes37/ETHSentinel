package capturedWhale

import "github.com/Crypto-ChainSentinel/modules/Analysis/common"

type CapturedWhale interface {
	DetectByHoldings()
	DetectByChainScan()
	DetectByTransactionPattern()
}

type Whale struct {
	common.Analyst
	Address string
}

func (W *Whale) DetectByHoldings() {

}

func (W *Whale) DetectByChainScan() {

}

func (W *Whale) DetectByTransactionPattern() {

}
