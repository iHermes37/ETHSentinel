use crossbeam_channel::Sender;
use tokio::net::TcpStream;
use tokio_tungstenite::{WebSocketStream, MaybeTlsStream};
use futures_util::{StreamExt};
use crate::Message;
use crate::subscription::EventTypes::WebsocketEventTypes;
use serde::ser::StdError;

pub async fn streamMapORM<T: WebsocketEventTypes + Send + 'static>(
    ws_stream:&mut WebSocketStream<MaybeTlsStream<TcpStream>>,
    sender: Sender<T>,
){
    while let Some(msg)=ws_stream.next().await{
        match msg?{
            Ok(Message::Text(text)) => {
                println!("收到消息: {}", text);
                if let Err(e) = process_text_message(text, &sender).await {
                    eprintln!("Failed to process text message: {:?}", e);
                }
            }
            _ => {}
        }
    }
}


async fn process_text_message<T : EventTypes+ Send + 'static>(
    text: String,
    eventSender: &Sender<T>,
)-> Result<(), Box<dyn StdError>>{
    
}

fn process_json_events(){

}

fn process_single_events(){
    
}