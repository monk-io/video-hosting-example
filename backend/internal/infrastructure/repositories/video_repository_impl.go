package repositories

import (
	"context"
	"fmt"

	"youtube-backend/internal/domain/entities"
	"youtube-backend/internal/domain/repositories"
	"youtube-backend/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VideoRepositoryImpl struct {
	collection *mongo.Collection
}

func NewVideoRepository(db *database.MongoDB) repositories.VideoRepository {
	return &VideoRepositoryImpl{
		collection: db.GetCollection("videos"),
	}
}

func (r *VideoRepositoryImpl) Create(ctx context.Context, video *entities.Video) error {
	_, err := r.collection.InsertOne(ctx, video)
	if err != nil {
		return fmt.Errorf("failed to create video: %w", err)
	}
	return nil
}

func (r *VideoRepositoryImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Video, error) {
	var video entities.Video
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&video)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("video not found")
		}
		return nil, fmt.Errorf("failed to get video: %w", err)
	}
	return &video, nil
}

func (r *VideoRepositoryImpl) Update(ctx context.Context, video *entities.Video) error {
	filter := bson.M{"_id": video.ID}
	update := bson.M{"$set": video}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("video not found")
	}

	return nil
}

func (r *VideoRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("video not found")
	}

	return nil
}

func (r *VideoRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.Video, error) {
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list videos: %w", err)
	}
	defer cursor.Close(ctx)

	var videos []*entities.Video
	for cursor.Next(ctx) {
		var video entities.Video
		if err := cursor.Decode(&video); err != nil {
			return nil, fmt.Errorf("failed to decode video: %w", err)
		}
		videos = append(videos, &video)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return videos, nil
}

func (r *VideoRepositoryImpl) GetByStatus(ctx context.Context, status entities.VideoStatus) ([]*entities.Video, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get videos by status: %w", err)
	}
	defer cursor.Close(ctx)

	var videos []*entities.Video
	for cursor.Next(ctx) {
		var video entities.Video
		if err := cursor.Decode(&video); err != nil {
			return nil, fmt.Errorf("failed to decode video: %w", err)
		}
		videos = append(videos, &video)
	}

	return videos, nil
}

func (r *VideoRepositoryImpl) GetByUploadedBy(ctx context.Context, uploadedBy string, limit, offset int) ([]*entities.Video, error) {
	filter := bson.M{"uploaded_by": uploadedBy}
	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get videos by uploader: %w", err)
	}
	defer cursor.Close(ctx)

	var videos []*entities.Video
	for cursor.Next(ctx) {
		var video entities.Video
		if err := cursor.Decode(&video); err != nil {
			return nil, fmt.Errorf("failed to decode video: %w", err)
		}
		videos = append(videos, &video)
	}

	return videos, nil
}

func (r *VideoRepositoryImpl) Search(ctx context.Context, query string, limit, offset int) ([]*entities.Video, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search videos: %w", err)
	}
	defer cursor.Close(ctx)

	var videos []*entities.Video
	for cursor.Next(ctx) {
		var video entities.Video
		if err := cursor.Decode(&video); err != nil {
			return nil, fmt.Errorf("failed to decode video: %w", err)
		}
		videos = append(videos, &video)
	}

	return videos, nil
}

func (r *VideoRepositoryImpl) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count videos: %w", err)
	}
	return count, nil
}
