package pubsub

import (
	"context"
	"encoding/json"
	"thanhldt060802/internal/rabbitmqclient"
	"thanhldt060802/model"
	"time"

	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

var RabbitMQPubInstance1 IRabbitMQPub[string]
var RabbitMQPubInstance2 IRabbitMQPub[*model.DataStruct]

type IRabbitMQPub[T any] interface {
	Publish(ctx context.Context, exchange string, routingKey string, data T) error
}

type RabbitMQPub[T any] struct {
	channel *rabbitmqclient.RabbitMQChannel
}

func NewRabbitMQPub[T any]() (IRabbitMQPub[T], error) {
	if channel, err := rabbitmqclient.RabbitMQClientConnInstance.NewChannel(); err != nil {
		log.Errorf("Create new channel failed: %v", err.Error())
		return nil, err
	} else {
		return &RabbitMQPub[T]{
			channel: channel,
		}, nil
	}
}

func (rabbitMqPub *RabbitMQPub[T]) Publish(ctx context.Context, exchange string, routingKey string, data T) error {
	body, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Marshal data failed: %v", err.Error())
		return err
	}
	if err := rabbitMqPub.channel.GetChannel().PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
			Timestamp:    time.Now(),
		},
	); err != nil {
		log.Errorf("Publish %v to %v of %v failed: %v", data, routingKey, exchange, err.Error())
		return err
	}

	log.Infof("Publish %v to %v of %v successful", data, routingKey, exchange)
	return nil
}
