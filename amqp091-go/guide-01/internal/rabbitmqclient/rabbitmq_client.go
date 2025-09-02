package rabbitmqclient

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

type IRabbitMQClientConn interface {
	GetConnection() *amqp091.Connection
	GetChannel() *amqp091.Channel
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
	channel    *amqp091.Channel
}

func NewRabbitMQClient(config RabbitMQConfig) IRabbitMQClientConn {
	client := &RabbitMQClientConn{}
	client.RabbitMQConfig = config
	if err := client.Connect(); err != nil {
		log.Fatalf("Connect to rabbitmq failed: %v", err.Error())
	}

	return client
}

func (c *RabbitMQClientConn) Connect() error {
	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v/", c.Username, c.Password, c.Host, c.Port))
	if err != nil {
		return err
	}
	c.connection = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	c.channel = ch

	return nil
}

func (c *RabbitMQClientConn) GetConnection() *amqp091.Connection {
	return c.connection
}

func (c *RabbitMQClientConn) GetChannel() *amqp091.Channel {
	return c.channel
}

func (c *RabbitMQClientConn) DeclareExchange(exchange string, kind string) error {
	return c.channel.ExchangeDeclare(exchange, kind, true, false, false, false, nil)
}
