use crossbeam_channel::{bounded};
use tokio_tungstenite::{
    connect_async, MaybeTlsStream, tungstenite::protocol::Message, WebSocketStream,
};
use futures_util::{SinkExt};

use crate::initialize::connect::ws_connect; 
use crate::models::notification::eventTypes::SolEventTypes;
use crate::utils::WsStreamTrans::streamMapORM;

pub async fn monit_whale()-> Result<(), Box<dyn std::error::Error>>{

    let mut ws_stream=match ws_connect().await{
            Ok(res)=>res,
            Err(e)=>{
                eprintln!("WebSocket 连接失败: {:?}", e);
                return Err(e);
            }
        };

        let subscribe_msg = r#"
            {
                "jsonrpc": "2.0",
                "id": 1,
                "method": "accountSubscribe",
                "params": [
                    "CM78CPUeXjn8o3yroDHxUtKsZZgoy4GPkPPXfouKNH12",
                    {
                    "encoding": "jsonParsed",
                    "commitment": "finalized"
                    }
                ]
            }
            "#;

        ws_stream.send(Message::Text(subscribe_msg.to_string())).await?;
        let (sender , receiver)=bounded::<SolEventTypes>(100);


        //对websocket接受到的消息进行预处理映射
        let ws_message_preprocess =tokio::spawn(async move||{
            streamMapORM::<SolEventTypes>(&mut ws_stream,sender).await;
        });

        //处理订阅的消息
        let handle_event=tokio::spawn(async move {
            while let Ok(event)=receiver.recv(){
                    match event{
                        SolEventTypes::LogNotification(ref log){

                        }
                    }
            }
        });


        match tokio::try_join!(
            ws_message_preprocess,
            handle_event
        ) {
            Ok(_) => println!("All tasks completed successfully"),
            Err(e) => eprintln!("A task exited with an error: {:?}", e),
        }

}