use actix_web::{App, HttpServer, web};

#[actix_web::main]
async fn main() -> std::io::Result<()>{
    println!("Server running at http://127.0.0.1:8080");
    
    HttpServer::new(||{
        APP.new().configure()
    })
    .bind("127.0.0.1")
    .run
    .await
}