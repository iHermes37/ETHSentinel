package service

import (
	"fmt"
	"github.com/Crypto-ChainSentinel/server/db"
	"github.com/Crypto-ChainSentinel/server/schemas"
	"github.com/Crypto-ChainSentinel/server/service/core/whaler"
)

type WhaleService struct {
	repo   *db.WhaleRepository
	Whaler *whaler.Whaler
}

func NewWhaleService() *WhaleService {
	return &WhaleService{
		repo: db.NewWhaleRepository(),
	}
}

// ==========================================================
func (w *WhaleService) QueryWhaleTrade(wt schemas.WhaleTradeParams, page int, pageSize int) ([]schemas.WhaleTradeResponse, error) {
	return w.repo.QueryWhaleTrade(wt, page, pageSize)
}

func (w *WhaleService) QueryCapturedWhales(wp schemas.WhaleQueryParams, page int, pageSize int) ([]schemas.WhaleResponse, error) {
	return w.repo.QueryCapturedWhales(wp, page, pageSize)
}

//===============================================

func (h *WhaleService) CapturedWhales(task schemas.CapturedWhale) error {
	switch task.Method {
	case schemas.HoldingsAnalysis:
		break
	case schemas.ChainScan:
		break
	case schemas.TransactionPattern:
		break
	default:
		return fmt.Errorf("不支持的检测方法: %s", task.Method)
	}
	return nil
}

func (h *WhaleService) TrackWhales(task schemas.TrackWhale) error {

}

// ===============================================

func (s *WhaleService) QueryWhaleAssets(addr string) ([]*schemas.WhaleAssetsResponse, error) {
	response, err := s.repo.FindAssetsByAddress(addr)
	if err == nil {
		return response, nil
	}
	response, err = s.Whaler.Sonar.DetectByAddr(addr)
	if err != nil {
		return nil, err
	}
	return response, nil
}
