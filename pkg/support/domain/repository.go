package domain

import (
	"context"
	"tts-poc-service/database/models"
)

type CommandRepository interface {
	InsertSupport(ctx context.Context, in *models.Support) error
}
