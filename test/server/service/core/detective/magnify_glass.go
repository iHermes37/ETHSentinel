package detective

import (
	"github.com/Crypto-ChainSentinel/internal/parser/comm"
	"github.com/Crypto-ChainSentinel/internal/scanner"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Search interface {
	SearchByDeFiTranFrequency()
	SearchByFlashLoan()
	SearchByGasPrice()
}

type MagnifyGlass struct {
	ThresholdFrequency  *big.Int
	ThresholdGasPrice   *big.Int
	ParseEngineSettings *map[comm.ProtocolType][]comm.ProtocolImpl
}

func (mg *MagnifyGlass) SearchByDeFiTranFrequency() {
	//	redis_curd
}

func (mg *MagnifyGlass) SearchByGasPrice(tx *types.Transaction) {
	gasprice := tx.GasPrice()
	if gasprice.Cmp(mg.ThresholdGasPrice) > 0 {
		//加入可疑池中
	}

}

func (mg *MagnifyGlass) SearchByFlashLoan(receipt *types.Receipt) {
	tranEvlist := scanner.ParseTranByLog(receipt, *mg.ParseEngineSettings)
	for _, tranEv := range tranEvlist {
		protocolType := tranEv.ProtocolType()
		eventType := tranEv.EventType()
		if protocolType == comm.Lending && eventType == comm.FlashLoan {
			//加入可疑池中
		}
	}

}
