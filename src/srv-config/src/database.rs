use std::collections::HashMap;
use std::path::PathBuf;
use std::sync::Arc;

use rustbreak::deser::Ron;
use rustbreak::PathDatabase;

use crate::models;

pub type BareDb = PathDatabase<models::ServerData, Ron>;
pub type Db = Arc<BareDb>;

pub fn load_or_init(path: PathBuf) -> Db {
    log::info!("Setup database ...");
    let db: BareDb = PathDatabase::load_from_path_or(
        path,
        models::ServerData {
            alarms: HashMap::new(),
        },
    )
    .unwrap();
    log::info!("Setup database successful");
    return Arc::from(db);
}
