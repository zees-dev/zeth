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

#[derive(serde::Deserialize, Debug, Clone)]
struct SurrealResponse<T> {
    time: String,
    status: String,
    result: T,
}

#[derive(Clone)]
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

    // TODO: read setup from surreal ql file instead
    pub async fn setup_db(&self) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        let result = self.post::<SurrealResult>("INFO FOR DB;").await?;
        if !result.tb.contains_key("user") {
            tracing::info!("creating user table...");
            let data = r#"
                    BEGIN TRANSACTION;
                    DEFINE TABLE user SCHEMAFULL
                    PERMISSIONS
                        FOR select, update WHERE id = $auth.id,
                        FOR create, delete NONE;
                    DEFINE FIELD user ON user TYPE string ASSERT is::email($value);
                    DEFINE FIELD pass ON user TYPE string;
                    DEFINE FIELD tags ON user TYPE array;
                    DEFINE INDEX idx_user ON user COLUMNS user UNIQUE;
                    COMMIT TRANSACTION;
                "#;
            self.post(data).await?;
        }
        if !result.sc.contains_key("allusers") {
            tracing::info!("creating allusers scope...");
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
        if !result.tb.contains_key("endpoint") {
            tracing::info!("creating endpoint table...");
            let data = r#"
                    BEGIN TRANSACTION;
                    DEFINE TABLE endpoint SCHEMAFULL
                        PERMISSIONS 
                            FOR select, create, update, delete WHERE user = $auth.id;
                    DEFINE FIELD user ON endpoint TYPE record(user) VALUE $auth.id;
                    DEFINE FIELD name ON endpoint TYPE string VALUE string::trim($value) ASSERT string::length($value) > 2 AND string::length($value) < 65;
                    DEFINE FIELD enabled ON endpoint TYPE bool VALUE $value or true;
                    DEFINE FIELD date_added ON endpoint TYPE datetime VALUE time::now();
                    DEFINE FIELD rpc_url ON endpoint TYPE string ASSERT $value = /^(http|https|ws|wss):\/\/.+/ AND string::length($value) < 129;
                    -- TODO: get base db url from query (maybe different for every user)
                    DEFINE FIELD proxy_url ON endpoint TYPE string VALUE string::concat(session::origin(), '/', $this.id, '/rpc');
                    DEFINE FIELD symbol ON endpoint TYPE string ASSERT string::length($value) > 2 AND string::length($value) < 13;
                    DEFINE FIELD block_explorer_url ON endpoint TYPE string;
                    DEFINE INDEX idx_endpoint_name ON endpoint COLUMNS user, name UNIQUE;
                    DEFINE INDEX idx_endpoint_rpc_url ON endpoint COLUMNS user, rpc_url UNIQUE;
                    COMMIT TRANSACTION;
                "#;
            self.post(data).await?;
        }
        Ok(())
    }

    pub async fn post<T: serde::de::DeserializeOwned + Clone>(
        &self,
        data: &str,
    ) -> Result<T, Box<dyn std::error::Error + Send + Sync>> {
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

        // convert body to str
        // let body_str = std::str::from_utf8(&body)?.to_string();

        let obj: Vec<SurrealResponse<T>> = serde_json::from_slice(&body)?;
        let result = &obj
            .first()
            .ok_or("No response from SurrealQL server")?
            .result;
        Ok(result.to_owned())
    }
}

#[derive(serde::Deserialize, Debug, Clone)]
struct SurrealResult {
    dl: HashMap<String, String>,
    dt: HashMap<String, String>,
    sc: HashMap<String, String>,
    tb: HashMap<String, String>,
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

        let res = client.post::<SurrealResult>("INFO FOR DB;").await.unwrap();
        println!("{:?}", res);

        assert_eq!(1, 2);
    }
}
