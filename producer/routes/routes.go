package routes

import (
	"test-microservices-rabbit/producer/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, messageHandler *handlers.MessageHandler) {

	router.GET("/send", messageHandler.SendMessageGET)
	router.POST("/send", messageHandler.SendMessage)
}
