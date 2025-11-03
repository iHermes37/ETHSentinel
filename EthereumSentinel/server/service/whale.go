package service

import (
	"fmt"

	"github.com/Crypto-ChainSentinel/server/model"
	"github.com/Crypto-ChainSentinel/server/repository"
)

type WhaleService struct {
	repo *repository.WhaleRepository
}

func NewWhaleService() *WhaleService {
	return &WhaleService{
		repo: repository.NewWhaleRepository(),
	}
}

func (w *WhaleService) QueryWhaleTrade(wt model.WhaleTradeParams, page int, pageSize int) ([]model.WhaleTradeResponse, error) {
	return w.repo.QueryWhaleTrade(wt, page, pageSize)
}

func (w *WhaleService) QueryCapturedWhales(wp model.WhaleQueryParams, page int, pageSize int) ([]model.WhaleResponse, error) {
	return w.repo.QueryCapturedWhales(wp, page, pageSize)
}

func (h *WhaleService) CapturedWhales(task model.CapturedWhale) error {
	switch task.Method {
	case model.HoldingsAnalysis:
		break
	case model.ChainScan:
		break
	case model.TransactionPattern:
		break
	default:
		return fmt.Errorf("不支持的检测方法: %s", task.Method)
	}
	return nil
}

func (h *WhaleService) TrackWhales(task model.TrackWhale) error {

}
