package service

import (
	"github.com/CryptoQuantX/chain_monitor/server/internal/model"
	"github.com/CryptoQuantX/chain_monitor/server/internal/repository"
)

type ArbitrageService struct {
	repo *repository.ArbitrageRepository
}

func NewArbitrageService() *ArbitrageService {
	return &ArbitrageService{
		repo: repository.NewArbitrageRepository(),
	}
}

func (s *ArbitrageService) QueryOpportunity(ao model.ArbitrageOpportunityParams, pageStr int, pageSizeStr int) ([]model.ArbitrageOpportunityResponse, error) {
	return s.repo.QueryOpportunity(ao, pageStr, pageSizeStr)
}
