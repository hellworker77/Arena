#[derive(Debug, Clone, Copy)]
pub enum ExchangeType {
    Direct,
    Topic,
    Fanout,
    Headers,
}