use salvo::prelude::*;
use salvo_cors::Cors;
use tonaris::build_router;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt().init();
    let cors = Cors::permissive().into_handler();
    let router = build_router().hoop(cors);
    let doc = OpenApi::new("Tonaris API", "0.0.1").merge_router(&router);
    let router = router
        .unshift(doc.into_router("/api-doc/openapi.json"))
        .unshift(Scalar::new("/api-doc/openapi.json").into_router("/scalar"));
    let port = std::env::var("PORT").unwrap_or_else(|_| "8698".to_string());
    let acceptor = TcpListener::new(format!("0.0.0.0:{port}")).bind().await;
    Server::new(acceptor).serve(router).await;
}
