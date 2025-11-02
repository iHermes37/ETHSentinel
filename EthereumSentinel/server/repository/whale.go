package repository

import (
	"github.com/Crypto-ChainSentinel/initialize"
	"github.com/Crypto-ChainSentinel/models"
	"github.com/Crypto-ChainSentinel/server/internal/model"
	"gorm.io/gorm"
)

type WhaleRepository struct {
	db *gorm.DB
}

func NewWhaleRepository() *WhaleRepository {
	return &WhaleRepository{
		db: initialize.InitMysql(),
	}
}

func (r *WhaleRepository) QueryCapturedWhales(wq model.WhaleQueryParams, page int, pageSize int) ([]model.WhaleResponse, error) {
	var whales []model.WhaleResponse
	query := r.db.Model(&models.Whale{})

	if wq.Address != nil {
		query = query.Where("address=?", wq.Address)
	}

	if wq.FirstSeen != nil {
		query = query.Where("first_seen=?", wq.Address)
	}

	// 分页处理
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&whales).Error

	return whales, err
}

func (r *WhaleRepository) QueryWhaleTrade(wq model.WhaleTradeParams, page int, pageSize int) ([]model.WhaleTradeResponse, error) {
	db := r.db.Model(&models.WhaleTransaction{})

	// ===== BaseWhaleTransaction =====
	if wq.Address != nil {
		db = db.Where("address = ?", *wq.Address)
	}
	if wq.Type != nil {
		db = db.Where("type = ?", *wq.Type)
	}
	if wq.TxHash != nil {
		db = db.Where("tx_hash = ?", *wq.TxHash)
	}
	if wq.To != nil {
		db = db.Where("to = ?", *wq.To)
	}
	if wq.StartTime != nil {
		db = db.Where("time >= ?", *wq.StartTime)
	}
	if wq.EndTime != nil {
		db = db.Where("time <= ?", *wq.EndTime)
	}

	// ===== DeFiTxDetail =====
	if wq.Exchange != nil {
		db = db.Where("defi_exchange = ?", *wq.Exchange)
	}
	if wq.Direction != nil {
		db = db.Where("defi_direction = ?", *wq.Direction)
	}
	if wq.TokenIn != nil {
		db = db.Where("defi_tokens_in_symbol = ?", *wq.TokenIn)
	}
	if wq.TokenOut != nil {
		db = db.Where("defi_tokens_out_symbol = ?", *wq.TokenOut)
	}
	if wq.AmountInMin != nil {
		db = db.Where("defi_amounts_in >= ?", *wq.AmountInMin)
	}
	if wq.AmountInMax != nil {
		db = db.Where("defi_amounts_in <= ?", *wq.AmountInMax)
	}
	if wq.AmountOutMin != nil {
		db = db.Where("defi_amounts_out >= ?", *wq.AmountOutMin)
	}
	if wq.AmountOutMax != nil {
		db = db.Where("defi_amounts_out <= ?", *wq.AmountOutMax)
	}

	// ===== ERC20TxDetail =====
	if wq.ERC20Token != nil {
		db = db.Where("erc20_token = ?", *wq.ERC20Token)
	}
	if wq.ERC20AmountMin != nil {
		db = db.Where("erc20_amount >= ?", *wq.ERC20AmountMin)
	}
	if wq.ERC20AmountMax != nil {
		db = db.Where("erc20_amount <= ?", *wq.ERC20AmountMax)
	}

	// ===== UserTxDetail =====
	if wq.UserAsset != nil {
		db = db.Where("user_asset = ?", *wq.UserAsset)
	}
	if wq.UserAmountMin != nil {
		db = db.Where("user_amount >= ?", *wq.UserAmountMin)
	}
	if wq.UserAmountMax != nil {
		db = db.Where("user_amount <= ?", *wq.UserAmountMax)
	}

	offset := (page - 1) * pageSize
	var trades []model.WhaleTradeResponse
	if err := db.Offset(offset).Limit(pageSize).Find(&trades).Error; err != nil {
		return nil, err
	}

	return trades, nil
}
