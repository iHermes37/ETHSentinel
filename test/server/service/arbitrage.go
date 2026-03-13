package service

import (
	"github.com/Crypto-ChainSentinel/server/db"
	"github.com/Crypto-ChainSentinel/server/schemas"
)

type ArbitrageService struct {
	repo *db.ArbitrageRepository
}

func NewArbitrageService() *ArbitrageService {
	return &ArbitrageService{
		repo: db.NewArbitrageRepository(),
	}
}

func (s *ArbitrageService) QueryOpportunity(ao model.ArbitrageOpportunityParams, pageStr int, pageSizeStr int) ([]model.ArbitrageOpportunityResponse, error) {
	return s.repo.QueryOpportunity(ao, pageStr, pageSizeStr)
}
