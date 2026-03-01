use salvo::prelude::*;

#[endpoint]
async fn hello() -> &'static str {
    "Hello World"
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt().init();
    let router = Router::new().push(Router::with_path("hello").get(hello));
    let doc = OpenApi::new("test api", "0.0.1").merge_router(&router);
    let router = router
        .unshift(doc.into_router("/api-doc/openapi.json"))
        .unshift(Scalar::new("/api-doc/openapi.json").into_router("/scalar"));
    let acceptor = TcpListener::new("0.0.0.0:8698").bind().await;
    println!("{router:?}");
    Server::new(acceptor).serve(router).await;
}
