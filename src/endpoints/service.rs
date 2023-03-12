use super::types;
use anyhow::anyhow;
use std::collections::{BTreeMap, HashMap};

use crate::surrealclient::SurrealHttpClient;

pub struct EndpointsService {
    client: SurrealHttpClient,
}

impl EndpointsService {
    pub fn new(client: SurrealHttpClient) -> Self {
        Self { client }
    }

    pub async fn get(&self, id: &str) -> Result<types::Endpoint, anyhow::Error> {
        // TODO: fix obvious SQL injection here
        let query = &format!("SELECT * FROM endpoint WHERE id == {};", id);
        let endpoint = self
            .client
            .post::<Vec<types::Endpoint>>(query)
            .await
            .map_err(|e| anyhow!("Could not get endpoint: {}", e))?;
        let endpoint = endpoint.first().unwrap();
        Ok(endpoint.to_owned())
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
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn get_endpoint() {
        let (surreal_rpc, admin_user, admin_pass, namespace, database) = (
            "http://localhost:8000/sql",
            "admin",
            "admin",
            "test",
            "test",
        );
        let surreal_http_client =
            SurrealHttpClient::new(surreal_rpc, admin_user, admin_pass, namespace, database);
        surreal_http_client.setup_db().await.unwrap();

        let query = &format!("SELECT * FROM endpoint WHERE id == endpoint:jplluhvaqi5oyajvnhh6;");
        let endpoint = surreal_http_client
            .post::<Vec<types::Endpoint>>(query)
            .await
            .map_err(|e| anyhow!("Could not get endpoint: {}", e));
        println!("{:#?}", endpoint);

        assert_eq!(1, 2);
    }
}
