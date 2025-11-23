// https://solana.com/zh/docs/rpc/websocket/accountsubscribe


// {
//   "jsonrpc": "2.0",
//   "id": 1,
//   "method": "accountSubscribe",
//   "params": [
//     "CM78CPUeXjn8o3yroDHxUtKsZZgoy4GPkPPXfouKNH12",
//     {
//       "encoding": "jsonParsed",
//       "commitment": "finalized"
//     }
//   ]
// }



// {
//   "jsonrpc": "2.0",
//   "method": "accountNotification",
//   "params": {
//     "result": {
//       "context": {
//         "slot": 5199307
//       },
//       "value": {
//         "data": {
//           "program": "nonce",
//           "parsed": {
//             "type": "initialized",
//             "info": {
//               "authority": "Bbqg1M4YVVfbhEzwA9SpC9FhsaG83YMTYoR4a8oTDLX",
//               "blockhash": "LUaQTmM7WbMRiATdMMHaRGakPtCkc2GHtH57STKXs6k",
//               "feeCalculator": {
//                 "lamportsPerSignature": 5000
//               }
//             }
//           }
//         },
//         "executable": false,
//         "lamports": 33594,
//         "owner": "11111111111111111111111111111111",
//         "rentEpoch": 635,
//         "space": 80
//       }
//     },
//     "subscription": 23784
//   }
// }







#[derive(Serialize, Debug, Clone)]
pub struct AccountSubscribeParams {
    pub pubkey: String,
    pub config: AccountConfig,
}

#[derive(Serialize, Debug, Clone)]
pub struct AccountConfig {
    pub encoding: String,
    pub commitment: String,
}





use serde::{Deserialize, Serialize};

// https://solana.com/zh/docs/rpc/websocket/logssubscribe
//--------------------请求-------------------------------
// {
//   "jsonrpc": "2.0",
//   "id": 1,
//   "method": "logsSubscribe",
//   "params": [
//     {
//       "mentions": ["11111111111111111111111111111111"]
//     },
//     {
//       "commitment": "finalized"
//     }
//   ]
// }


#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct LogsSubscribeRequest {
    pub jsonrpc: String,
    pub id: u64,
    pub method: String,
    pub params: Vec<serde_json::Value>,   // 注意：这里是 Vec<Value>，因为 params 本身是数组
}

//----------------------------响应----------------------------
// {
//   "jsonrpc": "2.0",
//   "method": "logsNotification",
//   "params": {
//     "result": {
//       "context": {
//         "slot": 5208469
//       },
//       "value": {
//         "signature": "5h6xBEauJ3PK6SWCZ1PGjBvj8vDdWG3KpwATGy1ARAXFSDwt8GFXM7W5Ncn16wmqokgpiKRLuS83KUxyZyv2sUYv",
//         "err": null,
//         "logs": [
//           "SBF program 83astBRguLMdt2h5U1Tpdq5tjFoJ6noeGwaY3mDLVcri success"
//         ]
//       }
//     },
//     "subscription": 24040
//   }
// }


#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct LogsNotificationContext{
    pub slot: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct LogsNotificationValue{
    signature: String,
    err: Option<serde_json::Value>,
    logs: Vec<String>,
    
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct LogsNotificationResult {
    pub context: LogsNotificationContext,
    pub value: LogsNotificationValue,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct LogsNotificationParams {
    pub result: LogsNotificationResult,
    pub subscription: u64,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct LogsNotificationResponse {
    pub jsonrpc: String,
    pub method: String,
    pub params: LogsNotificationParams,
}