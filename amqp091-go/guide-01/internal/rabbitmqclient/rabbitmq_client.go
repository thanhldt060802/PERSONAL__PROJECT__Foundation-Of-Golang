package rabbitmqclient

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

var RabbitMQClientConnInstance IRabbitMQClientConn

type IRabbitMQClientConn interface {
	GetConnection() *amqp091.Connection
	NewChannel() (*amqp091.Channel, error)
	DeclareExchange(exchange string, kind string) error
}

type RabbitMQConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type RabbitMQClientConn struct {
	RabbitMQConfig
	connection *amqp091.Connection
}

func NewRabbitMQClient(config RabbitMQConfig) IRabbitMQClientConn {
	client := &RabbitMQClientConn{}
	client.RabbitMQConfig = config
	if err := client.Connect(); err != nil {
		log.Fatalf("Connect to rabbitmq failed: %v", err.Error())
	}

	client.connectionWatcher(context.Background())

	return client
}

func (c *RabbitMQClientConn) Connect() error {
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v/", c.Username, c.Password, c.Host, c.Port))
	if err != nil {
		return err
	}
	c.connection = conn

	log.Infof("Connect to RabbitMQ successful")
	return nil
}

func (c *RabbitMQClientConn) GetConnection() *amqp091.Connection {
	return c.connection
}

func (c *RabbitMQClientConn) NewChannel() (*amqp091.Channel, error) {
	return c.connection.Channel()
}

func (c *RabbitMQClientConn) DeclareExchange(exchange string, kind string) error {
	channel, err := c.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.ExchangeDeclare(exchange, kind, true, false, false, false, nil)
}

func (c *RabbitMQClientConn) connectionWatcher(ctx context.Context) {
	go func() {
		for {
			closeChan := c.connection.NotifyClose(make(chan *amqp091.Error))

			select {
			case <-ctx.Done():
				c.connection.Close()
				log.Infof("Context canceled, stop watcher")
				return
			case rabbitErr := <-closeChan:
				if rabbitErr != nil {
					log.Errorf("RabbitMQ connection closed: %v. Retry in 5s...", rabbitErr.Error())
				} else {
					log.Warnf("RabbitMQ connection closed cleanly. Retry in 5s...")
				}
				time.Sleep(5 * time.Second)

				for {
					if err := c.Connect(); err != nil {
						log.Warnf("Reconnect to rabbitmq failed: %v. Retry in 5s...", err.Error())
						time.Sleep(5 * time.Second)
						continue
					}

					log.Infof("Reconnect to RabbitMQ successful")
					break
				}
			}
		}
	}()
}
