package routes

import (
	"abr_backend/controllers"
	"abr_backend/middleware"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func InitGin() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"}, // Allow all headers
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache the preflight response
	}))

	apiV1 := r.Group("/api/v1")

	r.GET("/ping", controllers.Ping)

	apiV1.POST("/initialize_upload", middleware.Authenticate(), controllers.InitializeUploadHandler)
	apiV1.POST("/generate_presign_url", middleware.Authenticate(), controllers.GetPresignUrlHandler)
	apiV1.POST("/complete_upload", middleware.Authenticate(), controllers.CompleteUploadHandler)
	apiV1.POST("/sns_handler", controllers.SnsHandler)
	apiV1.GET("/uploaded_videos", middleware.Authenticate(), controllers.GetVideosHandler)
	apiV1.POST("/register", controllers.RegisterHandler)
	apiV1.POST("/login", controllers.LoginHandler)

	return r
}
