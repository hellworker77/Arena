use infrastructure::factory::service_factory::ServiceFactory;

#[derive(Clone)]
pub struct AppState {
    pub(crate) factory: ServiceFactory,
}