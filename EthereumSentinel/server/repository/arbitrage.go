package repository

import (
	"github.com/Crypto-ChainSentinel/init"
	"github.com/Crypto-ChainSentinel/models"
	"github.com/Crypto-ChainSentinel/server/internal/model"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type ArbitrageRepository struct {
	db *gorm.DB
}

func NewArbitrageRepository() *ArbitrageRepository {
	return &ArbitrageRepository{
		db: init.InitMysql(),
	}
}

func (a *ArbitrageRepository) QueryOpportunity(cq model.ArbitrageOpportunityParams, page int, pageSize int) ([]model.ArbitrageOpportunityResponse, error) {
	var results []model.ArbitrageOpportunityResponse

	// 起始查询
	query := a.db.Model(&models.CrossPairData{})

	// ===== 动态条件构造 =====
	if cq.Token0 != (common.Address{}) {
		query = query.Where("pair_dex_a->'token0'->>'id' = ? OR pair_dex_b->'token0'->>'id' = ?", cq.Token0.Hex(), cq.Token0.Hex())
	}
	if cq.Token1 != (common.Address{}) {
		query = query.Where("pair_dex_a->'token1'->>'id' = ? OR pair_dex_b->'token1'->>'id' = ?", cq.Token1.Hex(), cq.Token1.Hex())
	}
	if cq.DexA != "" {
		query = query.Where("pair_dex_a->>'exchange' = ?", cq.DexA)
	}
	if cq.DexB != "" {
		query = query.Where("pair_dex_b->>'exchange' = ?", cq.DexB)
	}
	if cq.Direction != "" && cq.Direction != "Any" {
		query = query.Where("direction = ?", cq.Direction)
	}
	if cq.BorrowToken != (common.Address{}) {
		query = query.Where("borrow_token->>'id' = ?", cq.BorrowToken.Hex())
	}
	if cq.MinProfit > 0 {
		query = query.Where("net_profit >= ?", cq.MinProfit)
	}
	if cq.Amount > 0 {
		query = query.Where("x >= ?", cq.Amount)
	}

	// ===== 分页逻辑 =====
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
