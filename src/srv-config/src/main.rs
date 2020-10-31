#![deny(warnings)]

#[tokio::main]
async fn main() {
    srv_config::run().await;
}
