package domain

import (
	"context"
	"tts-poc-service/database/models"
)

type CommandRepository interface {
	InsertStoryBook(ctx context.Context, in *models.StoryBook) error
}

type QueryRepository interface {
	ListStoryBooks(ctx context.Context) ([]*models.StoryBook, error)
	GetStoryBookById(ctx context.Context, id string) (*models.StoryBook, error)
}
