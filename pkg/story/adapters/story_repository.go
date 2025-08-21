package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aarondl/sqlboiler/v4/boil"
	"tts-poc-service/database/models"
	"tts-poc-service/lib/baselogger"
	pkgError "tts-poc-service/pkg/common/error"
)

type storyRepository struct {
	db     *sql.DB
	logger *baselogger.Logger
}

func NewStoryRepository(db *sql.DB, logger *baselogger.Logger) *storyRepository {
	return &storyRepository{
		db:     db,
		logger: logger,
	}
}

func (p *storyRepository) InsertStoryBook(ctx context.Context, in *models.StoryBook) error {
	return in.Insert(ctx, p.db, boil.Infer())
}

func (p *storyRepository) ListStoryBooks(ctx context.Context) ([]*models.StoryBook, error) {
	return models.StoryBooks().All(ctx, p.db)
}

func (p *storyRepository) GetStoryBookById(ctx context.Context, id string) (*models.StoryBook, error) {
	isExist, err := models.StoryBookExists(ctx, p.db, id)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, fmt.Errorf(pkgError.STORY_BOOK_NOT_FOUND)
	}

	out, err := models.FindStoryBook(ctx, p.db, id)
	if err != nil {
		return nil, err
	}
	return out, nil
}
