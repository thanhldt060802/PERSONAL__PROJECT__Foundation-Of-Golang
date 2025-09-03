package rabbitmqclient

import (
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQChannel struct {
	channel *amqp091.Channel
}

func (ch *RabbitMQChannel) Qos(prefetchCount int) error {
	return ch.channel.Qos(prefetchCount, 0, false)
}

func (ch *RabbitMQChannel) QueueDeclare(queue string) error {
	_, err := ch.channel.QueueDeclare(queue, true, false, false, false, nil)
	return err
}

func (ch *RabbitMQChannel) QueueBind(queue string, routingKey string, exchange string) error {
	return ch.channel.QueueBind(queue, routingKey, exchange, false, nil)
}

func (ch *RabbitMQChannel) GetChannel() *amqp091.Channel {
	return ch.channel
}
