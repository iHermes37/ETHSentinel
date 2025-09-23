use serde::{Serialize, Deserialize};
use chrono::{DateTime, Utc};

// 交易方向
#[derive(Debug, Serialize, Deserialize)]
pub enum TransactionDirection {
    Incoming,   // 收入
    Outgoing,   // 支出
}

// Solana mempool 交易信息
#[derive(Debug, Serialize, Deserialize)]
pub struct MempoolTx {
    pub signature: String,               // 交易签名（唯一）
    pub from: String,                    // 发送方地址
    pub to: Option<String>,              // 接收方地址，可为空
    pub lamports: u64,                   // 交易金额（单位 lamport）
    pub token: Option<String>,           // Token mint 地址，可选
    pub block_slot: Option<u64>,         // 所属区块 Slot，可选
    pub direction: TransactionDirection, // 交易方向
    pub timestamp: DateTime<Utc>,        // 交易被广播时间
    pub instructions: Vec<String>,       // 交易指令摘要
}
