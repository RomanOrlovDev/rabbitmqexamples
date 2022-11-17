package main

import (
	"fmt"
	"log"
	"time"

	consumeronebyone "rabbitwithreconnections/consumingonebyone"
	"rabbitwithreconnections/rabbit"
	"rabbitwithreconnections/simplesendreceive"
	"rabbitwithreconnections/tmp"
)

func main5() {
	go func() {
		tmp.Receive()
	}()

	var forever chan struct{}
	tmp.Send()
	fmt.Println("waiting forever...")
	forever <- struct{}{}
}

func main4() {
	go func() {
		simplesendreceive.Receive()
	}()

	var forever chan struct{}
	simplesendreceive.Send()
	simplesendreceive.Send()
	fmt.Println("waiting forever...")
	forever <- struct{}{}
}

// running consumingonebyone
// here is important thing QoS - this is how many items could be loaded into channel
// - autoAck is in false because we need to mark finishing of work manually
func main2() {
	go func() {
		consumeronebyone.Receive()
	}()

	var forever chan struct{}
	consumeronebyone.Send()
	fmt.Println("waiting forever...")
	forever <- struct{}{}
}

func main() {
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
	fmt.Println("before")
	select {}
	fmt.Println("after")
}
