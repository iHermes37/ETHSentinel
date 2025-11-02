use solana_client::rpc_client::{RpcClient, GetConfirmedSignaturesForAddress2Config};
use solana_client::rpc_config::RpcTransactionConfig;
use solana_sdk::{pubkey::Pubkey, signature::Signature, commitment_config::CommitmentConfig};
use solana_transaction_status::UiTransactionEncoding;

pub fn raydium_test() -> Result<(), Box<dyn std::error::Error>> {
    let rpc_url = "https://api.mainnet-beta.solana.com";
    let client = RpcClient::new(rpc_url.to_string());

    let pool_address = "5quBtoiQqxF9Jv6KYKctB59NT3gtJD2Y65kdnB1Uev3h".parse::<Pubkey>()?;

    let signatures = client.get_signatures_for_address_with_config(
        &pool_address,
        GetConfirmedSignaturesForAddress2Config {
            before: None,
            until: None,
            limit: Some(10),
            commitment: Some(CommitmentConfig::finalized()),
        },
    )?;

    for sig_info in signatures.iter() {
        let sig = sig_info.signature.parse::<Signature>()?;
        let tx = client.get_transaction_with_config(
            &sig,
            RpcTransactionConfig {
                encoding: Some(UiTransactionEncoding::Json),
                commitment: Some(CommitmentConfig::finalized()),
                max_supported_transaction_version: Some(0), // v0
            },
        )?;

        println!("block_time is:"+tx.block_time)
   
    }

    Ok(())
}
