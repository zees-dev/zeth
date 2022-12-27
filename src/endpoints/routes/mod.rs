use std::sync::Arc;

use crate::AppState;
use axum::{
    extract::{Json, Path, State},
    http::StatusCode,
    response::{Html, IntoResponse, Response},
    routing::{delete, get, get_service, post, put},
    Router,
};
use serde_json::json;

use super::types;

pub fn router(state: Arc<AppState>) -> Router {
    Router::new()
        .route("/", get(get_endpoints).post(create_endpoint))
        .route(
            "/:endpoint_id",
            get(get_endpoint)
                .put(update_endpoint)
                .delete(delete_endpoint),
        )
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
