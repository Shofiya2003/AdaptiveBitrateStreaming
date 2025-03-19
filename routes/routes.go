package routes

import (
	"abr_backend/controllers"

	"github.com/gin-gonic/gin"
)

func InitGin() *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")

	r.GET("/ping", controllers.Ping)

	apiV1.POST("/initialize_upload", controllers.InitializeUploadHandler)
	apiV1.POST("/generate_presign_url", controllers.GetPresignUrlHandler)
	apiV1.POST("/complete_upload", controllers.CompleteUploadHandler)

	return r
}
