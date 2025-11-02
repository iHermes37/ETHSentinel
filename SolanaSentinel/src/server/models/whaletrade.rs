use serde::{Serialize, Deserialize};
use chrono::{DateTime, Utc};

#[derive(Debug, Serialize, Deserialize)]
pub enum TransactionDirection {
    Incoming,
    Outgoing,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct WhaleTransaction {
    pub id: i64,                           // 数据库主键，自增
    pub time: DateTime<Utc>,               // 交易时间
    pub address: String,                   // 巨鲸地址（Pubkey）
    pub tx_type: String,                    // 行为类型，例如 swap、transfer
    pub token_in: String,                  // 输入代币 Mint 地址
    pub token_out: Option<String>,         // 输出代币 Mint，可空
    pub amount_in: u64,                    // 输入代币数量（Lamports 或 Token 最小单位）
    pub amount_out: Option<u64>,           // 输出代币数量
    pub direction: TransactionDirection,   // 交易方向 Incoming / Outgoing
    pub signature: String,                 // Solana 交易签名（唯一）
    pub slot: Option<u64>,                 // 交易所在 Slot，可选
    pub exchange: Option<String>,          // 所属交易所，可空
    pub to: Option<String>,                // 接收方地址，可空
    pub created_at: DateTime<Utc>,         // 记录创建时间
    pub updated_at: DateTime<Utc>,         // 更新时间
}


pub struct WhaleTradeList{
    pub Trades:Vec![WhaleTransaction]
}
