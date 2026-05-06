package handlers

import (
	"net/http"

	"test-microservices-rabbit/events"
	"test-microservices-rabbit/services/rabbit"

	"github.com/gin-gonic/gin"
)

//const queueName = "ServiceQueue"

type MessageHandler struct {
	rabbit *rabbit.RabbitMQ
}

func NewMessageHandler(r *rabbit.RabbitMQ) *MessageHandler {
	return &MessageHandler{rabbit: r}
}

type MessageRequest struct {
	Message string `json:"message"`
}

func (h *MessageHandler) SendMessage(c *gin.Context) {

	var req MessageRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message required"})
		return
	}

	payload := events.MessageSent{
		Message: req.Message,
		Type:    "info",
	}

	err := rabbit.Publish(
		h.rabbit,
		"events",
		"message.sent",
		"producer-service",
		payload,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish"})
		return
	}

	//Old version to send a simple text message
	/*if err := h.rabbit.Publish("ServiceQueue", []byte(req.Message)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish"})
		return
	}*/

	c.JSON(http.StatusOK, gin.H{
		"status": "sent",
		"event":  "message.sent",
	})
}

func (h *MessageHandler) SendMessageGET(c *gin.Context) {

	msg := c.Query("msg")

	if msg == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
		return
	}

	payload := events.MessageSent{
		Message: msg,
		Type:    "info",
	}

	err := rabbit.Publish(
		h.rabbit,
		"events",
		"message.sent",
		"producer-service",
		payload,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	//old version to send a simple message text
	/*if err := h.rabbit.Publish(queueName, []byte(msg)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed"})
		return
	}*/

	c.JSON(http.StatusOK, gin.H{
		"message": msg,
		"status":  "success",
	})
}
