use tokio_tungstenite::{tungstenite::protocol::Message};
mod subscription;
mod utils;
mod models;
mod initialize;
mod modules;
use modules::analysis::monitWhale::monit_whale;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    monit_whale().await?;
}