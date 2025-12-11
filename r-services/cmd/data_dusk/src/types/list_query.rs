use serde::Deserialize;

#[derive(Deserialize, Clone)]
pub struct ListQuery {
    pub(crate) prefix: Option<String>,
}