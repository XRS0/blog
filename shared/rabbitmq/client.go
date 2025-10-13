package rabbitmq
package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client wraps RabbitMQ connection
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *slog.Logger
}

// Event types
const (
	EventArticleViewed = "article.viewed"
	EventArticleLiked  = "article.liked"
	EventArticleUnliked = "article.unliked"
	EventArticleCreated = "article.created"
)

// Event represents a message in the queue
type Event struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// NewClient creates a new RabbitMQ client
func NewClient(url string, logger *slog.Logger) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Client{
		conn:    conn,
		channel: channel,
		logger:  logger,
	}, nil
}

// DeclareExchange declares a topic exchange
func (c *Client) DeclareExchange(name string) error {
	return c.channel.ExchangeDeclare(
		name,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// DeclareQueue declares a queue
func (c *Client) DeclareQueue(name string) error {
	_, err := c.channel.QueueDeclare(
		name,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

// BindQueue binds a queue to an exchange with a routing key
func (c *Client) BindQueue(queueName, exchangeName, routingKey string) error {
	return c.channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false, // no-wait
		nil,   // arguments
	)
}

// Publish publishes an event to an exchange
func (c *Client) Publish(ctx context.Context, exchange, routingKey string, event Event) error {
	event.Timestamp = time.Now()
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = c.channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   event.Timestamp,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	c.logger.Debug("event published", "type", event.Type, "exchange", exchange, "routing_key", routingKey)
	return nil
}

// Consume starts consuming messages from a queue
func (c *Client) Consume(queueName string, handler func(Event) error) error {
	msgs, err := c.channel.Consume(
		queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for msg := range msgs {
			var event Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				c.logger.Error("failed to unmarshal event", "error", err)
				msg.Nack(false, false) // Don't requeue
				continue
			}

			c.logger.Debug("event received", "type", event.Type)
			
			if err := handler(event); err != nil {
				c.logger.Error("failed to handle event", "type", event.Type, "error", err)
				msg.Nack(false, true) // Requeue
				continue
			}

			msg.Ack(false)
		}
	}()

	return nil
}

// Close closes the connection
func (c *Client) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
