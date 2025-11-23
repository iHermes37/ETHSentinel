#[async_trait::async_trait]
pub trait StreamClient {
    type Request;
    type Response;

    async fn connect(&mut self) -> anyhow::Result<()>;     // 建立连接
    async fn subscribe(&mut self, req: Self::Request) -> anyhow::Result<()>; 
    async fn recv(&mut self) -> anyhow::Result<Self::Response>; // 统一接收
    async fn close(&mut self) -> anyhow::Result<()>;
}


pub trait JsonEvent: Sized {
    fn event_type(&self) -> String;
    fn deserialize_from_json(value: serde_json::Value) -> anyhow::Result<Self>;
}
