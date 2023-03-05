use super::types;
use anyhow::anyhow;
use std::{collections::BTreeMap, sync::Arc};
use surreal_simple_client::SurrealClient;

pub struct EndpointsService {
    client: Arc<SurrealClient>,
}

impl EndpointsService {
    pub fn new(client: Arc<SurrealClient>) -> Self {
        Self { client }
    }

    pub async fn get(&self, id: &str) -> Result<types::Endpoint, anyhow::Error> {
        todo!();

        // let (ds, sess) = &self.db;

        // let sql = "SELECT * FROM endpoints WHERE id == $id;";

        // let vars: BTreeMap<String, Value> = [("id".into(), thing(id)?.into())].into();

        // let res = ds.execute(sql, sess, Some(vars), false).await?;
        // let first_res = res
        //     .into_iter()
        //     .next()
        //     .ok_or(anyhow!("did not get a response"))?;

        // match first_res.result?.first() {
        //     Value::Object(obj) => obj.try_into(),
        //     _ => Err(anyhow!("object not found")),
        // }
    }

    pub async fn get_all(&self) -> Result<Vec<types::Endpoint>, anyhow::Error> {
        todo!();
        // let (ds, sess) = &self.db;

        // let sql = "SELECT * FROM endpoints;";

        // let res = ds.execute(sql, sess, None, false).await?;
        // let first_res = res
        //     .into_iter()
        //     .next()
        //     .ok_or(anyhow!("did not get a response"))?;

        // match first_res.result? {
        //     Value::Array(arr) => {
        //         let it = arr
        //             .into_iter()
        //             .map(|v| match v {
        //                 Value::Object(obj) => obj.try_into(),
        //                 _ => Err(anyhow!("A record was not an Object")),
        //             })
        //             .flatten();
        //         Ok(it.collect::<Vec<types::Endpoint>>())
        //         // Ok(it)
        //     }
        //     _ => Err(anyhow!("Could not convert to object list")),
        // }
    }

    pub async fn create(&self, id: &str) -> Result<types::Endpoint, anyhow::Error> {
        todo!();

        // let (ds, sess) = &self.db;

        // let sql = "SELECT * FROM endpoints WHERE id == $id;";

        // let vars: BTreeMap<String, Value> = [("id".into(), thing(id)?.into())].into();

        // let res = ds.execute(sql, sess, Some(vars), false).await?;
        // let first_res = res
        //     .into_iter()
        //     .next()
        //     .ok_or(anyhow!("did not get a response"))?;

        // match first_res.result?.first() {
        //     Value::Object(obj) => obj.try_into(),
        //     _ => Err(anyhow!("object not found")),
        // }
    }
}
