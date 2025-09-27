package service

import (
	"github.com/Crypto-ChainSentinel/server/internal/model"
	"github.com/Crypto-ChainSentinel/server/internal/repository"
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
