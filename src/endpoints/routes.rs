use crate::AppState;

use axum::{
    body::Body,
    extract::{
        ws::{Message, WebSocketUpgrade},
        Json, Path, State,
    },
    headers::UserAgent,
    http::{Request, StatusCode},
    response::{
        sse::{Event, Sse},
        IntoResponse,
    },
    routing::get,
    Router, TypedHeader,
};
use std::str::FromStr;
use std::sync::Arc;
use tokio_tungstenite::tungstenite;

pub fn router(state: Arc<AppState>) -> Router {
    Router::new()
        // .route("/", get(get_endpoints).post(create_endpoint))
        // .route(
        //     "/:endpoint_id",
        //     get(get_endpoint)
        //         .put(update_endpoint)
        //         .delete(delete_endpoint),
        // )
        .route(
            "/:endpoint_id/rpc",
            get(proxy_ws_rpc_request).post(proxy_http_rpc_request),
        )
        .route("/:endpoint_id/rpc/events", get(sse_handler))
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

async fn proxy_ws_rpc_request(
    state: State<Arc<AppState>>,
    ws: WebSocketUpgrade,
    user_agent: Option<TypedHeader<UserAgent>>,
    Path(endpoint_id): Path<String>,
) -> Result<impl IntoResponse, (StatusCode, String)> {
    let endpoint = state
        .endpoint_service
        .get(&endpoint_id)
        .await
        .map_err(|e| (StatusCode::NOT_FOUND, e.to_string()))?;

    let ws_url = match endpoint.rpc_ws {
        Some(url) => url,
        None => {
            return Err((
                StatusCode::BAD_REQUEST,
                "Endpoint does not support WS RPC".to_string(),
            ))
        }
    };

    let mut rpc_channel_map = state.rpc_channel_map.lock().unwrap();
    let sender = match rpc_channel_map.get(&endpoint_id) {
        Some((sender, _)) => sender.clone(),
        None => {
            let (sender, receiver) = tokio::sync::broadcast::channel::<hyper::body::Bytes>(2);
            rpc_channel_map.insert(endpoint_id.clone(), (sender.clone(), receiver));
            tracing::info!("created new event publisher for {}...", endpoint_id);
            sender
        }
    };

    // TODO: gracefully handle error conditions (dont unwrap)
    Ok(ws.on_upgrade(|mut socket| async move {
        let (mut endpoint_ws_stream, _) = tungstenite::connect(ws_url).unwrap();

        if let Some(TypedHeader(user_agent)) = user_agent {
            tracing::info!("`{}` connected", user_agent.as_str());
        }

        loop {
            // TODO: need some sort of keep-alive check to make sure the client is still connected
            if let Some(msg) = socket.recv().await {
                if let Ok(msg) = msg {
                    match msg {
                        Message::Text(t) => {
                            tracing::debug!("client sent str: {:?}", t);
                            endpoint_ws_stream
                                .write_message(tokio_tungstenite::tungstenite::Message::Text(t))
                                .unwrap();
                        }
                        Message::Binary(b) => {
                            tracing::debug!(
                                "client sent binary data, (utf8 string): {:?}",
                                String::from_utf8(b.clone()).unwrap()
                            );
                            endpoint_ws_stream
                                .write_message(tokio_tungstenite::tungstenite::Message::Binary(b))
                                .unwrap();
                        }
                        Message::Ping(_) | Message::Pong(_) => {
                            tracing::debug!("socket ping-pong");
                            socket.send(Message::Pong(Vec::new())).await.unwrap();
                        }
                        Message::Close(_) => {
                            tracing::debug!("disconnecting client...");
                            socket.close().await;
                            break;
                        }
                    }
                } else {
                    socket.close().await;
                    endpoint_ws_stream.close(None).unwrap();
                    break;
                }
            }
            if let Ok(msg) = endpoint_ws_stream.read_message() {
                match msg {
                    tokio_tungstenite::tungstenite::Message::Text(t) => {
                        tracing::debug!("endpoint sent str: {:?}", t);
                        socket.send(Message::Text(t.clone())).await.unwrap();
                        sender.send(hyper::body::Bytes::from(t)).unwrap();
                    }
                    tokio_tungstenite::tungstenite::Message::Binary(b) => {
                        tracing::debug!("endpoint sent binary data: {:?}", b);
                        socket.send(Message::Binary(b.clone())).await.unwrap();
                        sender.send(hyper::body::Bytes::from(b)).unwrap();
                    }
                    tokio_tungstenite::tungstenite::Message::Ping(_)
                    | tokio_tungstenite::tungstenite::Message::Pong(_) => {
                        tracing::debug!("endpoint ping-pong");
                        socket.send(Message::Pong(Vec::new())).await.unwrap();
                    }
                    tokio_tungstenite::tungstenite::Message::Close(_) => {
                        tracing::debug!("disconnecting endpoint...");
                        socket.close().await.unwrap();
                        return;
                    }
                    _ => {}
                }
            }
        }
        tracing::info!("client disconnected.");
    }))
}

async fn proxy_http_rpc_request(
    state: State<Arc<AppState>>,
    Path(endpoint_id): Path<String>,
    mut req: Request<Body>,
) -> Result<impl IntoResponse, (StatusCode, String)> {
    let endpoint = state
        .endpoint_service
        .get(&endpoint_id)
        .await
        .map_err(|e| (StatusCode::NOT_FOUND, e.to_string()))?;

    // req.into_parts();

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

    let body_bytes = hyper::body::to_bytes(res.into_body())
        .await
        .expect("failed to read body");

    let mut rpc_channel_map = state.rpc_channel_map.lock().unwrap();
    match rpc_channel_map.get(&endpoint_id) {
        Some((sender, _)) => {
            // send raw response body to existing channel
            sender
                .send(body_bytes.clone())
                .expect("failed to send message");
        }
        None => {
            let channel = tokio::sync::broadcast::channel::<hyper::body::Bytes>(2);
            rpc_channel_map.insert(endpoint_id.clone(), channel);
            tracing::info!("created new event publisher for {}...", endpoint_id);
        }
    }
    // TODO: handle cleanup/removal of channel (LRU cache?)

    Ok(body_bytes)
}

async fn sse_handler(
    state: State<Arc<AppState>>,
    TypedHeader(user_agent): TypedHeader<UserAgent>,
    Path(endpoint_id): Path<String>,
) -> Result<Sse<impl futures::Stream<Item = Result<Event, anyhow::Error>>>, (StatusCode, String)> {
    state
        .endpoint_service
        .get(&endpoint_id)
        .await
        .map_err(|e| (StatusCode::NOT_FOUND, e.to_string()))?;

    let rpc_channel_map = state.rpc_channel_map.lock().unwrap();
    let receiver = match rpc_channel_map.get(&endpoint_id) {
        Some((sender, _)) => sender.subscribe(),
        None => {
            return Err((
                StatusCode::NOT_FOUND,
                format!("no event publisher found for {}", endpoint_id),
            ))
        }
    };

    tracing::info!("`{}` connected", user_agent.as_str());

    let stream = futures::stream::unfold(receiver, |mut receiver| async move {
        match receiver.recv().await {
            Ok(bytes) => {
                let event = Event::default().data(String::from_utf8(bytes.to_vec()).unwrap());
                Some((Ok(event), receiver))
            }
            Err(err) => {
                tracing::error!("error receiving message: {}", err);
                None
            }
        }
    });

    Ok(Sse::new(stream))
}

#[cfg(test)]
mod tests {
    use crate::surreal::SurrealHttpClient;

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
    async fn test_proxy_http_rpc_request() {
        let ds = surrealdb::Datastore::new("memory").await.unwrap();
        let client = SurrealHttpClient::new(
            "http://localhost:8000/sql",
            "admin",
            "admin",
            "test",
            "test",
        );
        let state = Arc::new(AppState::new(ds, client).await);
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

        let res = proxy_http_rpc_request(State(state), Path(endpoint_id), req)
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

    #[tokio::test]
    async fn test_channels() {
        let (tx, mut rx1) = tokio::sync::broadcast::channel(2);
        let mut rx2 = tx.subscribe();

        tokio::spawn(async move {
            assert_eq!(rx1.recv().await.unwrap(), 10);
            assert_eq!(rx1.recv().await.unwrap(), 20);
        });

        tokio::spawn(async move {
            assert_eq!(rx2.recv().await.unwrap(), 110);
            assert_eq!(rx2.recv().await.unwrap(), 20);
        });

        tx.send(10).unwrap();
        tx.send(20).unwrap();
    }
}
