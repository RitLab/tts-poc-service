package domain

import (
	"context"
	"tts-poc-service/models"
)

type QueryRepository interface {
	GetBlockById(ctx context.Context, key string) (*models.Block, error)
}

type CommandRepository interface {
	InsertBlockKey(ctx context.Context, in *models.Block) error
}

type OpenAIRepository interface {
	SummarizeText(ctx context.Context, in *OpenAIRequest) (string, error)
}
