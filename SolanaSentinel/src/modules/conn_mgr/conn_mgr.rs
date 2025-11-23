use tokio_tungstenite::{connect_async, tungstenite::client::IntoClientRequest, WebSocketStream};
use tokio::net::TcpStream;
use tokio_tungstenite::MaybeTlsStream;
use std::error::Error;

pub async fn ws_connect()-> Result<WebSocketStream<MaybeTlsStream<TcpStream>>, Box<dyn Error>>{
    let request="wss://mainnet.helius-rpc.com/?api-key=41714ee1-e75b-45be-b8c3-7ffe8ae02f73".into_client_request()?;

    let (ws_stream,_)=connect_async(request).await.map_err(|e| {
            eprintln!("WebSocket 连接失败: {:?}", e);
            e
        })?;

    return Ok(ws_stream);
}


use solana_client::rpc_client::RpcClient;
use solana_client::rpc_config::RpcBlockConfig;
use serde_json::to_string_pretty;
use std::fs::File;
use serde_json::to_writer_pretty;

pub async fn scan_main(){
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

    // 使用 serde_json 输出成 JSON
    let json = to_string_pretty(&block).expect("序列化失败");
    println!("{}", json);

    // 保存到文件
    let file = File::create("block.json").expect("无法创建文件");
    to_writer_pretty(file, &block).expect("写入文件失败");

    println!("区块已保存到 block.json");

}


// ===========================================================
use std::fmt;



// 连接类型枚举
#[derive(Debug, Clone, PartialEq)]
pub enum Conn {
    WebSocket(String),  // URL
    Grpc(String),       // 端点地址
    Https(String),      // URL
}

// 自定义错误类型
#[derive(Debug)]
pub enum ConnError {
    InvalidUrl(String),
    ConnectionFailed(String),
    UnsupportedProtocol,
    Timeout,
}
impl std::error::Error for ConnError {}
impl fmt::Display for ConnError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            ConnError::InvalidUrl(url) => write!(f, "无效的URL: {}", url),
            ConnError::ConnectionFailed(msg) => write!(f, "连接失败: {}", msg),
            ConnError::UnsupportedProtocol => write!(f, "不支持的协议"),
            ConnError::Timeout => write!(f, "连接超时"),
        }
    }
}


// 连接管理器
pub struct ConnMgr {
    conn_type: Option<Conn>,
    is_initialized: bool,
}


impl ConnMgr{

    /// 创建新的连接管理器
    pub fn new() -> Self {}

    /// 设置连接类型
    pub fn set_conn(&mut self, conn_type: Conn) -> Result<(), ConnError> {}

    /// 初始化连接
    pub fn init(&mut self) -> Result<(), ConnError> {}

    /// 发送数据
    pub fn send_data(&self, data: &[u8]) -> Result<(), ConnError> {}

    /// 获取连接状态
    pub fn is_connected(&self) -> bool {}

    /// 获取连接
    pub fn get_conn(&self) -> Option<&Conn> {}

    /// 私有方法
    fn validate_conn(conn_type: &Conn) -> Result<(), ConnError> {}

}

