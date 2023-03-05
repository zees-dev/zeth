use axum::http::request::Builder;
use hyper::body::Buf;
use hyper::client::HttpConnector;
use hyper::service::Service;
use hyper::{Body, Client, HeaderMap, Request};
use hyper_tls::HttpsConnector;
use serde::Deserialize;
use serde_json::Value;
use std::collections::HashMap;
use std::str::FromStr;
use surreal_simple_client::SurrealClient;
use tokio::net::TcpStream;

pub async fn get_client<'a>(
    surreal_ws_rpc: &'a str,
    admin_user: &'a str,
    admin_pass: &'a str,
    namespace: &'a str,
    database: &'a str,
) -> Result<SurrealClient, anyhow::Error> {
    let mut surreal_client = SurrealClient::new(surreal_ws_rpc)
        .await
        .map_err(|e| anyhow::anyhow!(e))?;

    surreal_client
        .signin(admin_user, admin_pass)
        .await
        .map_err(|e| anyhow::anyhow!(e))?;

    surreal_client
        .use_namespace(namespace, database)
        .await
        .map_err(|e| anyhow::anyhow!(e))?;

    Ok(surreal_client)
}

#[derive(Deserialize, Debug)]
struct DBInfo {
    dl: HashMap<String, String>,
    dt: HashMap<String, String>,
    sc: HashMap<String, String>,
    tb: HashMap<String, String>,
}

pub async fn setup_surreal_db(
    surreal_client: &mut SurrealClient,
) -> Result<(), Box<dyn std::error::Error>> {
    // client
    //     .send_query(
    //         "create User set username = $username".to_owned(),
    //         json!({ "username": "John" }),
    //     )
    //     .await
    //     .unwrap();

    // let some_user: Option<User> = client
    //     .find_one("select * from User".to_owned(), Value::Null)
    //     .await
    //     .unwrap();

    // if let Some(user) = some_user {
    //     print!("found user: {:?}", user);
    // }

    // let info_res = self.post("INFO FOR DB;").await?;
    let info_res = surreal_client
        .send_query("INFO FOR DB;".into(), Value::Null)
        .await
        .expect("as");
    // .await
    // .map_err(|e| anyhow::anyhow!(e))?;

    println!("info_res: {:?}", info_res);

    // if let Some(result) = info_res.result {
    //     if !result.tb.contains_key("user") {
    //         let data = r#"
    //                 BEGIN TRANSACTION;
    //                 DEFINE TABLE user SCHEMAFULL
    //                 PERMISSIONS
    //                     FOR select, update WHERE id = $auth.id,
    //                     FOR create, delete NONE;
    //                 DEFINE FIELD user ON user TYPE string ASSERT is::email($value);
    //                 DEFINE FIELD pass ON user TYPE string;
    //                 DEFINE FIELD tags ON user TYPE array;
    //                 DEFINE INDEX idx_user ON user COLUMNS user UNIQUE;
    //                 COMMIT TRANSACTION;
    //             "#;
    //         self.post(data).await?;
    //     }
    //     if !result.sc.contains_key("allusers") {
    //         let data = r#"
    //                 BEGIN TRANSACTION;
    //                 DEFINE SCOPE allusers
    //                 -- the JWT session will be valid for 14 days
    //                 SESSION 14d
    //                 -- The optional SIGNUP clause will be run when calling the signup method for this scope
    //                 -- It is designed to create or add a new record to the database.
    //                 -- If set, it needs to return a record or a record id
    //                 -- The variables can be passed in to the signin method
    //                 SIGNUP ( CREATE user SET user = string::lowercase(string::trim($user)), pass = crypto::argon2::generate($pass), tags = $tags )
    //                 -- The optional SIGNIN clause will be run when calling the signin method for this scope
    //                 -- It is designed to check if a record exists in the database.
    //                 -- If set, it needs to return a record or a record id
    //                 -- The variables can be passed in to the signin method
    //                 SIGNIN ( SELECT * FROM user WHERE user = string::lowercase(string::trim($user)) AND crypto::argon2::compare(pass, $pass) );
    //                 -- this optional clause will be run when calling the signup method for this scope
    //                 COMMIT TRANSACTION;
    //             "#;
    //         self.post(data).await?;
    //     }
    //     if !result.tb.contains_key("endpoint") {
    //         let data = r#"
    //                 BEGIN TRANSACTION;
    //                 DEFINE TABLE endpoint SCHEMAFULL
    //                     PERMISSIONS
    //                         FOR select, create, update, delete WHERE user = $auth.id;
    //                 DEFINE FIELD user ON endpoint TYPE record(user);
    //                 DEFINE FIELD name ON endpoint TYPE string;
    //                 DEFINE FIELD enabled ON endpoint TYPE bool;
    //                 DEFINE FIELD date_added ON endpoint TYPE datetime;
    //                 DEFINE FIELD rpc_url ON endpoint TYPE string;
    //                 DEFINE FIELD proxy_url ON endpoint TYPE string;
    //                 -- TODO: ensure that the proxy_url cannot be updated by user
    //                 DEFINE FIELD type ON endpoint TYPE string;
    //                 DEFINE FIELD symbol ON endpoint TYPE string;
    //                 DEFINE FIELD block_explorer_url ON endpoint TYPE string;
    //                 DEFINE INDEX idx_endpoint_name ON endpoint COLUMNS user, name UNIQUE;
    //                 DEFINE INDEX idx_endpoint_rpc_url ON endpoint COLUMNS user, rpc_url UNIQUE;
    //                 -- TODO: define event for creation
    //                 -- TODO: define event for update
    //                 -- TODO: define event for deletion
    //                 COMMIT TRANSACTION;
    //             "#;
    //         self.post(data).await?;
    //     }
    // } else {
    //     Err("Failed to get info for DB".to_string())?;
    // }
    Ok(())
}
