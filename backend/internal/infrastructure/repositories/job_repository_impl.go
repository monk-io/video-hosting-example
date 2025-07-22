package repositories

import (
	"context"
	"fmt"
	"time"

	"youtube-backend/internal/domain/entities"
	"youtube-backend/internal/domain/repositories"
	"youtube-backend/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobRepositoryImpl struct {
	collection *mongo.Collection
}

func NewJobRepository(db *database.MongoDB) repositories.JobRepository {
	return &JobRepositoryImpl{
		collection: db.GetCollection("jobs"),
	}
}

func (r *JobRepositoryImpl) Create(ctx context.Context, job *entities.Job) error {
	_, err := r.collection.InsertOne(ctx, job)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}
	return nil
}

func (r *JobRepositoryImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Job, error) {
	var job entities.Job
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	return &job, nil
}

func (r *JobRepositoryImpl) Update(ctx context.Context, job *entities.Job) error {
	filter := bson.M{"_id": job.ID}
	update := bson.M{"$set": job}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}

func (r *JobRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}

func (r *JobRepositoryImpl) GetByVideoID(ctx context.Context, videoID primitive.ObjectID) ([]*entities.Job, error) {
	filter := bson.M{"video_id": videoID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by video ID: %w", err)
	}
	defer cursor.Close(ctx)

	var jobs []*entities.Job
	for cursor.Next(ctx) {
		var job entities.Job
		if err := cursor.Decode(&job); err != nil {
			return nil, fmt.Errorf("failed to decode job: %w", err)
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (r *JobRepositoryImpl) GetByStatus(ctx context.Context, status entities.JobStatus) ([]*entities.Job, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by status: %w", err)
	}
	defer cursor.Close(ctx)

	var jobs []*entities.Job
	for cursor.Next(ctx) {
		var job entities.Job
		if err := cursor.Decode(&job); err != nil {
			return nil, fmt.Errorf("failed to decode job: %w", err)
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (r *JobRepositoryImpl) GetPendingJobs(ctx context.Context, limit int) ([]*entities.Job, error) {
	filter := bson.M{"status": entities.JobStatusPending}
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: 1}}) // FIFO order

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending jobs: %w", err)
	}
	defer cursor.Close(ctx)

	var jobs []*entities.Job
	for cursor.Next(ctx) {
		var job entities.Job
		if err := cursor.Decode(&job); err != nil {
			return nil, fmt.Errorf("failed to decode job: %w", err)
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (r *JobRepositoryImpl) GetActiveJobs(ctx context.Context) ([]*entities.Job, error) {
	filter := bson.M{"status": entities.JobStatusProcessing}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get active jobs: %w", err)
	}
	defer cursor.Close(ctx)

	var jobs []*entities.Job
	for cursor.Next(ctx) {
		var job entities.Job
		if err := cursor.Decode(&job); err != nil {
			return nil, fmt.Errorf("failed to decode job: %w", err)
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (r *JobRepositoryImpl) UpdateProgress(ctx context.Context, id primitive.ObjectID, progress int) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"progress":   progress,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update job progress: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("job not found")
	}

	return nil
}
