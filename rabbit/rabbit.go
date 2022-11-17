package rabbit

import (
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
}

type Rabbit struct {
	config     Config
	connection *amqp.Connection
}

// NewRabbit returns a RabbitMQ instance.
func NewRabbit(config Config) *Rabbit {
	return &Rabbit{
		config: config,
	}
}

// Connect connects to RabbitMQ server.
func (r *Rabbit) Connect() error {
	if r.connection == nil || r.connection.IsClosed() {
		con, err := amqp.Dial(fmt.Sprintf(
			"amqp://%s:%s@%s:%s",
			r.config.Username,
			r.config.Password,
			r.config.Host,
			r.config.Port,
		))
		if err != nil {
			return err
		}
		r.connection = con
	}

	return nil
}

// Connection returns exiting `*amqp.Connection` instance.
func (r *Rabbit) Connection() (*amqp.Connection, error) {
	if r.connection == nil || r.connection.IsClosed() {
		return nil, errors.New("connection is not open")
	}

	return r.connection, nil
}

// Channel returns a new `*amqp.Channel` instance.
func (r *Rabbit) Channel() (*amqp.Channel, error) {
	chn, err := r.connection.Channel()
	if err != nil {
		return nil, err
	}

	return chn, nil
}
