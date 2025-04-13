package controllers

import (
	"abr_backend/config"
	"abr_backend/utils"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// upload → transcode → vmaf_check → thumbnail_gen → watermark → notify → analytics

func SnsHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// First try to get records to validate the message
	bucket, key, err := utils.GetRecords(body)
	if err != nil {
		log.Printf("Error getting records: %v", err)
		// Return 200 to acknowledge receipt but indicate processing failed
		// This will trigger SNS retry
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Message received but processing failed",
			"error":   err.Error(),
		})
		return
	}

	// Create a message for the queue
	message := struct {
		Bucket string `json:"bucket"`
		Key    string `json:"key"`
	}{
		Bucket: bucket,
		Key:    key,
	}

	// Convert struct to JSON bytes
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare message"})
		return
	}

	// Publish to queue
	err = utils.PublishEvent(config.Channel, config.Queue, messageBytes)
	if err != nil {
		log.Printf("Error publishing to queue: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue transcoding job"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transcoding job queued successfully"})
}

type SNSMessage struct {
	Type         string `json:"Type"`
	Message      string `json:"Message"`
	SubscribeURL string `json:"SubscribeURL,omitempty"`
}

func SnsSubscriber(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read body"})
		return
	}

	var msg SNSMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		log.Println("Failed to parse SNS message:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Handle subscription confirmation
	if msg.Type == "SubscriptionConfirmation" {
		log.Println("Confirming subscription:", msg.SubscribeURL)
		http.Get(msg.SubscribeURL) // Confirm the subscription
	}

	// Process SNS message
	if msg.Type == "Notification" {
		log.Println("SNS Notification received:", msg.Message)
		// Extract bucket & key from SNS message if it's an S3 event
	}

	c.JSON(http.StatusOK, gin.H{"message": "Received"})
}
