package services

import (
	"context"
	"fmt"

	"youtube-backend/internal/domain/entities"
	"youtube-backend/internal/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcessingService struct {
	jobRepo   repositories.JobRepository
	videoRepo repositories.VideoRepository
}

func NewProcessingService(jobRepo repositories.JobRepository, videoRepo repositories.VideoRepository) *ProcessingService {
	return &ProcessingService{
		jobRepo:   jobRepo,
		videoRepo: videoRepo,
	}
}

// GetJob retrieves a job by ID
func (s *ProcessingService) GetJob(ctx context.Context, id primitive.ObjectID) (*entities.Job, error) {
	return s.jobRepo.GetByID(ctx, id)
}

// GetPendingJobs retrieves pending jobs for processing
func (s *ProcessingService) GetPendingJobs(ctx context.Context, limit int) ([]*entities.Job, error) {
	return s.jobRepo.GetPendingJobs(ctx, limit)
}

// StartJob marks a job as started by a worker
func (s *ProcessingService) StartJob(ctx context.Context, jobID primitive.ObjectID, workerID string) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	if job.Status != entities.JobStatusPending {
		return fmt.Errorf("job is not in pending status: %s", job.Status)
	}

	job.Start(workerID)
	return s.jobRepo.Update(ctx, job)
}

// UpdateJobProgress updates the progress of a job
func (s *ProcessingService) UpdateJobProgress(ctx context.Context, jobID primitive.ObjectID, progress int) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("invalid progress value: %d", progress)
	}

	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	job.UpdateProgress(progress)
	return s.jobRepo.Update(ctx, job)
}

// CompleteJob marks a job as completed
func (s *ProcessingService) CompleteJob(ctx context.Context, jobID primitive.ObjectID) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	job.Complete()
	if err := s.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	// Check if all jobs for this video are completed
	return s.checkVideoCompletion(ctx, job.VideoID)
}

// FailJob marks a job as failed
func (s *ProcessingService) FailJob(ctx context.Context, jobID primitive.ObjectID, errorMessage string) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	job.Fail(errorMessage)
	if err := s.jobRepo.Update(ctx, job); err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	// Mark video as failed if any job fails
	return s.markVideoAsFailed(ctx, job.VideoID, errorMessage)
}

// GetJobsByVideoID retrieves all jobs for a video
func (s *ProcessingService) GetJobsByVideoID(ctx context.Context, videoID primitive.ObjectID) ([]*entities.Job, error) {
	return s.jobRepo.GetByVideoID(ctx, videoID)
}

// GetActiveJobs retrieves all currently processing jobs
func (s *ProcessingService) GetActiveJobs(ctx context.Context) ([]*entities.Job, error) {
	return s.jobRepo.GetActiveJobs(ctx)
}

// checkVideoCompletion checks if all jobs for a video are completed and updates video status
func (s *ProcessingService) checkVideoCompletion(ctx context.Context, videoID primitive.ObjectID) error {
	jobs, err := s.jobRepo.GetByVideoID(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to get jobs for video: %w", err)
	}

	allCompleted := true
	for _, job := range jobs {
		if !job.IsCompleted() {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		// Update video status to ready
		video, err := s.videoRepo.GetByID(ctx, videoID)
		if err != nil {
			return fmt.Errorf("failed to get video: %w", err)
		}

		video.UpdateStatus(entities.VideoStatusReady)
		return s.videoRepo.Update(ctx, video)
	}

	return nil
}

// markVideoAsFailed marks a video as failed
func (s *ProcessingService) markVideoAsFailed(ctx context.Context, videoID primitive.ObjectID, errorMessage string) error {
	video, err := s.videoRepo.GetByID(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	video.UpdateStatus(entities.VideoStatusFailed)
	return s.videoRepo.Update(ctx, video)
}
