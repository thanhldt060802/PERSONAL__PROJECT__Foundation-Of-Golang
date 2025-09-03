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
	PublishWithRetry(ctx context.Context, exchange string, routingKey string, data T)
}

type RabbitMQPub[T any] struct {
	channel *amqp091.Channel
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

func (rabbitMqPub *RabbitMQPub[T]) PublishWithRetry(ctx context.Context, exchange string, routingKey string, data T) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Marshal data failed: %v", err.Error())
		return
	}

	closeChan := rabbitMqPub.channel.NotifyClose(make(chan *amqp091.Error))

	for {
		err := rabbitMqPub.startPublish(ctx, exchange, routingKey, body)
		if err == nil {
			log.Infof("Publish %v to %v of %v successful", data, routingKey, exchange)
			return
		}

		log.Errorf("Publish %v to %v of %v failed: %v", data, routingKey, exchange, err.Error())

		select {
		case <-ctx.Done():
			log.Infof("Context canceled, stop publishing to %v of %v", routingKey, exchange)
			return
		case rabbitErr := <-closeChan:
			if rabbitErr != nil {
				log.Errorf("Channel of publisher for %v of %v closed: %v. Retry in 5s...", routingKey, exchange, rabbitErr.Error())
			} else {
				log.Warnf("Channel of publisher for %v of %v closed cleanly. Retry in 5s...", routingKey, exchange)
			}
			time.Sleep(5 * time.Second)

			for {
				newCh, chErr := rabbitmqclient.RabbitMQClientConnInstance.NewChannel()
				if chErr != nil {
					log.Errorf("Create new channel failed: %v. Retry in 5s...", chErr.Error())
					time.Sleep(5 * time.Second)
					continue
				}
				rabbitMqPub.channel = newCh

				closeChan = rabbitMqPub.channel.NotifyClose(make(chan *amqp091.Error))
				break
			}
		default:
			log.Warnf("Publish %v to %v of %v failed but channel still open. Retry in 5s...", data, routingKey, exchange)
			time.Sleep(5 * time.Second)
		}
	}
}

func (rabbitMqPub *RabbitMQPub[T]) startPublish(ctx context.Context, exchange string, routingKey string, body []byte) error {
	return rabbitMqPub.channel.PublishWithContext(
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
	)
}
