use anyhow::anyhow;
use chrono::TimeZone;
use chrono::Utc;
use std::collections::BTreeMap;
use surrealdb::sql::{Object, Value};

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct Endpoint {
    pub id: Option<String>,
    pub name: String,
    pub is_dev: bool,
    pub enabled: bool,
    pub date_added: chrono::DateTime<chrono::Utc>,
    pub explorer_url: Option<String>,
    pub rpc_http: String,
    pub rpc_ws: Option<String>,
}

impl Endpoint {
    pub fn validate(&self) -> Result<bool, String> {
        // TODO: Validate the explorer_url

        // TODO: Validate the rpc_http
        // TODO: Validate the rpc_ws
        Ok(true)
    }

    pub fn test_connection(&self) -> Result<bool, String> {
        // TODO: Test the connection via http HEAD request to the rpc_http endpoint

        // TODO: Test the connection via http HEAD request to the rpc_ws endpoint (ws-connect?)?
        Ok(true)
    }
}

// impl From<Endpoint> for BTreeMap<String, Value> {
//     fn from(ep: Endpoint) -> Self {
//         let mut map = BTreeMap::new();
//         if ep.id.is_some() {
//             map.insert("id".into(), ep.id.unwrap().into());
//         }
//         map.insert("name".into(), ep.name.into());
//         map.insert("is_dev".into(), ep.is_dev.into());
//         map.insert("enabled".into(), ep.enabled.into());
//         map.insert("date_added".into(), ep.date_added.into());
//         if ep.explorer_url.is_some() {
//             map.insert("explorer_url".into(), ep.explorer_url.unwrap().into());
//         }
//         map.insert("rpc_http".into(), ep.rpc_http.into());
//         if ep.rpc_ws.is_some() {
//             map.insert("rpc_ws".into(), ep.rpc_ws.unwrap().into());
//         }
//         map
//     }
// }

impl From<Endpoint> for Value {
    fn from(ep: Endpoint) -> Self {
        let mut map = BTreeMap::new();

        if ep.id.is_some() {
            map.insert("id".to_string(), ep.id.unwrap().into());
        }
        map.insert("name".to_string(), ep.name.into());
        map.insert("is_dev".to_string(), ep.is_dev.into());
        map.insert("enabled".to_string(), ep.enabled.into());
        map.insert("date_added".to_string(), ep.date_added.into());
        if ep.explorer_url.is_some() {
            map.insert("explorer_url".to_string(), ep.explorer_url.unwrap().into());
        }
        map.insert("rpc_http".to_string(), ep.rpc_http.into());
        if ep.rpc_ws.is_some() {
            map.insert("rpc_ws".to_string(), ep.rpc_ws.unwrap().into());
        }

        Value::Object(map.into())
    }
}

impl TryFrom<Object> for Endpoint {
    type Error = anyhow::Error;
    fn try_from(mut val: Object) -> Result<Endpoint, anyhow::Error> {
        let map = val.0;

        // Extract the values from the map and convert them to the appropriate types
        let id = map.get("id").map(|v| v.to_owned().as_string());
        let name = map
            .get("name")
            .map(|v| v.to_owned().as_string())
            .ok_or(anyhow!("missing name"))?;
        let is_dev = map
            .get("is_dev")
            .map(|v| v.to_owned().as_string() == "true")
            .ok_or(anyhow!("missing is_dev"))?;
        let enabled = map
            .get("enabled")
            .map(|v| v.to_owned().as_string() == "true")
            .ok_or(anyhow!("missing enabled"))?;

        let date_added = map
            .get("date_added")
            .map(|v| v.to_owned().as_datetime().0)
            .ok_or(anyhow!("missing date_added"))?;
        let explorer_url = map.get("explorer_url").map(|v| v.to_owned().as_string());
        let rpc_http = map
            .get("rpc_http")
            .map(|v| v.to_owned().as_string())
            .ok_or(anyhow!("missing rpc_http"))?;
        let rpc_ws = map.get("rpc_ws").map(|v| v.to_owned().as_string());

        let ep = Endpoint {
            id,
            name,
            is_dev,
            enabled,
            date_added,
            explorer_url,
            rpc_http,
            rpc_ws,
        };

        // println!("ep: {:#?}", ep);
        Ok(ep)
    }
}

pub trait Creatable: Into<Value> {}

impl Creatable for Endpoint {}
