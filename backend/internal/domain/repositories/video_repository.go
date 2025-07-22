package repositories

import (
	"context"

	"youtube-backend/internal/domain/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoRepository interface {
	Create(ctx context.Context, video *entities.Video) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entities.Video, error)
	Update(ctx context.Context, video *entities.Video) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Video, error)
	GetByStatus(ctx context.Context, status entities.VideoStatus) ([]*entities.Video, error)
	GetByUploadedBy(ctx context.Context, uploadedBy string, limit, offset int) ([]*entities.Video, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*entities.Video, error)
	Count(ctx context.Context) (int64, error)
}
