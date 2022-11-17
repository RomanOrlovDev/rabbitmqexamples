package simplesendreceive

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "hello"

type sender struct {
	cc *amqp.Connection
	ch *amqp.Channel
}

func Send() {
	s := createChannel()
	defer s.cc.Close()
	defer s.ch.Close()

	body := "Hello World!"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

func createChannel() *sender {
	s := sender{}
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	s.cc = conn

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	s.ch = ch

	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return &s
}
