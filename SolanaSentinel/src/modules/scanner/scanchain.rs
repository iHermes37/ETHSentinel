use solana_client::rpc_client::RpcClient;
use solana_client::rpc_config::RpcBlockConfig;

use crate::modules::scanner::parser::parse_transaction;

pub async fn scanmain(){
    let client=RpcClient::new(
        "https://mainnet.helius-rpc.com/?api-key=41714ee1-e75b-45be-b8c3-7ffe8ae02f73".to_string(),
    );
    let slot=client.get_slot().expect("failed to get slot");
    let block_config = RpcBlockConfig {
        max_supported_transaction_version: Some(0), // 限制客户端版本为 0
        ..RpcBlockConfig::default()
    };
    // let myblock=client.get_block(slot).expect("获取失败");
    // 获取区块
    let block=client.get_block_with_config(slot,block_config).expect("获取失败");


    parse_transaction::parse_trans(&block);

}