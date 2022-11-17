package consumeronebyone

import (
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type receiver struct {
	cc *amqp.Connection
	ch *amqp.Channel
}

func Receive() {
	r := createChannelForReceiver()
	defer r.cc.Close()
	defer r.ch.Close()

	msgs, err := r.ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		true,      // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	for i := 1; i < 3; i++ {
		i := i
		go func() {
			for d := range msgs {
				v, err := strconv.Atoi(string(d.Body))
				if err != nil {
					log.Fatalln("routine:", i, "err occurred while reading value from d.Body", v)
				}

				log.Printf("Received a message: %s by routine: %d, number of seconds: %d \n", d.Body, i, v)
				t := time.Duration(v)
				time.Sleep(t * time.Second)
				log.Printf("message (%d) is processed by routine: %d \n", v, i)
				d.Ack(false)
			}
		}()
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func createChannelForReceiver() *receiver {
	r := receiver{}
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	r.cc = conn

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	r.ch = ch

	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		2,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	return &r
}
