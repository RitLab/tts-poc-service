package command

import (
	"context"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/models"
	"tts-poc-service/pkg/common/decorator"
	"tts-poc-service/pkg/support/domain"
)

type InsertSupportCommand struct {
	Name    string
	Email   string
	Message string
}

type InsertSupportHandler decorator.CommandHandler[InsertSupportCommand]

type insertSupportRepository struct {
	repo   domain.CommandRepository
	logger *baselogger.Logger
}

func NewInsertSupportRepository(repo domain.CommandRepository, log *baselogger.Logger) decorator.CommandHandler[InsertSupportCommand] {
	return decorator.ApplyCommandDecorators[InsertSupportCommand](
		insertSupportRepository{repo: repo, logger: log},
		log)
}

func (g insertSupportRepository) Handle(ctx context.Context, in InsertSupportCommand) error {
	return g.repo.InsertSupport(ctx, &models.Support{
		Name:    in.Name,
		Email:   in.Email,
		Message: in.Message,
	})
}
