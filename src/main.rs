#![allow(unused)] // TODO remove for prod.

use anyhow::Result;
use axum::{
    extract::{Json, Path},
    http::StatusCode,
    response::{Html, IntoResponse, Response},
    routing::{get, get_service},
    Router,
};
use rust_embed::RustEmbed;
use serde_json::{json, Value};
use std::{net::SocketAddr, sync::Arc};
use surrealdb::Datastore;
use tower_http::services::ServeDir;

pub mod db;

mod endpoints;
use endpoints::{routes as endpoints_routes, service::EndpointsService};

pub struct AppState {
    endpoint_service: EndpointsService,
    channel: (
        crossbeam_channel::Sender<()>,
        crossbeam_channel::Receiver<()>,
    ),
}

impl AppState {
    fn new(ds: Datastore) -> Self {
        // TODO inject DB trait?

        let endpoint_service = EndpointsService::new(ds);

        let (tx, rx) = crossbeam_channel::unbounded::<()>();
        Self {
            endpoint_service,
            channel: (tx, rx),
        }
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    // TODO - use clap for CLI

    tracing_subscriber::fmt::init();

    // docker run --rm -it --name surrealdb -p 127.0.0.1:8000:8000 surrealdb/surrealdb:latest start --log trace --user root --pass root memory
    // let ds = Datastore::new("memory").await?;
    let ds = Datastore::new("file://Zeth/temp.db").await?;

    let state = Arc::new(AppState::new(ds));

    let app = Router::new()
        // .route("/", get(index_handler))
        // .route("/:dir/:asset", get(asset_handler))
        .nest_service(
            "/",
            get_service(ServeDir::new("./client/dist")).handle_error(|_| not_found()),
        )
        .nest("/api/v1/endpoints", endpoints_routes::router(state.clone()))
        // .route("/:endpoint_id/rpc", any_service(service_fn(proxy_rpc_req)))
        // .route(
        //     "/:endpoint_id/rpc",
        //     any_service(service_fn(move |req: Request<Body>| {
        //         // let router_svc = router_svc.clone();
        //         async move {
        //             // if req.method() == Method::CONNECT {
        //             proxy(req).await
        //             // } else {
        //             //     router_svc.oneshot(req).await.map_err(|err| match err {})
        //             // }
        //         }
        //     })),
        // )
        .route("/health", get(health))
        .route("/version", get(version))
        .fallback_service(get(not_found));
    // .with_state(state);

    let addr = std::net::SocketAddr::from(([0, 0, 0, 0], 3000)); // TODO - get from config
    tracing::info!("listening on {}...", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service_with_connect_info::<SocketAddr>())
        .await?;

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
