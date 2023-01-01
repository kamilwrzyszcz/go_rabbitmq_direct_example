package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("Failed to connect to RabbitMQ")
		panic(err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Println("Failed to set up a channel")
		panic(err)
	}
	defer channel.Close()

	exchangeName := "ping-pong"
	err = channel.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		false,        // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Println("Failed to set up exchange")
		panic(err)
	}

	q, err := channel.QueueDeclare(
		"",    // name random
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Println("Failed to declare a queue")
		panic(err)
	}

	log.Printf("Binding queue %s to exchange %s with routing key %s",
		q.Name, exchangeName, os.Args[1])
	err = channel.QueueBind(
		q.Name,       // queue name
		os.Args[1],   // routing key
		exchangeName, // exchange
		false,
		nil)
	if err != nil {
		log.Println("Failed to bind a queue")
		panic(err)
	}

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		log.Println("Failed to register a consumer")
		panic(err)
	}

	replyMap := map[string]string{"Ping": "Pong", "Pong": "Ping"}
	ctxSig, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for {
		select {
		case msg := <-msgs:
			log.Printf("%s received!\n", msg.Body)
			log.Printf("%s!\n", replyMap[string(msg.Body)])
		case <-ctxSig.Done():
			log.Println("Quitting...")
			return
		}
	}
}
