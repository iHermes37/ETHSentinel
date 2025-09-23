use serde::{Serialize, Deserialize};
use chrono::{DateTime, Utc};

// Solana 巨鲸账户信息
#[derive(Debug, Serialize, Deserialize)]
pub struct Whale {
    pub id: i64,                       // 数据库主键，自增
    pub address: String,                // 巨鲸地址（Solana Pubkey）
    pub first_seen: DateTime<Utc>,      // 第一次发现时间
    pub note: Option<String>,           // 可选备注，如发现来源或策略标签
    pub created_at: DateTime<Utc>,      // 记录创建时间
    pub updated_at: DateTime<Utc>,      // 记录更新时间
}
