package controllers

import (
	"abr_backend/service/upload"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitializeUploadHandler(c *gin.Context) {

	var req struct {
		Strategy string `json:"strategy"`
		Name     string `json:"name"`
		Bucket   string `json:"bucket"`
		FileType string `json:"file_type"`
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "incomplete or wrong information",
		})
		return
	}

	// check if bucket exists
	// check if the file name is repeated

	// GetUploadstrategy
	uploadStrategy, err := upload.GetUploadStrategy(req.Strategy)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "incorrect upload strategy",
		})
	}

	if req.Strategy == "single" {
		uploadUrl, err := uploadStrategy.InitializeUpload(req.Bucket, req.Name, req.FileType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "could not initialize upload",
			})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"upload_url": uploadUrl,
			"strategy":   req.Strategy,
		})

		return
	}

	// strategy.InitializeUpload

}
