package repositories

import (
	"context"

	"youtube-backend/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobRepository interface {
	Create(ctx context.Context, job *entities.Job) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Job, error)
	Update(ctx context.Context, job *entities.Job) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByVideoID(ctx context.Context, videoID primitive.ObjectID) ([]*entities.Job, error)
	GetByStatus(ctx context.Context, status entities.JobStatus) ([]*entities.Job, error)
	GetPendingJobs(ctx context.Context, limit int) ([]*entities.Job, error)
	GetActiveJobs(ctx context.Context) ([]*entities.Job, error)
	UpdateProgress(ctx context.Context, id primitive.ObjectID, progress int) error
}
