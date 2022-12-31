use super::types;
use crate::AppState;

use axum::{
    body::Body,
    extract::{Json, Path, State},
    http::{Request, StatusCode},
    response::IntoResponse,
    routing::{get, post},
    Router,
};
use std::str::FromStr;
use std::sync::Arc;

pub fn router(state: Arc<AppState>) -> Router {
    Router::new()
        .route("/", get(get_endpoints).post(create_endpoint))
        .route(
            "/:endpoint_id",
            get(get_endpoint)
                .put(update_endpoint)
                .delete(delete_endpoint),
        )
        .route("/:endpoint_id/rpc", post(proxy_rpc_request))
        .with_state(state)
}

async fn get_endpoints(
    state: State<Arc<AppState>>,
) -> Result<impl IntoResponse, (StatusCode, String)> {
    let endpoints = state.endpoint_service.get_all().await;
    match endpoints {
        Ok(eps) => Ok(Json(eps)),
        Err(e) => {
            tracing::error!("[get_endpoints] error getting endpoints: {}", e);
            Err((StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))
        }
    }
}

async fn create_endpoint(
    state: State<Arc<AppState>>,
    Json(endpoint): Json<types::Endpoint>,
) -> Result<impl IntoResponse, (StatusCode, String)> {
    // TODO: can be optimized by using a single transaction
    // checking if name already exists, requires getting all endpoints (for now)
    let all_endpoints = state.endpoint_service.get_all().await;
    match all_endpoints {
        Ok(endpoints) => {
            // check if name (case-insensitive, trimmed) already exists
            if endpoints
                .iter()
                .any(|ep| ep.name.trim().eq_ignore_ascii_case(&endpoint.name.trim()))
            {
                return Err((
                    StatusCode::BAD_REQUEST,
                    "Endpoint with that name already exists".to_string(),
                ));
            }
        }
        Err(e) => {
            tracing::error!("[create_endpoint] error getting endpoints: {}", e);
            return Err((StatusCode::INTERNAL_SERVER_ERROR, e.to_string()));
        }
    }

    let endpoint_id = state.endpoint_service.create(&endpoint).await;
    match endpoint_id {
        Ok(id) => Ok(Json(id)),
        Err(e) => {
            tracing::error!("[create_endpoint] error creating endpoint: {}", e);
            Err((StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))
        }
    }
}

async fn get_endpoint(
    state: State<Arc<AppState>>,
    Path(endpoint_id): Path<String>,
) -> Result<impl IntoResponse, (StatusCode, String)> {
    let endpoint = state.endpoint_service.get(&endpoint_id).await;
    match endpoint {
        Ok(ep) => Ok(Json(ep)),
        Err(e) => {
            tracing::error!("[get_endpoint] error getting endpoint: {}", e);
            Err((StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))
        }
    }
}

async fn update_endpoint(
    state: State<Arc<AppState>>,
    Path(endpoint_id): Path<String>,
    Json(endpoint): Json<types::Endpoint>,
) -> Result<impl IntoResponse, (StatusCode, String)> {
    let endpoint = state.endpoint_service.update(&endpoint_id, &endpoint).await;
    match endpoint {
        Ok(ep) => Ok(Json(ep)),
        Err(e) => {
            tracing::error!("[update_endpoint] error updating endpoint: {}", e);
            Err((StatusCode::INTERNAL_SERVER_ERROR, e.to_string()))
        }
    }
}

async fn delete_endpoint(
    state: State<Arc<AppState>>,
    Path(endpoint_id): Path<String>,
) -> (StatusCode, String) {
    // TODO: can be optimized by using a single transaction
    // check if it exists first
    let result = state.endpoint_service.get(&endpoint_id).await;
    if let Some(err) = result.err() {
        tracing::error!("[delete_endpoint] not found: {}", err);
        return (StatusCode::NOT_FOUND, err.to_string());
    }

    let result = state.endpoint_service.delete(&endpoint_id).await;
    match result {
        Ok(_) => (StatusCode::OK, "deleted".to_string()),
        Err(e) => {
            tracing::error!("[delete_endpoint] error deleting endpoint: {}", e);
            (StatusCode::INTERNAL_SERVER_ERROR, e.to_string())
        }
    }
}

