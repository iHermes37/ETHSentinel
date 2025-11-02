use crossbeam_channel::{bounded};
use tokio_tungstenite::{
    tungstenite::protocol::Message
};
use futures_util::{SinkExt, StreamExt};

use crate::initialize::connect::ws_connect; 
use crate::models::notification::event_types::SolEventTypes;
use crate::modules::subscription::ws_stream_trans::stream_maporm;

pub async fn track_main()-> Result<(), Box<dyn std::error::Error>>{

    eprintln!("进入追踪");

    let mut ws_stream=match ws_connect().await{
            Ok(res) => {
                 println!("WebSocket 连接成功");
                 res
            },
            Err(e) => {
                eprintln!("WebSocket 连接失败: {:?}", e);
                return Err(e);
            }
        };

        let subscribe_msg = r#"
            {
                        "jsonrpc": "2.0",
                        "id": 1,
                        "method": "logsSubscribe",
                        "params": [
                            {
                                "mentions": log_program_ids
                            },
                            {
                                "encoding": "jsonParsed",
                                "commitment": "finalized"
                            }
                        ]
                    }
            "#;

        ws_stream.send(Message::Text(subscribe_msg.to_string().into())).await?;

        // 接收消息
        while let Some(msg) = ws_stream.next().await {
            match msg {
                Ok(Message::Text(text)) => println!("收到消息: {}", text),
                Ok(_) => {},
                Err(e) => eprintln!("错误: {:?}", e),
            }
        }

        let (sender , receiver)=bounded::<SolEventTypes>(100);


        //对websocket接受到的消息进行预处理映射
        let ws_message_preprocess =tokio::spawn(async move{
            stream_maporm::<SolEventTypes>(&mut ws_stream,sender).await;
        });

        //处理订阅的消息
        // let handle_event=tokio::spawn(async move {
        //     while let Ok(event)=receiver.recv(){
        //             match event{

        //                 SolEventTypes::LogNotification(ref log)=>{
        //                         println!("LogNotification")
        //                 }

        //                 // SolEventTypes::AccountNotification(notification) => {
        //                 //     let signature = notification;
        //                 //     println!("[[SOLANA TASK]] GOT ACCOUNT NOTIFICATION {:?}", signature)
        //                 // }
                            
        //                 _ => {
        //                 println!("Stand by")
        //                 }
        //             }
        //         }
        // });


        match tokio::try_join!(
            ws_message_preprocess,
            // handle_event
        ) {
            Ok(_) => {
                println!("All tasks completed successfully");
            },
            Err(e) => {
                eprintln!("A task exited with an error: {:?}", e);
            },
        }

         Ok(())

}