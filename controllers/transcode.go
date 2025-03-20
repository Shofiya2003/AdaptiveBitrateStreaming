package controllers

import (
	"abr_backend/utils"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SnsHandler(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	fmt.Println("body", string(body))

	go utils.TranscodeVideo(body)

	// eventBody := data.VideoEvent{
	// 	VideoURL: url,
	// }

	// eventBodyJson, err := json.Marshal(eventBody)
	// if err != nil {
	// 	log.Println("Failed to marshal event:", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal event"})
	// 	return
	// }
	// utils.PublishEvent(config.Channel, config.Queue, eventBodyJson)

	c.JSON(http.StatusOK, gin.H{"message": "Transcoding in progress"})
}

// type SNSMessage struct {
// 	Type         string `json:"Type"`
// 	Message      string `json:"Message"`
// 	SubscribeURL string `json:"SubscribeURL,omitempty"`
// }

// func SnsHandler(c *gin.Context) {
// 	body, err := ioutil.ReadAll(c.Request.Body)
// 	if err != nil {
// 		log.Println("Failed to read request body:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read body"})
// 		return
// 	}

// 	var msg SNSMessage
// 	if err := json.Unmarshal(body, &msg); err != nil {
// 		log.Println("Failed to parse SNS message:", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
// 		return
// 	}

// 	// Handle subscription confirmation
// 	if msg.Type == "SubscriptionConfirmation" {
// 		log.Println("Confirming subscription:", msg.SubscribeURL)
// 		http.Get(msg.SubscribeURL) // Confirm the subscription
// 	}

// 	// Process SNS message
// 	if msg.Type == "Notification" {
// 		log.Println("SNS Notification received:", msg.Message)
// 		// Extract bucket & key from SNS message if it's an S3 event
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Received"})
// }
