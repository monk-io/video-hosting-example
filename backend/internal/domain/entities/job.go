package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobType string
type JobStatus string

const (
	JobTypeTranscode JobType = "transcode"
	JobTypeThumbnail JobType = "thumbnail"
)

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VideoID      primitive.ObjectID `json:"video_id" bson:"video_id"`
	Type         JobType            `json:"type" bson:"type"`
	Status       JobStatus          `json:"status" bson:"status"`
	Progress     int                `json:"progress" bson:"progress"` // 0-100
	ErrorMessage string             `json:"error_message,omitempty" bson:"error_message,omitempty"`
	WorkerID     string             `json:"worker_id,omitempty" bson:"worker_id,omitempty"`
	Payload      map[string]any     `json:"payload" bson:"payload"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	StartedAt    *time.Time         `json:"started_at,omitempty" bson:"started_at,omitempty"`
	CompletedAt  *time.Time         `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
}

// NewJob creates a new job entity
func NewJob(videoID primitive.ObjectID, jobType JobType, payload map[string]any) *Job {
	now := time.Now()
	return &Job{
		ID:        primitive.NewObjectID(),
		VideoID:   videoID,
		Type:      jobType,
		Status:    JobStatusPending,
		Progress:  0,
		Payload:   payload,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Start marks the job as started
func (j *Job) Start(workerID string) {
	now := time.Now()
	j.Status = JobStatusProcessing
	j.WorkerID = workerID
	j.StartedAt = &now
	j.UpdatedAt = now
}

// UpdateProgress updates the job progress
func (j *Job) UpdateProgress(progress int) {
	j.Progress = progress
	j.UpdatedAt = time.Now()
}

// Complete marks the job as completed
func (j *Job) Complete() {
	now := time.Now()
	j.Status = JobStatusCompleted
	j.Progress = 100
	j.CompletedAt = &now
	j.UpdatedAt = now
}

// Fail marks the job as failed
func (j *Job) Fail(errorMessage string) {
	now := time.Now()
	j.Status = JobStatusFailed
	j.ErrorMessage = errorMessage
	j.CompletedAt = &now
	j.UpdatedAt = now
}

// IsCompleted checks if job is completed
func (j *Job) IsCompleted() bool {
	return j.Status == JobStatusCompleted
}

// IsFailed checks if job is failed
func (j *Job) IsFailed() bool {
	return j.Status == JobStatusFailed
}

// IsProcessing checks if job is currently processing
func (j *Job) IsProcessing() bool {
	return j.Status == JobStatusProcessing
}
