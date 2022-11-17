package rabbit

import (
	"context"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "user_create"

type sender struct {
	cc *amqp.Connection
	ch *amqp.Channel
}

func Produce() {
	s := createChannel()
	defer s.cc.Close()
	defer s.ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 7; i++ {
		message := 4
		if i%2 == 0 {
			message = 1
		}
		// message := 1 + rand.Intn(3)
		err := s.ch.PublishWithContext(ctx,
			"",
			queueName,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(strconv.Itoa(message)),
			})

		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %d", message)
	}
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
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return &s
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
