use crate::db;
use surrealdb::{Datastore, Session};

pub struct SettingsService {
    db: DB,
}

impl SettingsService {
    pub fn new(ds: Datastore) -> Self {
        let session = Session::for_db("settings", "zeth_db");
        Self { db: (ds, session) }
    }
}
