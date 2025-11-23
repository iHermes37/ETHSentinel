mod utils;
pub mod modules;

// use modules::analysis::monitWhale::monit_whale;
// use modules::scanner::scanchain::scanmain;

use solana_client::pubsub_client::PubsubClient;
use solana_client::rpc_config::RpcProgramAccountsConfig;
use solana_sdk::pubkey::Pubkey;
use solana_account_decoder::UiAccountEncoding;

// use solana_streamer_sdk::dex_parser::core::event_parser::EventParser;
// use solana_streamer_sdk::dex_parser::Protocol;
// use solana_streamer_sdk::dex_parser::UnifiedEvent;
use std::str::FromStr;
use std::sync::Arc;

fn main() {
    let rpc_url = "wss://mainnet.helius-rpc.com/?api-key=41714ee1-e75b-45be-b8c3-7ffe8ae02f73";
    use solana_sdk::pubkey::Pubkey;
    let amm_program_id: Pubkey = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8".parse().unwrap();
    // 订阅 AMM 程序账户变化
    let (mut client, receiver) = PubsubClient::program_subscribe(
        rpc_url,
        &amm_program_id,
        Some(RpcProgramAccountsConfig {
            account_config: solana_client::rpc_config::RpcAccountInfoConfig {
                encoding: Some(UiAccountEncoding::Base64),
                ..Default::default()
            },
            ..Default::default()
        })
    ).unwrap();
    for msg in receiver {
        let account_info = msg.value.account;
         let pubkey = &msg.value.pubkey;   // 这里才有 pubkey
        // 在这里解析 account_info.data，判断是不是新池子
        println!("Detected AMM account change: {:?}", pubkey);
        // println!("Account data length: {}", account_info.data.len());
    }
}



