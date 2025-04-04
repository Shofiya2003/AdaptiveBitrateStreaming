package controllers

import (
	"abr_backend/config"
	"abr_backend/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetVideosHandler(c *gin.Context) {

	clientID := c.GetString("client_id")
	videos, err := utils.GetVideos(config.GetDB(), clientID, 1)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Fetched Videos:", videos)

	c.JSON(http.StatusAccepted, gin.H{
		"videos": videos,
	})
}
