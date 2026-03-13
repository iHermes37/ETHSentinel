package service

import (
	"github.com/Crypto-ChainSentinel/server/db"
	"github.com/Crypto-ChainSentinel/server/schemas"
)

type ContractService struct {
	repo *db.ContractRepository
}

func NewContractService() *ContractService {
	return &ContractService{
		repo: db.NewContractRepository(),
	}
}

func (s *ContractService) QueryContracts(cq model.ContractQueryParams, pageStr int, pageSizeStr int) ([]model.ContractQueryResponse, error) {
	return s.repo.Query(cq, pageStr, pageSizeStr)
}
