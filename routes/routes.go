package routes

import (
	"abr_backend/controllers"
	"abr_backend/middleware"

	"github.com/gin-gonic/gin"
)

func InitGin() *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")

	r.GET("/ping", controllers.Ping)

	apiV1.POST("/initialize_upload", middleware.Authenticate(), controllers.InitializeUploadHandler)
	apiV1.POST("/generate_presign_url", middleware.Authenticate(), controllers.GetPresignUrlHandler)
	apiV1.POST("/complete_upload", middleware.Authenticate(), controllers.CompleteUploadHandler)
	apiV1.POST("/sns_handler", controllers.SnsHandler)
	apiV1.POST("/register", controllers.RegisterHandler)
	apiV1.POST("/login", controllers.LoginHandler)

	return r
}
