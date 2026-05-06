package main

import (
	"log"
	"os"
	"test-microservices-rabbit/producer/handlers"
	"test-microservices-rabbit/producer/routes"
	"test-microservices-rabbit/services/rabbit"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	//new connection
	rabbitService, err := rabbit.NewRabbitMQ(url)
	if err != nil {
		panic(err)
	}

	defer rabbitService.Close()

	err = rabbitService.DeclareExchange("events")
	if err != nil {
		panic(err)
	}

	// handler (controller)
	messageHandler := handlers.NewMessageHandler(rabbitService)

	router := gin.Default()

	// routes
	routes.RegisterRoutes(router, messageHandler)

	//Add route to send message to the queue
	/*router.GET("/send", func(c *gin.Context) {
		msg := c.Query("msg")
		if msg == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
			return
		}

		//publish a message to the queue
		if err := rabbitService.Publish(queueName, []byte(msg)); err != nil {
			log.Printf("Failed to publish message: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": msg, "status": "success"})
	})

	router.POST("/send", func(c *gin.Context) {
		var req MessageRequest

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Message == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "message required"})
			return
		}

		err := rabbitService.Publish(queueName, []byte(req.Message))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "publish failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "sent",
			"message": req.Message,
			"type":    req.Type,
		})
	})*/

	//start gin server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(router.Run(":" + port))
}
