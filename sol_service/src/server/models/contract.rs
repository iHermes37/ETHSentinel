use serde::{Serialize, Deserialize};
use chrono::{DateTime, Utc};
use sqlx::types::BigDecimal;

// 合约类型枚举
#[derive(Debug, Serialize, Deserialize, sqlx::Type)]
#[sqlx(type_name = "varchar")] // SQL 数据库类型
pub enum ContractType {
    Token,
    DeFi,
    NFT,
    Other,
}

// Solana Program 部署信息
#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct ContractInfo {
    pub id: i64,                                      // 数据库主键，自增
    pub address: String,                              // 合约地址，Solana Pubkey
    pub contract_type: ContractType,                  // 合约类型
    pub contract_age: Option<i64>,                    // 合约年龄（秒，可选）
    pub deploy_tx_signature: String,                  // 部署交易签名
    pub slot: Option<u64>,                            // 部署所在区块 Slot
    pub deploy_time: Option<DateTime<Utc>>,           // 部署时间
    pub created_at: DateTime<Utc>,                    // 记录创建时间
    pub updated_at: DateTime<Utc>,                    // 记录更新时间
}