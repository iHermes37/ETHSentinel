use crossbeam_channel::Sender;
use tokio::net::TcpStream;
use tokio_tungstenite::{WebSocketStream, MaybeTlsStream};
use futures_util::{StreamExt};
use crate::models::websocket::websocket_event_type::WebsocketEventTypes;
use serde::ser::StdError;
use tokio_tungstenite::tungstenite::Message;
use serde_json;
use serde_json::Value;

pub async fn stream_maporm<T: WebsocketEventTypes + Send + 'static>(
    ws_stream:&mut WebSocketStream<MaybeTlsStream<TcpStream>>,
    sender: Sender<T>,
){
    while let Some(msg)=ws_stream.next().await{
        match msg{
            Ok(Message::Text(text)) => {
                // println!("收到消息: {}", text);
                if let Err(e) = process_text_message(text, &sender).await {
                    eprintln!("Failed to process text message: {:?}", e);
                }
            }
            _ => {}
        }
    }
}


async fn process_text_message<T:WebsocketEventTypes+ Send + 'static>(
    text: String,
    event_sender: &Sender<T>,
)-> Result<(), Box<dyn StdError>>{
    let event_json:Result<Value, _>=serde_json::from_str(&text);
    match event_json{
        Ok(events) => {
            eprintln!("解析消息: {}", text);
            process_json_events(events,event_sender)?;
        }
        Err(e) => {
            eprintln!("Error parsing JSON: {:?}", e);
        }
    }
    Ok(())
}

fn process_json_events<T : WebsocketEventTypes+ Send + 'static>(
    events: Value,
    event_sender: &Sender<T>,
)->Result<(), Box<dyn StdError>>{
    if events.is_array() {
        for event in events.as_array().unwrap() {
            process_single_event(event, event_sender)?;
        }
    } else {
        process_single_event(&events, event_sender)?;
    }
    Ok(())
}

fn process_single_event<T: WebsocketEventTypes + Send + 'static>(
    event: &Value,
    sender: &Sender<T>,
) -> Result<(), Box<dyn StdError>> {
    match T::deserialize_event(event) {
        Ok(event) => {
            sender.send(event).map_err(|e| e.into())
        }
        Err(e) => {
            // eprintln!("consume_stream.process_single_event: Error deserializing message: {:?}", e);
            Err(e.into())
        }
    }
}