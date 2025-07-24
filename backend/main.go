package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"youtube-backend/internal/infrastructure/database"
	"youtube-backend/internal/infrastructure/queue"
	"youtube-backend/internal/infrastructure/storage"
	httphandlers "youtube-backend/internal/interfaces/http"
	"youtube-backend/pkg/config"
	"youtube-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	log := logger.New()
	defer log.Sync()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Disconnect()

	// Initialize Redis
	redisClient := queue.NewRedisClient(cfg.RedisURI)
	defer redisClient.Close()

	// Initialize MinIO
	minioClient, err := storage.NewMinIOClient(cfg.MinIO, log)
	if err != nil {
		log.Fatal("Failed to connect to MinIO", zap.Error(err))
	}

	// Initialize router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup routes
	httphandlers.SetupRoutes(router, db, redisClient, minioClient, log)

	// Create HTTP server
	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Info("Starting server", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}
