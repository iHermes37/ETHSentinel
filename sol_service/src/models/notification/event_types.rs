use serde::{Deserialize, Serialize};

use crate::models::notification::chainlog::LogsNotificationResponse;

#[derive(Serialize, Deserialize, Debug, Clone)]
pub enum SolEventTypes {
    LogNotification(LogsNotificationResponse),
    // AccountNotification(AccountNotificationResponse),
    // Program(ProgramNotificationResponse)
}