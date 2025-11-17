package whaler

import (
	"github.com/Crypto-ChainSentinel/modules/Analysis/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type SonarSettings struct {
	DetectByHoldings           bool
	DetectByETH                bool
	DetectByDefi               bool
	DetectByToken              bool
	DetectByTransactionPattern bool
}

type Detect interface {
	DetectByHoldings()
	DetectByETH()
	DetectByDefi()
	DetectByToekn()
	DetectByTransactionPattern()
}

type Sonar struct {
	common.Analyst
	Address string
}

func (W *Sonar) DetectByHoldings() {

}

func (W *Sonar) DetectByETH(tx *types.Transaction) {

}

func (W *Sonar) DetectByTransactionPattern() {

}
