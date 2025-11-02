use serde::Serialize;

#[derive(Default, Serialize)]
pub  struct whale_trans{
    pub block_hash : String,
    pub block_time:Option<i64>,
    pub block_height:Option<u64>

}