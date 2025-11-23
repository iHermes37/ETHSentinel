// 账户订阅模块


// 日志订阅请求


// 合约订阅模块


use serde::{Deserialize, Serialize};

use crate::models::notification::chainlog::LogsNotificationResponse;

#[derive(Serialize, Deserialize, Debug, Clone)]
pub enum SolEventTypes {
    LogNotification(LogsNotificationResponse),
    // AccountNotification(AccountNotificationResponse),
    // Program(ProgramNotificationResponse)
}


// ================================================
#[derive(Debug, Clone)]
pub enum WsRequestMethod {
    AccountSubscribe(AccountSubscribeParams),
    ProgramSubscribe(ProgramSubscribeParams),
    LogsSubscribe(LogsSubscribeParams),
    SlotSubscribe(SlotSubscribeParams),
    SignatureSubscribe(SignatureSubscribeParams),
}


impl WsRequestMethod {
    pub fn method_name(&self) -> &'static str {
        match self {
            WsRequestMethod::AccountSubscribe(_) => "accountSubscribe",
            WsRequestMethod::ProgramSubscribe(_) => "programSubscribe",
            WsRequestMethod::LogsSubscribe(_) => "logsSubscribe",
            WsRequestMethod::SlotSubscribe(_) => "slotSubscribe",
            WsRequestMethod::SignatureSubscribe(_) => "signatureSubscribe",
        }
    }
}


#[derive(Serialize)]
pub struct JsonRpcRequest<T> {
    pub jsonrpc: &'static str,
    pub id: u64,
    pub method: &'static str,
    pub params: T,
}

impl WsRequestMethod {
    pub fn into_json_request(self, id: u64) -> JsonRpcRequest<serde_json::Value> {
        let method = self.method_name();
        let params = match self {
            WsRequestMethod::AccountSubscribe(p) => serde_json::json!([p.pubkey, p.config]),
            WsRequestMethod::ProgramSubscribe(p) => serde_json::to_value(p).unwrap(),
            WsRequestMethod::LogsSubscribe(p) => serde_json::to_value(p).unwrap(),
            WsRequestMethod::SlotSubscribe(p) => serde_json::to_value(p).unwrap(),
            WsRequestMethod::SignatureSubscribe(p) => serde_json::to_value(p).unwrap(),
        };

        JsonRpcRequest {
            jsonrpc: "2.0",
            id,
            method,
            params,
        }
    }
}
