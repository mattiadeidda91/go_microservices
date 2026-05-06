package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"test-microservices-rabbit/consumer/handlers"
	"test-microservices-rabbit/events"
	"test-microservices-rabbit/services/rabbit"

	"github.com/joho/godotenv"
)

const queueName = "ServiceQueue"

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found, using system env")
	}
}

func main() {

	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatal("RABBITMQ_URL not set")
	}

	rabbit, err := rabbit.NewRabbitMQ(url)
	if err != nil {
		panic(err)
	}
	defer rabbit.Close()

	err = rabbit.DeclareExchange("events")
	if err != nil {
		panic(err)
	}

	err = rabbit.DeclareQueue(queueName)
	if err != nil {
		panic(err)
	}

	err = rabbit.BindQueue(queueName, "message.*", "events")
	if err != nil {
		panic(err)
	}

	// dispatcher
	dispatcher := events.NewDispatcher()

	// handler
	messageHandler := handlers.NewMessageEventHandler()

	// registrazione eventi
	dispatcher.Register("message.sent", messageHandler.HandleMessageSent)

	//subscribe to get messages from the queue
	messages, err := rabbit.Consume(queueName)
	if err != nil {
		panic(err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg := <-messages:
			var event events.Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Println("Invalid event:", err)
				continue
			}

			log.Printf("Event received: %s\n", event.Type)

			// dispatch automatico
			dispatcher.Dispatch(event)

			/* old managed without event dispatcher
			switch event.Type {
			case "message.sent":
				var payload events.MessageSent
				if err := json.Unmarshal(event.Data, &payload); err != nil {
					log.Println("Invalid payload:", err)
					continue
				}

				log.Printf("Message received: %s\n", payload.Message)

			default:
				log.Println("Unknown event:", event.Type)
			}*/

		case <-signalChannel:
			log.Printf("Interrupt detected!")
			os.Exit(0)
		}
	}
}
