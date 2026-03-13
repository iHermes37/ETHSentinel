package whaler

import (
	"github.com/Crypto-ChainSentinel/internal/parser/comm"
	"github.com/Crypto-ChainSentinel/internal/scanner"
	"github.com/Crypto-ChainSentinel/server/schemas"
	"github.com/Crypto-ChainSentinel/test/db"
	"github.com/Crypto-ChainSentinel/utils"
	"github.com/Crypto-ChainSentinel/utils/address_mapper"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type SonarSettings struct {
	DetectByHoldings             bool
	DetectByETH                  bool //普通ETH转账检测
	DetectByDefi                 bool //Defi交互检测
	DetectByToken                bool //代币转账检测
	DetectByTransactionFrequency bool
	ParseEngineSettings          *map[comm.ProtocolType][]comm.ProtocolImpl
}

type Detect interface {
	DetectByHoldings()
	DetectByETH()
	DetectByDefi()
	DetectByToken()
	DetectByTransactionFrequency()
}

type Sonar struct {
	*SonarSettings
	ThresholdWei       *big.Int
	ThresholdFrequency *big.Int
	RedisMgr           *db.RedisMgr
}

func NewSonar(thresholdETH *big.Int) *Sonar {
	redisMgr := db.NewRedisMgr() // 提前初始化

	return &Sonar{
		ThresholdWei: utils.ETHToWei(thresholdETH),
		RedisMgr:     redisMgr, // 注入依赖
	}
}

func (s *Sonar) DetectByHoldings() {

}

func (s *Sonar) DetectByETH(tx *types.Transaction) {
	eth_value := tx.Value()
	if eth_value.Cmp(s.ThresholdWei) >= 0 {
		// 普通ETH转账,进行拦截
		from_whale := utils.GetFromAddr(tx)
		s.RedisMgr.StoreWhaleToMonitorPool(from_whale)
	}
}

func (s *Sonar) DetectByToken(recipt *types.Receipt) {
	tran_evlist := scanner.ParseTranByLog(recipt, *s.SonarSettings.ParseEngineSettings)
	for _, tran_ev := range tran_evlist {
		base := tran_ev.CoreEvent()
		ref_tokens := base.RefTokens
		from_whale := base.From
		for _, ref_token := range ref_tokens {
			eth_value := utils.TokenToWei(&ref_token.Name, &ref_token.Amount)
			if eth_value.Cmp(s.ThresholdWei) >= 0 {
				// 代币转账,进行拦截
				s.RedisMgr.StoreWhaleToMonitorPool(&from_whale)
			}
		}

	}
}

// https://etherscan.io/tx/0x2e6532d635a880886a6082f358b93a5cb515d39525e2309d654aea4fbd72f375#eventlog
func (s *Sonar) DetectByDefi(recipt *types.Receipt) {
	tran_evlist := scanner.ParseTranByLog(recipt, *s.SonarSettings.ParseEngineSettings)
	for _, tran_ev := range tran_evlist {
		base := tran_ev.CoreEvent()
		ref_tokens := base.RefTokens
		from_whale := base.From

		if address_mapper.IsDefiProtocol(&from_whale) != nil {
			return
		}

		for _, ref_token := range ref_tokens {
			eth_value := utils.TokenToWei(&ref_token.Name, &ref_token.Amount)
			if eth_value.Cmp(s.ThresholdWei) >= 0 {
				// defi交互,进行拦截
				s.RedisMgr.StoreWhaleToMonitorPool(&from_whale)
				break
			}
		}
	}
}

func (s *Sonar) DetectByTransactionFrequency(recipt *types.Receipt) {
	//判断监控池中的巨鲸的频率是否超过阈值
}

// ================================================================
// 统计阈值
func (s *Sonar) CollectTranFrequency(receipt *types.Receipt) {
	tran_evlist := scanner.ParseTranByLog(receipt, *s.SonarSettings.ParseEngineSettings)
	for _, tran_ev := range tran_evlist {
		base := tran_ev.CoreEvent()
		from_whale := base.From
		var defi_interact bool
		defi_interact = false

		if address_mapper.IsDefiProtocol(&from_whale) != nil {
			defi_interact = true
			return
		}

		if s.RedisMgr.IsWhaleInMonitorPool(&from_whale) {
			if defi_interact {
				s.RedisMgr.UpdateWhaleDeFiTransactionCount(&from_whale)
			}
			s.RedisMgr.UpdateWhaleTransactionCount(&from_whale)
		} else {
			s.RedisMgr.StoreWhaleToMonitorPool(&from_whale)
		}
	}
}

// ==============================================================================
func (s *Sonar) DetectByAddr(addr string) ([]*schemas.WhaleAssetsResponse, error) {
	return nil, nil
}
