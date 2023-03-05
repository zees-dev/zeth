#![allow(unused)] // TODO remove for prod.

use axum::{
    extract::{Json, Path},
    http::StatusCode,
    response::{Html, IntoResponse, Response},
    routing::{get, get_service},
    Router,
};
use rust_embed::RustEmbed;
use serde_json::{json, Value};
use std::{
    collections::HashMap,
    net::SocketAddr,
    sync::{Arc, Mutex},
};
use surreal_simple_client::SurrealClient;
use tokio::sync::broadcast::{Receiver, Sender};
use tower_http::services::ServeDir;

mod endpoints;
use endpoints::{routes as endpoints_routes, service::EndpointsService};

mod surreal;

pub struct AppState {
    endpoint_service: EndpointsService,
    rpc_channel_map:
        Mutex<HashMap<String, (Sender<hyper::body::Bytes>, Receiver<hyper::body::Bytes>)>>,
    surreal_client: Arc<SurrealClient>,
}

impl AppState {
    async fn new(surreal_client: Arc<SurrealClient>) -> Self {
        // TODO inject surreal http client trait?

        let rpc_channel_map = Mutex::new(HashMap::<
            String,
            (Sender<hyper::body::Bytes>, Receiver<hyper::body::Bytes>),
        >::new());

        let endpoint_service = EndpointsService::new(surreal_client.clone());

        Self {
            endpoint_service,
            rpc_channel_map,
            surreal_client,
        }
    }
}

#[tokio::main]
async fn main() -> Result<(), String> {
    // TODO - use clap for CLI

    tracing_subscriber::fmt::init();
    tracing::info!("setting up client...");

    // TODO - get from config
    let (surreal_ws_rpc, admin_user, admin_pass, namespace, database) =
        ("ws://127.0.0.1:8000/rpc", "admin", "admin", "test", "test");

    let mut surreal_client =
        surreal::get_client(surreal_ws_rpc, admin_user, admin_pass, namespace, database)
            .await
            .map_err(|e| format!("failed to get surreal client: {}", e))?;
    surreal::setup_surreal_db(&mut surreal_client)
        .await
        .map_err(|e| format!("failed to setup surreal db: {}", e))?;

    let state = Arc::new(AppState::new(Arc::new(surreal_client)).await);

    let app = Router::new()
        // .route("/", get(index_handler))
        // .route("/:dir/:asset", get(asset_handler))
        .nest_service(
            "/",
            get_service(ServeDir::new("./client/dist")).handle_error(|_| not_found()),
        )
        .nest("/api/v1/endpoints", endpoints_routes::router(state.clone()))
        .route("/health", get(health))
        .route("/version", get(version))
        .fallback_service(get(not_found));
    // .with_state(state);

    let addr = std::net::SocketAddr::from(([0, 0, 0, 0], 3000)); // TODO - get from config, parse str instead
    tracing::info!("listening on {}...", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service_with_connect_info::<SocketAddr>())
        .await
        .unwrap();

    Ok(())
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
    tracing::info!("retreiving asset: {asset_path}");

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

// `curl -X GET http://localhost:3000/version`
async fn version() -> Json<Value> {
    Json(json!({ "version": "v0.0.1" })) // todo - get from env?/file?
}

// Finally, we use a fallback route for anything that didn't match.
async fn not_found() -> impl IntoResponse {
    (StatusCode::NOT_FOUND, Html("<h1>404</h1><p>Not Found</p>"))
}
