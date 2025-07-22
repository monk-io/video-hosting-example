package queue

import (
	"context"
	"encoding/json"

	"youtube-backend/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobPublisher struct {
	redisClient *RedisClient
}

type JobMessage struct {
	ID      string                 `json:"id"`
	VideoID string                 `json:"video_id"`
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func NewJobPublisher(redisClient *RedisClient) *JobPublisher {
	return &JobPublisher{
		redisClient: redisClient,
	}
}

func (jp *JobPublisher) PublishJob(ctx context.Context, job *entities.Job) error {
	jobMessage := JobMessage{
		ID:      job.ID.Hex(),
		VideoID: job.VideoID.Hex(),
		Type:    string(job.Type),
		Payload: job.Payload,
	}

	jobData, err := json.Marshal(jobMessage)
	if err != nil {
		return err
	}

	return jp.redisClient.Enqueue("video_jobs", jobData)
}

func (jp *JobPublisher) PublishTranscodeJob(ctx context.Context, videoID primitive.ObjectID, jobID primitive.ObjectID, quality string) error {
	jobMessage := JobMessage{
		ID:      jobID.Hex(),
		VideoID: videoID.Hex(),
		Type:    "transcode",
		Payload: map[string]interface{}{
			"quality": quality,
		},
	}

	jobData, err := json.Marshal(jobMessage)
	if err != nil {
		return err
	}

	return jp.redisClient.Enqueue("video_jobs", jobData)
}

func (jp *JobPublisher) PublishThumbnailJob(ctx context.Context, videoID primitive.ObjectID, jobID primitive.ObjectID) error {
	jobMessage := JobMessage{
		ID:      jobID.Hex(),
		VideoID: videoID.Hex(),
		Type:    "thumbnail",
		Payload: map[string]interface{}{},
	}

	jobData, err := json.Marshal(jobMessage)
	if err != nil {
		return err
	}

	return jp.redisClient.Enqueue("video_jobs", jobData)
}
