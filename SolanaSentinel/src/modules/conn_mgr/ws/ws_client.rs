use serde_json::Value;
use std::error::Error;
use crate::models::notification::event_types::SolEventTypes;
use crate::models::notification::chainlog::LogsNotificationResponse;
use crate::utils::serde_helper::deserialize_into;


pub enum WsClient{
    LogNotification(LogsNotificationResponse),
    // AccountNotification(AccountNotificationResponse),
    // Program(ProgramNotificationResponse)
}


impl StreamClient for WsClient{

    fn event_type(&self) -> String{
        match self{
            SolEventTypes::LogNotification(_)=>"LogNotification".to_string()
        }
    }

    fn deserialize_event(value: &Value) -> Result<Self, Box<dyn Error>>{
            let method_name=value["method"].as_str().ok_or_else(||"Missing method in event")?;
            let result=match method_name{
                    "logsNotification" => {
                        let log_notification=deserialize_into::<LogsNotificationResponse>(value)?;
                        Ok(SolEventTypes::LogNotification(log_notification))
                    }
                    _=>Err(format!("Unsupported event type: {}", method_name).into())
            };
            result
    }  
}