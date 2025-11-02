use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct TokenBalance{
    pub mint: String,      // Token Mint 地址
    pub amount: u64,       // Token 数量（原始 lamports）
    pub decimals: u8,      // Token 小数位
}

/// 巨鲸账户信息
#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct WhaleAccount {
    pub pubkey: String,                       // 巨鲸公钥
    pub lamports: u64,                        // SOL 余额
    pub last_slot: u64,                        // 最新更新 slot
    pub tokens: HashMap<String, TokenBalance>, // key: mint 地址, value: TokenBalance
}