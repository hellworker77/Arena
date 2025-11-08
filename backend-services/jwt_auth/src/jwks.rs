use serde::Deserialize;
use reqwest::Client;
use jsonwebtoken::DecodingKey;

#[derive(Debug, Deserialize, Clone)]
pub struct JwkKey {
    pub kty: String,
    pub kid: String,
    #[serde(rename = "use")]
    pub use_: String,
    pub alg: String,
    pub n: String,
    pub e: String,
}

#[derive(Debug, Deserialize, Clone)]
pub struct Jwks {
    pub keys: Vec<JwkKey>,
}

impl Jwks {
    pub async fn fetch(authority: &str) -> anyhow::Result<Self> {
        let url = format!("{}/.well-known/jwks.json", authority);
        let jwks = Client::new().get(&url).send().await?.json::<Jwks>().await?;
        Ok(jwks)
    }

    pub fn get_decoding_key(&self, kid: &str) -> Option<Result<DecodingKey, jsonwebtoken::errors::Error>> {
        self.keys
            .iter()
            .find(|k| k.kid == kid)
            .map(|jwk| DecodingKey::from_rsa_components(&jwk.n, &jwk.e))
    }
}