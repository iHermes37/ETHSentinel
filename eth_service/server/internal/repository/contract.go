package repository

import (
	"github.com/Crypto-ChainSentinel/initialize"
	"github.com/Crypto-ChainSentinel/models"
	"github.com/Crypto-ChainSentinel/server/internal/model"
	"gorm.io/gorm"
)

type ContractRepository struct {
	db *gorm.DB
}

func NewContractRepository() *ContractRepository {
	return &ContractRepository{
		db: initialize.InitMysql(),
	}
}

// 按字段查询，可支持 address、contractType、txHash
func (r *ContractRepository) Query(cq model.ContractQueryParams, page int, pageSize int) ([]model.ContractQueryResponse, error) {
	var contracts []model.ContractQueryResponse
	query := r.db.Model(&models.ConstractInfo{})

	if cq.Address != nil {
		query = query.Where("address = ?", cq.Address)
	}
	if cq.ContractType != nil {
		query = query.Where("contract_type = ?", cq.ContractType)
	}
	if cq.TxHash != nil {
		query = query.Where("tx_hash = ?", cq.TxHash)
	}
	if cq.ContractAge != nil {
		query = query.Where("contract_age = ?", cq.TxHash)
	}
	if cq.DeployTime != nil {
		query = query.Where("deploy_time = ?", cq.TxHash)
	}

	// 分页处理
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&contracts).Error
	return contracts, err
}
