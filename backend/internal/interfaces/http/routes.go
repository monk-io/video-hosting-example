package http

import (
	"youtube-backend/internal/application/handlers"
	"youtube-backend/internal/domain/services"
	"youtube-backend/internal/infrastructure/database"
	"youtube-backend/internal/infrastructure/queue"
	"youtube-backend/internal/infrastructure/repositories"
	"youtube-backend/internal/infrastructure/storage"
	"youtube-backend/internal/interfaces/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRoutes(router *gin.Engine, db *database.MongoDB, redis *queue.RedisClient, minio *storage.MinIOClient, logger *zap.Logger) {
	// Add middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger(logger))

	// Initialize repositories
	videoRepo := repositories.NewVideoRepository(db)
	jobRepo := repositories.NewJobRepository(db)
	// userRepo := repositories.NewUserRepository(db) // TODO: Implement user handlers

	// Initialize job publisher
	jobPublisher := queue.NewJobPublisher(redis)

	// Initialize services
	videoService := services.NewVideoService(videoRepo, jobRepo, jobPublisher)
	processingService := services.NewProcessingService(jobRepo, videoRepo)

	// Initialize handlers
	videoHandler := handlers.NewVideoHandler(videoService, minio, logger)
	jobHandler := handlers.NewJobHandler(processingService, logger)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "youtube-backend",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Video routes
		videos := v1.Group("/videos")
		{
			videos.POST("/upload", videoHandler.UploadVideo)
			videos.GET("", videoHandler.GetVideos)
			videos.GET("/:id", videoHandler.GetVideo)
			videos.GET("/:id/stream", videoHandler.StreamVideo)
			videos.GET("/:id/thumbnail", videoHandler.GetThumbnail)
			videos.POST("/:id/process", videoHandler.ProcessVideo)
		}

		// Job routes
		jobs := v1.Group("/jobs")
		{
			jobs.GET("/:id", jobHandler.GetJob)
			jobs.GET("/video/:videoId", jobHandler.GetJobsByVideoID)
			jobs.GET("/active", jobHandler.GetActiveJobs)
		}

		// User routes (basic implementation for future use)
		users := v1.Group("/users")
		{
			users.POST("", func(c *gin.Context) {
				c.JSON(501, gin.H{"error": "User creation not implemented yet"})
			})
			users.GET("/:id", func(c *gin.Context) {
				c.JSON(501, gin.H{"error": "User retrieval not implemented yet"})
			})
		}
	}
}
