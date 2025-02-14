package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/models"
	pkgError "tts-poc-service/pkg/common/error"
)

type pdfRepository struct {
	db     *sql.DB
	logger *baselogger.Logger
}

func NewPdfRepository(db *sql.DB, logger *baselogger.Logger) *pdfRepository {
	return &pdfRepository{
		db:     db,
		logger: logger,
	}
}

func (p *pdfRepository) InsertBlockKey(ctx context.Context, in *models.Block) error {
	return in.Insert(ctx, p.db, boil.Infer())
}

func (p *pdfRepository) GetBlockById(ctx context.Context, key string) (*models.Block, error) {
	isExist, err := models.BlockExists(ctx, p.db, key)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, fmt.Errorf(pkgError.PRIVATE_KEY_NOT_FOUND)
	}

	out, err := models.FindBlock(ctx, p.db, key)
	if err != nil {
		return nil, err
	}
	return out, nil
}
