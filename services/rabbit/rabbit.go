package rabbit

import (
	"encoding/json"
	"test-microservices-rabbit/events"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

func Publish[T any](r *RabbitMQ, exchange, routingKey, source string, payload T) error {

	data, _ := json.Marshal(payload)

	event := events.Event{
		ID:            uuid.New().String(),
		Type:          routingKey,
		Source:        source,
		Timestamp:     time.Now().Unix(),
		Data:          data,
		Version:       "1.0.0",
		CorrelationID: uuid.New().String(),
	}

	body, _ := json.Marshal(event)

	return r.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Works but without topic and send text/plain message
func (r *RabbitMQ) Publish(queue string, body []byte) error {
	return r.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) Consume(queue string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) DeclareQueue(name string) error {
	_, err := r.channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (r *RabbitMQ) DeclareExchange(name string) error {
	return r.channel.ExchangeDeclare(
		name,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) BindQueue(queue, routingKey, exchange string) error {
	return r.channel.QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil,
	)
}
