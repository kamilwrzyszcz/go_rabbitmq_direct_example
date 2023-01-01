package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	rand.Seed(time.Now().Unix())
}

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

	choice := []string{"Ping", "Pong"}

	ctxSig, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for {
		select {
		case <-time.After(2 * time.Second):
			key := choice[rand.Intn(len(choice))]

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			err = channel.PublishWithContext(ctx,
				exchangeName, // exchange
				key,          // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(key),
				})
			if err != nil {
				log.Println("Failed to publish a message")
				panic(err)
			}
			log.Printf("Message sent: %s!\n", key)
		case <-ctxSig.Done():
			log.Println("Quitting...")
			return
		}
	}
}
