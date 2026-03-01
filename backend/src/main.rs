use salvo::prelude::*;

#[handler]
async fn hello() -> &'static str {
    "Hello World"
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt().init();
    let acceptor = TcpListener::new("0.0.0.0:8698").bind().await;
    let router = Router::new().get(hello);
    println!("{router:?}");
    Server::new(acceptor).serve(router).await;
}
