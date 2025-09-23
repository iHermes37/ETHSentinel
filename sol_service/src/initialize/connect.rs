use tokio_tungstenite::{connect_async, tungstenite::client::IntoClientRequest, WebSocketStream};
use tokio::net::TcpStream;
use tokio_tungstenite::MaybeTlsStream;
use std::error::Error;

pub async fn ws_connect()-> Result<WebSocketStream<MaybeTlsStream<TcpStream>>, Box<dyn Error>>{
    let request="wss://mainnet.helius-rpc.com/?api-key=41714ee1-e75b-45be-b8c3-7ffe8ae02f73".into_client_request()?;

    let (ws_stream,response)=connect_async(request).await.map_err(|e| {
            eprintln!("WebSocket 连接失败: {:?}", e);
            e
        })?;

    println!("WebSocket 已连接");

    return Ok(ws_stream);
}