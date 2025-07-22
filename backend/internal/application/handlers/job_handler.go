package handlers

import (
	"context"
	"net/http"
	"time"

	"youtube-backend/internal/domain/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type JobHandler struct {
	processingService *services.ProcessingService
	logger            *zap.Logger
}

type JobResponse struct {
	ID           string         `json:"id"`
	VideoID      string         `json:"video_id"`
	Type         string         `json:"type"`
	Status       string         `json:"status"`
	Progress     int            `json:"progress"`
	ErrorMessage string         `json:"error_message,omitempty"`
	WorkerID     string         `json:"worker_id,omitempty"`
	Payload      map[string]any `json:"payload"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	StartedAt    *time.Time     `json:"started_at,omitempty"`
	CompletedAt  *time.Time     `json:"completed_at,omitempty"`
}

func NewJobHandler(processingService *services.ProcessingService, logger *zap.Logger) *JobHandler {
	return &JobHandler{
		processingService: processingService,
		logger:            logger,
	}
}

// GetJob returns the status and progress of a specific job
func (h *JobHandler) GetJob(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := h.processingService.GetJob(ctx, objectID)
	if err != nil {
		h.logger.Error("Failed to get job", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	jobResponse := JobResponse{
		ID:           job.ID.Hex(),
		VideoID:      job.VideoID.Hex(),
		Type:         string(job.Type),
		Status:       string(job.Status),
		Progress:     job.Progress,
		ErrorMessage: job.ErrorMessage,
		WorkerID:     job.WorkerID,
		Payload:      job.Payload,
		CreatedAt:    job.CreatedAt,
		UpdatedAt:    job.UpdatedAt,
		StartedAt:    job.StartedAt,
		CompletedAt:  job.CompletedAt,
	}

	c.JSON(http.StatusOK, jobResponse)
}

// GetJobsByVideoID returns all jobs for a specific video
func (h *JobHandler) GetJobsByVideoID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	videoID := c.Param("videoId")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	jobs, err := h.processingService.GetJobsByVideoID(ctx, objectID)
	if err != nil {
		h.logger.Error("Failed to get jobs for video", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get jobs"})
		return
	}

	jobResponses := make([]JobResponse, len(jobs))
	for i, job := range jobs {
		jobResponses[i] = JobResponse{
			ID:           job.ID.Hex(),
			VideoID:      job.VideoID.Hex(),
			Type:         string(job.Type),
			Status:       string(job.Status),
			Progress:     job.Progress,
			ErrorMessage: job.ErrorMessage,
			WorkerID:     job.WorkerID,
			Payload:      job.Payload,
			CreatedAt:    job.CreatedAt,
			UpdatedAt:    job.UpdatedAt,
			StartedAt:    job.StartedAt,
			CompletedAt:  job.CompletedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":  jobResponses,
		"count": len(jobResponses),
	})
}

// GetActiveJobs returns all currently processing jobs
func (h *JobHandler) GetActiveJobs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	jobs, err := h.processingService.GetActiveJobs(ctx)
	if err != nil {
		h.logger.Error("Failed to get active jobs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get active jobs"})
		return
	}

	jobResponses := make([]JobResponse, len(jobs))
	for i, job := range jobs {
		jobResponses[i] = JobResponse{
			ID:           job.ID.Hex(),
			VideoID:      job.VideoID.Hex(),
			Type:         string(job.Type),
			Status:       string(job.Status),
			Progress:     job.Progress,
			ErrorMessage: job.ErrorMessage,
			WorkerID:     job.WorkerID,
			Payload:      job.Payload,
			CreatedAt:    job.CreatedAt,
			UpdatedAt:    job.UpdatedAt,
			StartedAt:    job.StartedAt,
			CompletedAt:  job.CompletedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":  jobResponses,
		"count": len(jobResponses),
	})
}
