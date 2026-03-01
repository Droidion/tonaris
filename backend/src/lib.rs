use salvo::prelude::*;

#[endpoint]
async fn hello() -> &'static str {
    "Hello World"
}

pub fn build_router() -> Router {
    Router::new().push(Router::with_path("hello").get(hello))
}
