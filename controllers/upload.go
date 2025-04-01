package controllers

import (
	"abr_backend/config"
	"abr_backend/service/upload"
	"abr_backend/utils"
	"fmt"
	"log"
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
		log.Fatal("incomplete or wrong information", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "incomplete or wrong information",
		})
		return
	}

	// check if the file name is repeated

	clientID := c.GetString("client_id")

	// GetUploadstrategy
	uploadStrategy, err := upload.GetUploadStrategy(req.Strategy)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "incorrect upload strategy",
		})
	}

	uploadUrl, uploadId, err := uploadStrategy.InitializeUpload(req.Bucket, req.Name, req.FileType, clientID)

	videoID := utils.GenerateVideoID()
	db := config.GetDB()
	key := fmt.Sprintf("%s/%s", clientID, req.Name)

	if err != nil {
		log.Fatal("could not initialize upload: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not initialize upload",
		})
		return
	}

	err = utils.AddVideo(db, videoID, clientID, "initialized", key, req.Bucket, req.Strategy)

	if err != nil {
		log.Println("Error creating a video record ", err)
	}

	if req.Strategy == "single" {
		c.PureJSON(http.StatusOK, gin.H{
			"upload_url": uploadUrl,
			"strategy":   req.Strategy,
		})
	} else {
		c.PureJSON(http.StatusOK, gin.H{
			"upload_url": uploadUrl,
			"strategy":   req.Strategy,
			"partNumber": 1,
			"uploadId":   uploadId,
		})
	}

}

func GetPresignUrlHandler(c *gin.Context) {
	var req struct {
		Name       string `json:"name"`
		Bucket     string `json:"bucket"`
		FileType   string `json:"file_type"`
		PartNumber int32  `json:"part_number"`
		UploadId   string `json:"upload_id"`
	}

	clientID := c.GetString("client_id")

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		log.Fatal("incomplete or wrong information", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "incomplete or wrong information",
		})
		return
	}

	uploadUrl, err := upload.GetPresignedUrl(req.Bucket, req.Name, req.UploadId, req.PartNumber, clientID)

	if err != nil {
		log.Fatal("could not get pre-signed url: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not initialize upload",
		})
		return
	}

	c.PureJSON(http.StatusOK, gin.H{
		"upload_url": uploadUrl,
		"strategy":   "multipart",
		"partNumber": req.PartNumber,
	})

}

func CompleteUploadHandler(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Bucket   string `json:"bucket"`
		UploadId string `json:"upload_id"`
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		log.Fatal("incomplete or wrong information", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "incomplete or wrong information",
		})
		return
	}

	clientID := c.GetString("client_id")

	err := upload.CompleteUpload(req.Bucket, req.Name, req.UploadId, clientID)
	if err != nil {
		fmt.Printf("could not complete upload: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not complete upload: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "successfully completed upload",
		"uploadID": req.UploadId,
	})

}
