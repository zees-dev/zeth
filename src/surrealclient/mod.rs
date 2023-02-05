use axum::http::request::Builder;
use hyper::body::Buf;
use hyper::client::HttpConnector;
use hyper::service::Service;
use hyper::{Body, Client, HeaderMap, Request};
use hyper_tls::HttpsConnector;
use serde::Deserialize;
use std::collections::HashMap;
use std::str::FromStr;
use tokio::net::TcpStream;

pub struct SurrealHttpClient {
    url: String,
    username: String,
    password: String,
    namespace: String,
    database: String,
    client: Client<HttpsConnector<HttpConnector>>,
}

impl SurrealHttpClient {
    pub fn new(
        url: &str,
        username: &str,
        password: &str,
        namespace: &str,
        database: &str,
    ) -> SurrealHttpClient {
        let connector = hyper_tls::HttpsConnector::new();
        let client = hyper::Client::builder().build::<_, hyper::Body>(connector);

        SurrealHttpClient {
            url: url.to_string(),
            username: username.to_string(),
            password: password.to_string(),
            namespace: namespace.to_string(),
            database: database.to_string(),
            client,
        }
    }

    pub async fn setup_db(&self) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        let info_res = self.post("INFO FOR DB;").await?;
        if let Some(result) = info_res.result {
            if !result.tb.contains_key("user") {
                let data = r#"
                    BEGIN TRANSACTION;
                    DEFINE TABLE user SCHEMAFULL
                    PERMISSIONS
                        FOR select, update WHERE id = $auth.id,
                        FOR create, delete NONE;
                    DEFINE FIELD user ON user TYPE string;
                    DEFINE FIELD pass ON user TYPE string;
                    DEFINE FIELD tags ON user TYPE array;
                    DEFINE INDEX idx_user ON user COLUMNS user UNIQUE;
                    COMMIT TRANSACTION;
                "#;
                self.post(data).await?;
            }
            if !result.sc.contains_key("allusers") {
                let data = r#"
                    BEGIN TRANSACTION;
                    DEFINE SCOPE allusers
                    -- the JWT session will be valid for 14 days
                    SESSION 14d
                    -- The optional SIGNUP clause will be run when calling the signup method for this scope
                    -- It is designed to create or add a new record to the database.
                    -- If set, it needs to return a record or a record id
                    -- The variables can be passed in to the signin method
                    SIGNUP ( CREATE user SET user = string::lowercase(string::trim($user)), pass = crypto::argon2::generate($pass), tags = $tags )
                    -- The optional SIGNIN clause will be run when calling the signin method for this scope
                    -- It is designed to check if a record exists in the database.
                    -- If set, it needs to return a record or a record id
                    -- The variables can be passed in to the signin method
                    SIGNIN ( SELECT * FROM user WHERE user = string::lowercase(string::trim($user)) AND crypto::argon2::compare(pass, $pass) );
                    -- this optional clause will be run when calling the signup method for this scope
                    COMMIT TRANSACTION;
                "#;
                self.post(data).await?;
            }
            if !result.tb.contains_key("endpoints") {
                let data = r#"
                    BEGIN TRANSACTION;
                    DEFINE TABLE endpoints SCHEMAFULL
                        PERMISSIONS 
                            FOR select, update, delete WHERE id = $auth.id,
                            FOR create FULL;
                    DEFINE FIELD name ON endpoints TYPE string;
                    DEFINE FIELD enabled ON endpoints TYPE bool;
                    DEFINE FIELD date_added ON endpoints TYPE datetime;
                    DEFINE FIELD rpc_url ON endpoints TYPE string;
                    DEFINE FIELD type ON endpoints TYPE string;
                    DEFINE INDEX idx_endpoints_name ON endpoints COLUMNS name UNIQUE;
                    DEFINE INDEX idx_endpoints_rpc_url ON endpoints COLUMNS rpc_url UNIQUE;
                    -- TODO: define event for creation
                    -- TODO: define event for update
                    -- TODO: define event for deletion
                    COMMIT TRANSACTION;
                "#;
                self.post(data).await?;
            }
        } else {
            Err("Failed to get info for DB".to_string())?;
        }
        Ok(())
    }

    pub async fn post(
        &self,
        data: &str,
    ) -> Result<SurrealResponse, Box<dyn std::error::Error + Send + Sync>> {
        let req = Request::builder()
            .method("POST")
            .uri(&self.url)
            .header("Accept", "application/json")
            .header("NS", &self.namespace)
            .header("DB", &self.database)
            .header(
                "Authorization",
                format!(
                    "Basic {}",
                    base64::encode(format!("{}:{}", self.username, self.password))
                ),
            )
            .body(Body::from(data.to_owned()))
            .unwrap();

        let res = self.client.request(req).await?;
        let body = hyper::body::to_bytes(res.into_body()).await?;
        let res = serde_json::from_slice::<Vec<SurrealResponse>>(&body)?
            .first()
            .unwrap()
            .clone();
        Ok(res)
    }
}

#[derive(serde::Deserialize, Debug, Clone)]
struct SurrealResult {
    dl: HashMap<String, String>,
    dt: HashMap<String, String>,
    sc: HashMap<String, String>,
    tb: HashMap<String, String>,
}

#[derive(serde::Deserialize, Debug, Clone)]
pub struct SurrealResponse {
    time: String,
    status: String,
    result: Option<SurrealResult>,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_info_request() {
        let client = SurrealHttpClient::new(
            "http://localhost:8000/sql",
            "admin",
            "admin",
            "test",
            "test",
        );
        client.setup_db().await.unwrap();

        let res = client.post("INFO FOR DB;").await.unwrap();
        println!("{:?}", res);

        assert_eq!(1, 2);
    }
}
