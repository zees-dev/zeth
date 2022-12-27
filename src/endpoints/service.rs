use anyhow::anyhow;
use std::collections::BTreeMap;
use surrealdb::{
    sql::{thing, Array, Object, Value},
    Datastore, Error, Response, Session,
};

use crate::db;

use super::types;

pub struct EndpointsService {
    db: db::DB,
}

impl EndpointsService {
    pub fn new(ds: Datastore) -> Self {
        let session = Session::for_db("endpoints", "zeth-db");
        Self { db: (ds, session) }
    }

    pub async fn create(&self, ep: &types::Endpoint) -> Result<String, anyhow::Error> {
        let (ds, sess) = &self.db;

        let sql = "CREATE endpoints CONTENT $data;";

        let vars: BTreeMap<String, Value> = [("data".into(), ep.clone().into())].into();

        let res = ds.execute(sql, sess, Some(vars), false).await?;

        let id = response_to_id(res)?;

        Ok(id)
    }

    pub async fn get(&self, id: &str) -> Result<types::Endpoint, anyhow::Error> {
        let (ds, sess) = &self.db;

        let sql = "SELECT * FROM endpoints WHERE id == $id;";

        let vars: BTreeMap<String, Value> = [("id".into(), thing(id)?.into())].into();

        let res = ds.execute(sql, sess, Some(vars), false).await?;
        let first_res = res
            .into_iter()
            .next()
            .ok_or(anyhow!("did not get a response"))?;

        match first_res.result?.first() {
            Value::Object(obj) => obj.try_into(),
            _ => Err(anyhow!("object not found")),
        }
    }

    pub async fn get_all(&self) -> Result<Vec<types::Endpoint>, anyhow::Error> {
        let (ds, sess) = &self.db;

        let sql = "SELECT * FROM endpoints;";

        let res = ds.execute(sql, sess, None, false).await?;
        let first_res = res
            .into_iter()
            .next()
            .ok_or(anyhow!("did not get a response"))?;

        match first_res.result? {
            Value::Array(arr) => {
                let it = arr
                    .into_iter()
                    .map(|v| match v {
                        Value::Object(obj) => obj.try_into(),
                        _ => Err(anyhow!("A record was not an Object")),
                    })
                    .flatten();
                Ok(it.collect::<Vec<types::Endpoint>>())
                // Ok(it)
            }
            _ => Err(anyhow!("Could not convert to object list")),
        }
    }

    pub async fn update(&self, id: &str, ep: &types::Endpoint) -> Result<String, anyhow::Error> {
        let (ds, sess) = &self.db;

        // can add properties here - on the fly if needed
        let sql = "UPDATE $th MERGE $data RETURN id;";
        let vars: BTreeMap<String, Value> = [
            ("th".into(), thing(id)?.into()),
            ("data".into(), ep.clone().into()),
        ]
        .into();

        let res = ds.execute(sql, sess, Some(vars), true).await?;

        let id = response_to_id(res)?;

        Ok(id)
    }

    pub async fn delete(&self, id: &str) -> Result<Vec<Response>, Error> {
        let (ds, sess) = &self.db;

        let sql = "DELETE $th;";
        let vars: BTreeMap<String, Value> = [("th".into(), thing(id)?.into())].into();
        let res = ds.execute(sql, sess, Some(vars), true).await?;

        Ok(res)
    }
}

