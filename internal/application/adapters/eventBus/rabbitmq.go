package eventBus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	eventBus "platform/internal/application/ports/services"
	domainEvents "platform/internal/domain/events"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// RabbitMQEventBus handles events using RabbitMQ.
type rabbitMQEventBus struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queues   map[string]string // subscriptionId -> queue name
}

func NewRabbitMQEventBus(exchangeName string) eventBus.EventBus {
	rabbitMQUser := os.Getenv("RABBITMQ_USER")
	rabbitMQPass := os.Getenv("RABBITMQ_PASS")
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	rabbitMQPort := os.Getenv("RABBITMQ_PORT")

	rabbitMQ_URL := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPass, rabbitMQHost, rabbitMQPort)

	// Establish a connection to RabbitMQ
	conn, err := amqp.Dial(rabbitMQ_URL)
	if err != nil {
		zap.L().Fatal("Failed to connect to RabbitMQ: %v", zap.Error(err))
	}
	zap.L().Info("Successfully connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		zap.L().Error("Failed to open a channel", zap.Error(err))
	}
	zap.L().Info("Channel opened successfully")

	err = ch.ExchangeDeclarePassive(
		exchangeName, // exchange name
		"topic",      // exchange type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		zap.L().Info("Exchange does not exist; declaring a new one", zap.String("exchange", exchangeName))
		err = ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil)
		if err != nil {
			ch.Close()
			conn.Close()
			zap.L().Fatal("Failed to declare exchange", zap.Error(err))
		}
		zap.L().Info("Exchange declared successfully", zap.String("exchange", exchangeName))
	} else {
		zap.L().Info("Exchange already exists", zap.String("exchange", exchangeName))
	}

	return &rabbitMQEventBus{
		conn:     conn,
		channel:  ch,
		exchange: exchangeName,
		queues:   make(map[string]string),
	}
}

// Publish sends an event to the exchange.
func (b *rabbitMQEventBus) Publish(ctx context.Context, event domainEvents.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Message should send to all services listening to this event
	// routingKey := event.Name
	return b.channel.Publish(
		b.exchange, // Exchange name
		event.Name, // Routing key
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Subscribe listens for messages with a routing key.
func (b *rabbitMQEventBus) Subscribe(subscriber, eventName string, handler func(ctx context.Context, event domainEvents.Event) error) (string, error) {
	queueName := eventName + "." + subscriber
	q, err := b.channel.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		true,      // no-wait
		nil,       // arguments
	)
	if err != nil {
		return "", err
	}

	err = b.channel.QueueBind(q.Name, eventName, b.exchange, false, nil)
	if err != nil {
		return "", err
	}

	msgs, err := b.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return "", err
	}

	subId := uuid.New().String()
	b.queues[subId] = q.Name

	go func() {
		for d := range msgs {
			var event domainEvents.Event
			if err := json.Unmarshal(d.Body, &event); err == nil {
				handler(context.Background(), event)
			}
		}
	}()

	return subId, nil
}

// Unsubscribe deletes the subscriber's queue.
func (b *rabbitMQEventBus) Unsubscribe(subscriptionId string) error {
	queue, exists := b.queues[subscriptionId]
	if !exists {
		return errors.New("subscription not found")
	}
	delete(b.queues, subscriptionId)
	_, err := b.channel.QueueDelete(queue, false, false, false)
	return err
}

// Close gracefully shuts down the RabbitMQ connection
func (b *rabbitMQEventBus) Close() {
	if b.channel != nil {
		b.channel.Close()
	}
	if b.conn != nil {
		b.conn.Close()
	}
	zap.L().Info("RabbitMQ connection closed")
}
