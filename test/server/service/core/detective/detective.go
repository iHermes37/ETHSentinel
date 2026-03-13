package detective

import (
	"github.com/Crypto-ChainSentinel/internal/parser/comm"
	"github.com/Crypto-ChainSentinel/internal/scanner"
	"github.com/Crypto-ChainSentinel/test/db"
	"github.com/Crypto-ChainSentinel/utils/address_mapper"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Interval struct {
	StartBlock *big.Int
	EndBlock   *big.Int
}

type Detective struct {
	Interval
	RedisMgr            *db.RedisMgr
	ParseEngineSettings *map[comm.ProtocolType][]comm.ProtocolImpl
}

func NewDetective(ThresholdFrequency *big.Int, ThresholdGasPrice *big.Int) *Detective {
	redisMgr := db.NewRedisMgr() // 提前初始化
	return &Detective{
		RedisMgr: redisMgr,
	}
}

func (d *Detective) CollectNewContract(receipt *types.Receipt) {
	newContractAddr := receipt.ContractAddress
	d.RedisMgr.StoreContractToMonitorPool(&newContractAddr)
}

func (d *Detective) CollectTranFrequency(receipt *types.Receipt) {
	tranEvlist := scanner.ParseTranByLog(receipt, *d.ParseEngineSettings)
	for _, tranEv := range tranEvlist {
		base := tranEv.CoreEvent()
		fromWhale := base.From
		var defiInteract bool
		defiInteract = false

		if address_mapper.IsDefiProtocol(&fromWhale) != nil {
			defiInteract = true
			return
		}

		if d.RedisMgr.IsContractInMonitorPool(&fromWhale) {
			if defiInteract {
				d.RedisMgr.UpdateContractDeFiTranCount(&fromWhale)
			}
			d.RedisMgr.UpdateContractTransactionCount(&fromWhale)
		}
	}
}
