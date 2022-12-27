use anyhow::{anyhow, Result};
use surrealdb::{
    sql::{Object, Value},
    Datastore, Response, Session,
};

pub type DB = (Datastore, Session);

// TODO: Remove anyhow
pub fn into_iter_objects(ress: Vec<Response>) -> Result<impl Iterator<Item = Result<Object>>> {
    let res = ress.into_iter().next().map(|rp| rp.result).transpose()?;

    match res {
        Some(Value::Array(arr)) => {
            let it = arr.into_iter().map(|v| match v {
                Value::Object(object) => Ok(object),
                _ => Err(anyhow!("A record was not an Object")),
            });
            Ok(it)
        }
        _ => Err(anyhow!("No records found.")),
    }
}

// TODO: This is a bit of a hack. I need to figure out how to do this better (error handling without unwraps)
pub fn iter_first_to_string(res: Result<impl Iterator<Item = Result<Object>>>) -> String {
    res.unwrap()
        .next()
        .transpose()
        .unwrap()
        .and_then(|obj| obj.get("id").map(|id| id.to_string()))
        .unwrap()
}
