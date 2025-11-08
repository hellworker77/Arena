use std::collections::HashMap;
use async_trait::async_trait;
use anyhow::Result;
use tokio::sync::mpsc::Receiver;
use crate::models::delivery_mode::DeliveryMode;
use crate::models::broker_message::Message;
use crate::models::exchange_type::ExchangeType;

/// Ru
/// Репозиторий для взаимодействия с брокером сообщений (RabbitMQ-like).
/// Определяет методы для создания обменников и очередей,
/// публикации сообщений, подписки, подтверждения и маршрутизации.
/// Eng
/// Repository for interacting with a message broker (RabbitMQ-like).
/// Defines methods for creating exchanges and queues,
/// publishing messages, subscribing, acknowledging, and routing.
#[async_trait]
pub trait BrokerRepository: Send + Sync {
    /// Declares a new exchange (topic) in the message broker.
    /// # Arguments
    /// * `topic` - The name of the topic to create.
    /// # Returns
    /// * `Result<()>` - An empty result indicating success or failure.
    async fn declare_exchange(&self, topic: &str, kind: ExchangeType) -> Result<()>;

    /// Deletes an existing exchange (topic) from the message broker.
    /// # Arguments
    /// * `topic` - The name of the topic to delete.
    /// # Returns
    /// * `Result<()>` - An empty result indicating success or failure.
    async fn delete_exchange(&self, topic: &str) -> Result<()>;

    /// Declares a new queue in the message broker.
    /// # Arguments
    /// * `queue` - The name of the queue to create.
    /// # Returns
    /// * `Result<()>` - An empty result indicating success or failure.
    async fn declare_queue(&self, queue: &str) -> Result<()>;

    /// Deletes an existing queue from the message broker.
    /// # Arguments
    /// * `queue` - The name of the queue to delete.
    /// # Returns
    /// * `Result<()>` - An empty result indicating success or failure.
    async fn delete_queue(&self, queue: &str) -> Result<()>;

    /// Binds a queue to a topic with a specific routing key.
    /// # Arguments
    /// * `queue` - The name of the queue to bind.
    /// * `topic` - The name of the topic to bind to.
    /// * `routing_key` - The routing key for the binding.
    /// # Returns
    /// * `Result<()>` - An empty result indicating success or failure.
    async fn bind_queue(&self, queue: &str, topic: &str, routing_key: &str) -> Result<()>;

    /// Publishes a message to a specific exchange with a routing key.
    /// # Arguments
    /// * `exchange` - The name of the exchange to publish to.
    /// * `routing_key` - The routing key for the message.
    /// * `payload` - The message payload as a byte slice.
    /// * `headers` - A map of headers to include with the message.
    /// * `mode` - The delivery mode for the message (e.g., persistent or transient).
    /// # Returns
    /// * `Result<()>` - An empty result indicating success or failure.
    async fn publish(&self, exchange: &str, routing_key: &str, payload: &[u8], headers: HashMap<String, String>, mode: DeliveryMode) -> Result<()>;

    /// Consumes messages from a specific queue.
    /// # Arguments
    /// * `queue` - The name of the queue to consume from.
    /// # Returns
    /// * `Result<Receiver<Message>>` - A receiver channel for incoming messages.
    async fn consume(&self, queue: &str) -> Result<Receiver<Message>>;

    /// Acknowledges a message as successfully processed.
    /// # Arguments
    /// * `queue` - The name of the queue the message was consumed from.
    /// * `message_id` - The identifier of the message to acknowledge.
    /// # Returns
    /// * `Result<()>` - An empty result indicating success or failure.
    async fn ack(&self, queue: &str, message_id: &str) -> Result<()>;

    /// Negatively acknowledges a message, optionally requeuing it.
    /// # Arguments
    /// * `queue` - The name of the queue the message was consumed from.
    /// * `message_id` - The identifier of the message to negatively acknowledge.
    /// * `requeue` - A boolean indicating whether to requeue the message.
    /// # Returns
    /// * `Result<()>` - An empty result indicating success or failure.
    async fn nack(&self, queue: &str, message_id: &str, requeue: bool) -> Result<()>;

    /// Purges all messages from a specific queue.
    /// # Arguments
    /// * `queue` - The name of the queue to purge.
    /// # Returns
    /// * `Result<()>` - An empty result indicating success or failure.
    async fn purge_queue(&self, queue: &str) -> Result<()>;

    /// Checks if a specific queue exists in the message broker.
    /// # Arguments
    /// * `queue` - The name of the queue to check.
    /// # Returns
    /// * `Result<bool>` - A boolean indicating whether the queue exists.
    async fn queue_exists(&self, queue: &str) -> Result<bool>;
}