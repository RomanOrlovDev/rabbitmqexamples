package main

import (
	"fmt"
	"log"
	"time"

	consumeronebyone "rabbitwithreconnections/consumingonebyone"
	"rabbitwithreconnections/rabbit"
	"rabbitwithreconnections/simplesendreceive"
)

// simpleSendReceive example of simple sending messages from main routine (u can create as many producers as u wish)
// and consuming from a routine
// todo add description
func simpleSendReceive() {
	go func() {
		simplesendreceive.Receive()
	}()

	var forever chan struct{}
	simplesendreceive.Send()
	simplesendreceive.Send()
	fmt.Println("waiting forever...")
	forever <- struct{}{}
}

// consumeonebyone running package consumingonebyone
// Producer produces messages which may have either number one or four. It produces consequently 1 than 4.
// Consumer creates 2 goroutines. When a routine consumes a message
// it will wait as long seconds as consumed number.
// Another interesting thing - when we create a channel for consumers,
// we apply QoS property which controls prefetch property. This is how many items could be loaded by one consumer.
// Property "autoAck" is in false because we need to mark finishing of work manually
func consumeonebyone() {
	go func() {
		consumeronebyone.Receive()
	}()

	var forever chan struct{}
	consumeronebyone.ProduceOneOrFourSeconds()
	fmt.Println("waiting forever...")
	forever <- struct{}{}
}

func main() {
	// rabbitClientWithReconnections()
	consumeonebyone()
}

// this is a demonstration of custom client which is put to rabbit package.
// This is able to reconnect when rabbit server is down for some reasons.
// Also, it holds one tcp connection per initialization, so u can initialize rabbit client separately and per this client
// crate as many consumers/producers as u wish. Only one TCP connection will be in use
func rabbitClientWithReconnections() {
	go func() {
		rabbit.Produce()
	}()
	// RabbitMQ
	rbt := rabbit.NewRabbit(rabbit.Config{
		Username: "guest",
		Password: "guest",
		Host:     "localhost",
		Port:     "5672",
	})
	if err := rbt.Connect(); err != nil {
		log.Fatalln("unable to connect to rabbit", err)
	}
	//

	// Consumer
	cc := rabbit.ConsumerConfig{
		// ExchangeName:  "user",
		// ExchangeType:  "direct",
		// RoutingKey:    "create",
		QueueName:     "user_create",
		ConsumerName:  "my_app_name",
		ConsumerCount: 3,
	}
	cc.Reconnect.MaxAttempt = 10
	cc.Reconnect.Interval = 1 * time.Second
	csm := rabbit.NewConsumer(cc, rbt)
	if err := csm.Start(); err != nil {
		log.Fatalln("unable to start consumer", err)
	}
	//
	select {}
}