fn response_to_id(res: Vec<Response>) -> Result<String, anyhow::Error> {
    let value = res
        .into_iter()
        .next()
        .map(|rp| rp.result)
        .transpose()?
        .ok_or(anyhow!("did not get a response"))?;

    let id = match value {
        Value::Array(arr) => {
            let it = arr.into_iter().map(|v| match v {
                Value::Object(object) => Ok(object),
                _ => Err(anyhow!("record was not an Object")),
            });
            Ok(it)
        }
        _ => Err(anyhow!("no records found.")),
    }?
    .next()
    .transpose()
    .map_err(|e| anyhow!(e))?
    .and_then(|obj| obj.get("id").map(|id| id.to_string()))
    .ok_or(anyhow!("did not get an id"))?;

    Ok(id)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::db;
    use crate::endpoints::types;
    use chrono::Utc;
    use std::env;

    #[tokio::test]
    async fn test_create() {
        let ds = Datastore::new("memory").await.unwrap();
        let service = EndpointsService::new(ds);

        let ep = types::Endpoint {
            id: Some("test".into()),
            name: "test".into(),
            is_dev: false,
            enabled: true,
            date_added: Utc::now(),
            explorer_url: Some("test".to_string()),
            rpc_http: "test".to_string(),
            rpc_ws: Some("test".to_string()),
        };
        let id = service.create(&ep).await.unwrap();

        let endpoints_prefix = id.split(":").next().unwrap();

        // println!("{:?}", id);
        assert_eq!(endpoints_prefix, "endpoints");
    }

    #[tokio::test]
    async fn test_get() {
        let ds = Datastore::new("memory").await.unwrap();
        let service = EndpointsService::new(ds);

        let id = service
            .create(&types::Endpoint {
                id: Some("test".into()), // this is irrelevant
                name: "test".to_string(),
                is_dev: false,
                enabled: true,
                date_added: Utc::now(),
                explorer_url: Some("test".to_string()),
                rpc_http: "test".to_string(),
                rpc_ws: Some("test".to_string()),
            })
            .await
            .unwrap();

        let endpoint = service.get(&id).await.unwrap();

        // println!("{:#?}", endpoint);
        assert_eq!(endpoint.name, "test");
    }

    #[tokio::test]
    async fn test_get_all() {
        let ds = Datastore::new("memory").await.unwrap();
        let service = EndpointsService::new(ds);

        let _ = service
            .create(&types::Endpoint {
                id: Some("test".into()), // this is irrelevant
                name: "test-1".into(),
                is_dev: false,
                enabled: true,
                date_added: Utc::now(),
                explorer_url: Some("test".into()),
                rpc_http: "test".into(),
                rpc_ws: Some("test".into()),
            })
            .await
            .unwrap();
        let _ = service
            .create(&types::Endpoint {
                id: None,
                name: "test-2".into(),
                is_dev: false,
                enabled: true,
                date_added: Utc::now(),
                explorer_url: None,
                rpc_http: "test".into(),
                rpc_ws: None,
            })
            .await
            .unwrap();

        let endpoints = service.get_all().await.unwrap();

        // println!("{:?}", endpoints);
        assert_eq!(endpoints.len(), 2);
        assert_eq!(
            endpoints
                .clone()
                .into_iter()
                .find(|ep| ep.name == "test-1")
                .is_some(),
            true
        );
        assert_eq!(
            endpoints
                .clone()
                .into_iter()
                .find(|ep| ep.name == "test-2")
                .is_some(),
            true
        );
    }

    #[tokio::test]
    async fn test_update() {
        let ds = Datastore::new("memory").await.unwrap();
        let service = EndpointsService::new(ds);

        let ep = types::Endpoint {
            id: Some("test".into()), // this is irrelevant
            name: "test".to_string(),
            is_dev: false,
            enabled: true,
            date_added: Utc::now(),
            explorer_url: Some("test".to_string()),
            rpc_http: "test".to_string(),
            rpc_ws: Some("test".to_string()),
        };

        let id = service.create(&ep).await.unwrap();

        let ep_upd = types::Endpoint {
            name: "test-updated".to_string(),
            ..ep.clone()
        };
        let upd = service.update(&id, &ep_upd).await.unwrap();

        let endpoint = service.get(&id).await.unwrap();

        // println!("{:#?}", endpoint);
        assert_eq!(endpoint.name, "test-updated");
    }

    #[tokio::test]
    async fn test_delete() {
        let ds = Datastore::new("memory").await.unwrap();
        let service = EndpointsService::new(ds);

        let ep = types::Endpoint {
            id: Some("test".into()), // this is irrelevant
            name: "test".to_string(),
            is_dev: false,
            enabled: true,
            date_added: Utc::now(),
            explorer_url: Some("test".to_string()),
            rpc_http: "test".to_string(),
            rpc_ws: Some("test".to_string()),
        };
        let id = service.create(&ep).await.unwrap();
        let _ = service.delete(&id).await.unwrap();

        let response = service.get(&id).await;
        // println!("{:?}", response);
        assert_eq!(response.is_err(), true);
        assert_eq!(response.err().unwrap().to_string(), "object not found");

        let endpoints = service.get_all().await.unwrap();
        // println!("{:?}", endpoints);
        assert_eq!(endpoints.len(), 0);
    }
}