async fn proxy_rpc_request(
    state: State<Arc<AppState>>,
    Path(endpoint_id): Path<String>,
    mut req: Request<Body>,
) -> Result<impl IntoResponse, (StatusCode, String)> {
    let endpoint = state
        .endpoint_service
        .get(&endpoint_id)
        .await
        .map_err(|e| (StatusCode::NOT_FOUND, e.to_string()))?;

    // remove headers - which would break proxied request
    req.headers_mut().remove("host");
    req.headers_mut().remove("content-length");

    // reset url to proxy request
    let rpc_url = endpoint.rpc_http;
    *req.uri_mut() = hyper::http::Uri::from_str(&rpc_url).expect("invalid url");

    tracing::info!("request: {:?}", req);

    let connector = hyper_tls::HttpsConnector::new();
    let client = hyper::Client::builder().build::<_, hyper::Body>(connector);
    let res = client.request(req).await.expect("failed to make request");

    tracing::info!("response status: {}", res.status());

    Ok(res)
}

#[cfg(test)]
mod tests {
    use super::*;
    use hyper::Method;

    #[tokio::test]
    async fn test_rpc_using_hyper_client() {
        let connector = hyper_tls::HttpsConnector::new();
        let client = hyper::Client::builder().build::<_, hyper::Body>(connector);
        let url = "https://rpc.flashbots.net";
        let body = r#"{"jsonrpc":"2.0","id":1,"method":"eth_chainId"}"#;
        let req = hyper::http::Request::post(url)
            .header(
                hyper::http::header::CONTENT_TYPE,
                hyper::http::header::HeaderValue::from_static("application/json"),
            )
            .body(Body::from(body))
            .unwrap();

        // req=Request { method: POST, uri: https://rpc.flashbots.net, version: HTTP/1.1, headers: {"content-type": "application/json"}, body: Body(Full(b"{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_chainId\"}")) }

        let res = client.request(req).await.unwrap();

        assert_eq!(res.status(), StatusCode::OK);

        let bytes = hyper::body::to_bytes(res.into_body()).await.unwrap();
        let body_string = std::str::from_utf8(&bytes).unwrap();

        assert_eq!(
            body_string,
            "{\"id\":1,\"result\":\"0x1\",\"jsonrpc\":\"2.0\"}\n"
        );
    }

    #[tokio::test]
    async fn test_proxy_rpc_request() {
        let ds = surrealdb::Datastore::new("memory").await.unwrap();
        let state = Arc::new(AppState::new(ds));
        let endpoint_id = state
            .endpoint_service
            .create(&types::Endpoint {
                id: None,
                name: "test".to_string(),
                rpc_http: "https://rpc.flashbots.net".to_string(),
                is_dev: false,
                enabled: true,
                date_added: chrono::Utc::now(),
                explorer_url: None,
                rpc_ws: None,
            })
            .await
            .unwrap();

        let req = Request::builder()
            .method(Method::POST)
            .uri("https://rpc.flashbots.net")
            .header(
                hyper::http::header::CONTENT_TYPE,
                hyper::http::header::HeaderValue::from_static("application/json"),
            )
            .body(Body::from(
                r#"{"jsonrpc":"2.0","id":1,"method":"eth_chainId"}"#,
            ))
            .unwrap();

        let res = proxy_rpc_request(State(state), Path(endpoint_id), req)
            .await
            .unwrap()
            .into_response();

        assert_eq!(res.status(), StatusCode::OK);

        let bytes = hyper::body::to_bytes(res.into_body()).await.unwrap();
        let body_string = std::str::from_utf8(&bytes).unwrap();
        assert_eq!(
            body_string,
            "{\"id\":1,\"result\":\"0x1\",\"jsonrpc\":\"2.0\"}\n"
        );
    }
}
