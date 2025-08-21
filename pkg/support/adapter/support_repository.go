package adapter

import (
	"context"
	"database/sql"
	"github.com/aarondl/sqlboiler/v4/boil"
	"tts-poc-service/database/models"
	"tts-poc-service/lib/baselogger"
)

type supportRepository struct {
	db     *sql.DB
	logger *baselogger.Logger
}

func NewSupportRepository(db *sql.DB, logger *baselogger.Logger) *supportRepository {
	return &supportRepository{
		db:     db,
		logger: logger,
	}
}

func (s *supportRepository) InsertSupport(ctx context.Context, in *models.Support) error {
	err := in.Insert(ctx, s.db, boil.Infer())
	return err
}
