package models

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type Whale struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"` // 数据库主键
	Address   common.Address `json:"address" gorm:"type:varchar(42);uniqueIndex"`
	FirstSeen time.Time      `json:"firstSeen" gorm:"type:datetime;index"`    // 第一次发现时间
	Note      string         `json:"note,omitempty" gorm:"type:varchar(255)"` // 可选：备注，记录发现来源
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
}
