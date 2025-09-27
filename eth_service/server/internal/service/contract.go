package service

import (
	"github.com/Crypto-ChainSentinel/server/internal/model"
	"github.com/Crypto-ChainSentinel/server/internal/repository"
)

type ContractService struct {
	repo *repository.ContractRepository
}

func NewContractService() *ContractService {
	return &ContractService{
		repo: repository.NewContractRepository(),
	}
}

func (s *ContractService) QueryContracts(cq model.ContractQueryParams, pageStr int, pageSizeStr int) ([]model.ContractQueryResponse, error) {
	return s.repo.Query(cq, pageStr, pageSizeStr)
}
