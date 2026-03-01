use salvo::oapi::OpenApi;
use tonaris::build_router;

fn main() {
    let router = build_router();
    let doc = OpenApi::new("Tonaris API", "0.0.1").merge_router(&router);
    let json = serde_json::to_string_pretty(&doc).expect("Failed to serialize OpenAPI schema");

    let out = std::path::Path::new(env!("CARGO_MANIFEST_DIR")).join("../shared/openapi.json");
    std::fs::create_dir_all(out.parent().unwrap()).expect("Failed to create shared directory");
    std::fs::write(&out, json).expect("Failed to write schema file");

    println!("Schema written to {}", out.display());
}
