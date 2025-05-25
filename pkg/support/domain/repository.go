package domain

import (
	"context"
	"tts-poc-service/models"
)

type CommandRepository interface {
	InsertSupport(ctx context.Context, in *models.Support) error
}
