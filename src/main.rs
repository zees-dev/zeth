// #![allow(unused)] // For beginning only.
use axum::{
    extract::{Json, Path, State},
    http::StatusCode,
    response::{Html, IntoResponse, Response},
    routing::{get, get_service, post, put},
    Router,
};
use rust_embed::RustEmbed;
use serde_json::{json, Value};
use tower_http::services::ServeDir;

#[derive(Clone)]
struct AppState {
    tx: crossbeam_channel::Sender<()>,
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    let (tx, rx) = crossbeam_channel::unbounded::<()>();

    let state = AppState { tx: tx.clone() };

    let app = Router::new()
        // .route("/", get(index_handler))
        // .route("/:dir/:asset", get(asset_handler))
        .nest_service(
            "/",
            get_service(ServeDir::new("./client/dist")).handle_error(|_| not_found()),
        )
        .route("/health", get(health))
        .fallback_service(get(not_found))
        .with_state(state);

    let addr = std::net::SocketAddr::from(([0, 0, 0, 0], 3000)); // TODO - get from config
    tracing::info!("listening on {}...", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

#[derive(RustEmbed)]
#[folder = "client/dist/"]
struct Assets;

// `curl -X GET http://localhost:3000/`
// TODO - use this instead of ServeDir - to serve single binary - with SPA and all assets
async fn index_handler() -> impl IntoResponse {
    let file = Assets::get("index.html").unwrap();
    let body = axum::body::boxed(axum::body::Full::from(file.data));

    Response::builder()
        .status(StatusCode::OK)
        .header("content-type", "text/html")
        .body(body)
        .unwrap()
}

async fn asset_handler(Path((dir, asset)): Path<(String, String)>) -> impl IntoResponse {
    let asset_path = format!("{}/{}", dir, asset);
    tracing::info!("retreiving asset: {}", asset_path);

    let file = Assets::get(&asset_path).unwrap();
    let body = axum::body::boxed(axum::body::Full::from(file.data));

    Response::builder()
        .status(StatusCode::OK)
        .header("content-type", "image/svg+xml") // <-- TODO - get from file name
        .body(body)
        .unwrap()
}

// `curl -X GET http://localhost:3000/health`
async fn health() -> Json<Value> {
    Json(json!({ "status": "up" }))
}

// Finally, we use a fallback route for anything that didn't match.
async fn not_found() -> impl IntoResponse {
    (StatusCode::NOT_FOUND, Html("<h1>404</h1><p>Not Found</p>"))
}
